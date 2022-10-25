package bootstrap

import (
	"blog-web3/config"
	"blog-web3/pkg/logger"
	"fmt"
)

func SetupLogger() {
	err := logger.Init(
		config.Logger.Filename,
		config.Logger.MaxSize,
		config.Logger.MaxBackup,
		config.Logger.MaxAge,
		config.Logger.Compress,
		config.Logger.LogType,
		config.Logger.Level,
	)
	if err != nil {
		panic(fmt.Sprintf("Log initialize fail. %v", err))
	}
	fmt.Println("Logger load success")
}
