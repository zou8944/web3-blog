package config

import (
	"blog-web3/pkg/logger"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	EnvLocal      = "local"
	EnvDev        = "dev"
	EnvTest       = "test"
	EnvQA         = "qa"
	EnvProduction = "prod"
)

// LoadConfig Rule
// if env present, load {env}.yaml, load default.yaml else
func LoadConfig() {
	env := determineEnv()

	cfgDir := "config"
	cfgName := "default"
	switch env {
	case EnvLocal, EnvDev, EnvTest, EnvQA, EnvProduction:
		cfgPath := fmt.Sprintf("%s/%s.yaml", cfgDir, env)
		if _, err := os.Stat(cfgPath); err != nil {
			logger.Warn(fmt.Sprintf("File %s not exist, load config/default.yaml", cfgPath))
			break
		}
		cfgName = env
	default:
		panic("Unknown env: " + env)
	}
	cfgPath := fmt.Sprintf("%s/%s.yaml", cfgDir, cfgName)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Read Config file failed. %v", err))
	}
}

// retrieve env from --env, environment variable
func determineEnv() string {
	if envFlag := flag.Lookup("env"); envFlag != nil {
		return envFlag.Value.String()
	}
	return os.Getenv("ENV")
}
