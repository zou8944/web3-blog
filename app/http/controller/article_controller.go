package controller

import "github.com/gin-gonic/gin"

type ArticleController struct{}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

func (ac *ArticleController) ListArticle(c *gin.Context) {

}
