package app

import (
	"github.com/project5e/web3-blog/config"
	pkgConfig "github.com/project5e/web3-blog/pkg/config"
)

func IsLocal() bool {
	return config.ENV == pkgConfig.EnvLocal
}

func IsDev() bool {
	return config.ENV == pkgConfig.EnvDev
}

func IsTest() bool {
	return config.ENV == pkgConfig.EnvTest
}

func IsProduction() bool {
	return config.ENV == pkgConfig.EnvProduction
}
