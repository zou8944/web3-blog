package my_words

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

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
