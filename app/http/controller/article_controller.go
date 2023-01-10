package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/project5e/web3-blog/app/http/requests"
	"github.com/project5e/web3-blog/app/mail"
	"github.com/project5e/web3-blog/app/models"
	"github.com/project5e/web3-blog/pkg/logger"
	"github.com/project5e/web3-blog/pkg/response"
	"gorm.io/datatypes"
	"net/http"
)

type ArticleController struct{}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

func (ac *ArticleController) Create(c *gin.Context) {
	var body requests.ArticleCreateRequest
	if ok := requests.BindAndValidate(c, &body); !ok {
		return
	}
	model := models.Article{
		Title:   body.Title,
		Content: body.Content,
	}
	if ok := model.Create(); !ok {
		response.AbortWith500(c)
	}
	response.Created(c, model)
}

func (ar *ArticleController) Update(c *gin.Context) {
	var body requests.ArticleUpdateRequest
	if ok := requests.BindAndValidate(c, &body); !ok {
		return
	}
	articleId := c.Param("id")
	model := models.GetArticleById(articleId)
	if ok := model.Update(); !ok {
		response.AbortWith500(c)
	}
	response.Created(c, model)
}

func (ar *ArticleController) Delete(c *gin.Context) {
	articleId := c.Param("id")
	model := models.GetArticleById(articleId)
	if ok := model.Delete(); !ok {
		response.AbortWith500(c)
	}
	response.Created(c, model)
}

func (ac *ArticleController) ListPage(c *gin.Context) {
	articles := models.ListArticle()
	// let article content be article abstract
	for i, article := range articles {
		aContentRune := []rune(article.Content)
		if len(aContentRune) > 100 {
			articles[i].Content = string(aContentRune[:100])
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Articles": articles,
	})
}

func (ac *ArticleController) DetailPage(c *gin.Context) {
	articleId := c.Param("id")
	article := models.GetArticleById(articleId)
	if article == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
	} else {
		c.HTML(http.StatusOK, "detail.html", gin.H{
			"Article": article,
		})
	}
}

func (ac *ArticleController) HandleEmail(b *mail.BlogMail) error {
	switch b.Action {
	case mail.Create:
		tags, _ := json.Marshal(b.Tags)
		article := models.Article{
			Title:   b.Title,
			Content: b.Content,
			Tags:    datatypes.JSON(tags),
			Visible: b.Visible,
		}
		if ok := article.Create(); !ok {
			return errors.New("create article fail")
		}
	case mail.Update:
		tags, _ := json.Marshal(b.Tags)
		article := models.GetArticleByTitle(b.Title)
		if article == nil {
			logger.Warnf("No exist article found with title %s, ignore", article)
			return nil
		}
		article.Content = b.Content
		article.Tags = datatypes.JSON(tags)
		article.Visible = b.Visible
	case mail.Delete:
		logger.Warnf("Unsupported action: delete, ignore")
	default:
		logger.Errorf("Unknown action: %v, ignore", b.Action)
	}
	return nil
}
