package my_words

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/project5e/web3-blog/app/models"
	"github.com/project5e/web3-blog/pkg/logger"
	"github.com/project5e/web3-blog/pkg/util"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Import(gitAddress string) error {
	repoZip, err := download(gitAddress)
	if err != nil {
		return errors.WithStack(err)
	}
	repoDir, err := extract(repoZip)
	defer os.RemoveAll(repoDir)
	if err != nil {
		return errors.WithStack(err)
	}
	articles, err := scanArticles(repoDir)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, article := range articles {
		if article.Slug == "" {
			logger.Errorf("Article has no slug, ignore: %s", article.Title)
			continue
		}
		persistent(&article)
	}
	return nil
}

// download git repository as a temp zip file, return file path
func download(gitAddress string) (string, error) {
	zipUrl := fmt.Sprintf("%s/archive/refs/heads/main.zip", gitAddress)
	resp, err := http.Get(zipUrl)
	if err != nil {
		return "", errors.WithStack(err)
	}
	zipFile, err := os.CreateTemp(os.TempDir(), "web3-blog-import-my-words-*")
	defer zipFile.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}
	_, _ = zipFile.Write(bytes)
	return zipFile.Name(), nil
}

// extract zip file to a dir, return dir name
func extract(zipFileName string) (string, error) {
	zor, err := zip.OpenReader(zipFileName)
	defer zor.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	zipDir := fmt.Sprintf("%s/%s-extract", os.TempDir(), filepath.Base(zipFileName))
	if err := os.Mkdir(zipDir, os.ModePerm); err != nil {
		return "", errors.WithStack(err)
	}
	for _, f := range zor.File {
		if f.FileInfo().IsDir() {
			err := os.Mkdir(zipDir+"/"+f.Name, os.ModePerm)
			if err != nil {
				return "", errors.WithStack(err)
			}
		} else {
			_f, err := os.Create(zipDir + "/" + f.Name)
			if err != nil {
				return "", errors.WithStack(err)
			}
			fr, err := f.Open()
			if err != nil {
				return "", errors.WithStack(err)
			}
			fbs, err := io.ReadAll(fr)
			if err != nil {
				return "", errors.WithStack(err)
			}
			err = ioutil.WriteFile(_f.Name(), fbs, os.ModeAppend)
			if err != nil {
				return "", errors.WithStack(err)
			}
		}
	}
	return zipDir, nil
}

// scan all markdown article
func scanArticles(dir string) ([]Article, error) {
	var articles []Article
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".MD") {
			file, err := os.Open(path)
			if err != nil {
				logger.Errorf("File read fail: %s. error: %+v", path, err)
				return nil
			}
			article, err := readArticle(file)
			if err != nil {
				logger.Errorf("Article read fail: %s. error: %+v", path, err)
				return nil
			}
			articles = append(articles, *article)
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].CreatedAt.Before(articles[j].CreatedAt)
	})
	return articles, nil
}

type Article struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

// read from file, convert to Article
func readArticle(f *os.File) (*Article, error) {
	title := filepath.Base(f.Name())
	scanner := bufio.NewScanner(f)
	var prefix string
	var meta string
	var content string
	// 扫描到第一个---时meta开始，此时meta前面不应该有任何内容，否则该文章没有meta；扫描到第二个---时meta结束
	indicator := "prefix"
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch indicator {
		case "prefix":
			if line == "---" {
				if len(prefix) != 0 {
					return nil, errors.New(fmt.Sprintf("Article has no meta field: %s", title))
				}
				indicator = "meta"
				continue
			}
			prefix = prefix + line
		case "meta":
			if line == "---" {
				indicator = "content"
				continue
			}
			meta = meta + line + "\n"
		case "content":
			content = content + line + "\n"
		}
	}
	content = strings.Trim(content, "\n")
	var article Article
	article.Title = title
	article.Content = content
	if err := fillArticleWithMeta(&article, meta); err != nil {
		return nil, err
	}
	return &article, nil
}

func fillArticleWithMeta(article *Article, metaString string) error {
	scanner := bufio.NewScanner(strings.NewReader(metaString))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		switch strings.Split(line, ":")[0] {
		case "created_at":
			t, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(strings.Split(line, "created_at:")[1]))
			if err != nil {
				return errors.WithStack(err)
			}
			article.CreatedAt = t
		case "updated_at":
			t, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(strings.Split(line, "updated_at:")[1]))
			if err != nil {
				return errors.WithStack(err)
			}
			article.UpdatedAt = t
		case "slug":
			article.Slug = strings.TrimSpace(strings.Split(line, "slug:")[1])
		}
	}
	return nil
}

func persistent(article *Article) {
	articleMD5 := util.MD5String(article.Title + article.Slug + article.Content)
	if _article := models.GetArticleBySlug(article.Slug); _article == nil {
		_newArticle := models.Article{
			Title:     article.Title,
			Content:   article.Content,
			Slug:      article.Slug,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		}
		if ok := _newArticle.Create(); !ok {
			logger.Errorf("Article create fail. title: %s", article.Title)
		}
	} else {
		_articleMD5 := util.MD5String(_article.Title + _article.Slug + _article.Content)
		if articleMD5 == _articleMD5 {
			logger.Infof("Article exist and no modification, ignore: %s", article.Title)
			return
		}
		_article.Title = article.Title
		_article.Content = article.Content
		_article.CreatedAt = article.CreatedAt
		if ok := _article.Update(); !ok {
			logger.Errorf("Article update fail. title: %s", article.Title)
		} else {
			logger.Infof("Article updated: %s", article.Title)
		}
	}
}
