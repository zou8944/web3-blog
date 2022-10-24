package main

import (
	"blog-web3/app/route"
	"blog-web3/bootstrap"
	"blog-web3/config"
	pkgConfig "blog-web3/pkg/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// load config from file to viper
	pkgConfig.LoadConfig()
	// parse config from viper to object, get config item by invoking config.Database, etc
	config.Parse()

	g := gin.New()
	bootstrap.SetupAll(g)
	route.RegisterRoutes(g)

	log.Fatalln("Http server start fail", g.Run(":9000"))
}
