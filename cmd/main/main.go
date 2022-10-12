package main

import (
	"blog-web3/internal/configs"
	"blog-web3/pkg/infra"
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
}

func main() {
	initAll()
}
