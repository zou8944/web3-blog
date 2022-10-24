package route

import (
	"blog-web3/app/http/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine) {
	ac := controller.NewArticleController()
	uc := controller.NewUserController()

	g.POST("/users", uc.CreateUser)
	g.PUT("/users/:publicAddress", uc.OverrideUser)
	g.GET("/users/:publicAddress", uc.GetUser)
	g.POST("/users/login/metamask", uc.LoginWithMetaMask)
	g.GET("/articles", ac.ListArticle)
}
