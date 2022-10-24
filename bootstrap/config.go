package bootstrap

import (
	"blog-web3/config"
	pkgConfig "blog-web3/pkg/config"
	"fmt"
)

func SetupConfig() {
	// load config from file to viper
	if err := pkgConfig.LoadConfigFile(); err != nil {
		panic(fmt.Sprintf("Load config file fail. %v", err))
	}
	// parse config from viper to object, get config item by invoking config.Database, etc
	if err := config.Parse(); err != nil {
		panic(fmt.Sprintf("Parse config fail. %v", err))
	}
}
