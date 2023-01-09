package route

import (
	"github.com/gin-gonic/gin"
	"github.com/project5e/web3-blog/app/http/controller"
)

func RegisterRoutes(g *gin.Engine) {
	ac := controller.NewArticleController()
	uc := controller.NewUserController()

	g.POST("/users", uc.CreateUser)
	g.PUT("/users/:publicAddress", uc.OverrideUser)
	g.GET("/users/:publicAddress", uc.GetUser)
	g.POST("/users/login/metamask", uc.LoginWithMetaMask)
	g.GET("/articles", ac.List)
	g.POST("/articles", ac.Create)
	g.PUT("/articles/:id", ac.Update)
	g.DELETE("/articles/:id", ac.Delete)
}
