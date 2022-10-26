package controller

import (
	"blog-web3/app/mail"
	"blog-web3/app/models"
	"blog-web3/pkg/logger"
	"blog-web3/pkg/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/datatypes"
)

type ArticleController struct{}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

func (ac *ArticleController) ListArticle(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		response.AbortWith400(c, errors.New("username in query is required"))
		return
	}
	articles := models.ListArticleByUser(username)
	response.SuccessWithData(c, articles)
}

func (ac *ArticleController) HandleEmail(b *mail.BlogMail) error {
	switch b.Action {
	case mail.Create:
		tags, _ := json.Marshal(b.Tags)
		article := models.Article{
			Username: b.UserName,
			Title:    b.Title,
			Content:  b.Content,
			Tags:     datatypes.JSON(tags),
			Visible:  b.Visible,
		}
		if ok := article.Create(); !ok {
			return errors.New("create article fail")
		}
	case mail.Update:
		tags, _ := json.Marshal(b.Tags)
		article := models.GetArticleByTitle(b.UserName, b.Title)
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
