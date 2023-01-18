package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/project5e/web3-blog/app/route"
	"github.com/project5e/web3-blog/bootstrap"
	"github.com/project5e/web3-blog/config"
	"github.com/project5e/web3-blog/pkg/logger"
)

func main() {
	g := gin.New()
	bootstrap.SetupAll(g)
	route.RegisterRoutes(g)
	//if err := my_words.Import("https://github.com/zou8944/my-words"); err != nil {
	//	panic(err)
	//}
	address := fmt.Sprintf("0.0.0.0:%d", config.Server.Port)
	logger.ErrorIf(g.Run(address))
}
