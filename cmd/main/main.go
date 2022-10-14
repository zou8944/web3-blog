package main

import (
	"blog-web3/internal/configs"
	"blog-web3/internal/model"
	"blog-web3/internal/route"
	"blog-web3/pkg/infra"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func initAll() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "internal/configs/dev.yaml"
	}
	if err := configs.LoadConfigFile(configPath); err != nil {
		log.Fatalf("%+v", err)
	}
	if err := infra.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	if err := model.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func main() {
	initAll()
	g := gin.New()
	route.RegisterRoutes(g)
	log.Fatalln("Http server start fail", g.Run(":9000"))
}
