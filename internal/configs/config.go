package configs

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AWS      AWSConf  `json:"aws"`
	Business Business `json:"business"`
}

type Business struct {
	SupportEmail string `json:"support_email"`
}

type AWSConf struct {
	Region    string     `json:"region"`
	AccessKey string     `json:"access_key"`
	SecretKey string     `json:"secret_key"`
	SNS       AWSSQSConf `json:"sns"`
}

type AWSSQSConf struct {
	QueueName string `json:"queue_name"`
	Timeout   int32  `json:"timeout"`
}

var Conf AppConfig

func LoadConfigFile(filepath string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return errors.WithStack(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
