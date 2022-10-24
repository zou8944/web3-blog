package main

import (
	"blog-web3/app/route"
	"blog-web3/bootstrap"
	"blog-web3/config"
	"blog-web3/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	bootstrap.SetupAll(g)
	route.RegisterRoutes(g)

	address := fmt.Sprintf("0.0.0.0:%d", config.Server.Port)
	logger.ErrorIf(g.Run(address))
}
