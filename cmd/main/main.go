package main

import (
	"blog-web3/internal/configs"
	"blog-web3/pkg/infra"
	"os"
)

func initAll() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "internal/configs/dev.yaml"
	}
	if err := configs.LoadConfigFile(configPath); err != nil {
		panic(err)
	}
	if err := infra.Init(); err != nil {
		panic(err)
	}
}

func main() {
	initAll()
}
