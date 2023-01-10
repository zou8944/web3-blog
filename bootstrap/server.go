package bootstrap

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/project5e/web3-blog/pkg/app"
	"time"
)

func SetupServer(engine *gin.Engine) {
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:           []string{"*"},
		AllowCredentials:       true,
		MaxAge:                 12 * time.Hour,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
	}))
	engine.LoadHTMLGlob("templates/*")
	if app.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	fmt.Println("Logger load success")
}
