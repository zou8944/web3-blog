package route

import (
	"blog-web3/internal/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func RegisterRoutes(g *gin.Engine) {
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:           []string{"*"},
		AllowCredentials:       true,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
	}))

	g.POST("/users", controller.CreateUser)
	g.PUT("/users/:publicAddress", controller.OverrideUser)
	g.GET("/users/:publicAddress", controller.GetUser)
	g.POST("/users/login/metamask", controller.LoginWithMetaMask)
	g.GET("/articles", controller.ListArticle)

}
