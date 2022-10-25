package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	EnvDefault    = "default"
	EnvLocal      = "local"
	EnvDev        = "dev"
	EnvTest       = "test"
	EnvQA         = "qa"
	EnvProduction = "prod"
)

// LoadConfigFile
// if env present, load {env}.yaml, else, load default.yaml
func LoadConfigFile() error {
	env := determineEnv()

	cfgDir := "config"
	cfgName := EnvDefault
	switch env {
	case "":
		viper.Set("env", EnvDefault)
	case EnvLocal, EnvDev, EnvTest, EnvQA, EnvProduction:
		cfgPath := fmt.Sprintf("%s/%s.yaml", cfgDir, env)
		if _, err := os.Stat(cfgPath); err != nil {
			// todo it should be print with logger
			fmt.Printf("file %s not exist, load config/default.yaml", cfgPath)
			break
		}
		cfgName = env
		viper.Set("env", env)
	default:
		return fmt.Errorf("unknown env: %s", env)
	}
	cfgPath := fmt.Sprintf("%s/%s.yaml", cfgDir, cfgName)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Read Config file failed. %v", err)
	}
	return nil
}

// retrieve env from --env, environment variable
func determineEnv() string {
	if envFlag := flag.Lookup("env"); envFlag != nil {
		return envFlag.Value.String()
	}
	return os.Getenv("ENV")
}
