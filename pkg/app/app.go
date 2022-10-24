package app

import (
	"blog-web3/config"
	pkgConfig "blog-web3/pkg/config"
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
