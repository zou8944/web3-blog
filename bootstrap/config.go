package bootstrap

import (
	"fmt"
	"github.com/project5e/web3-blog/config"
	pkgConfig "github.com/project5e/web3-blog/pkg/config"
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
	fmt.Println("Config load success")
}
