package configs

import "github.com/spf13/viper"

type AppConfig struct {
	AWS AWSConf `json:"aws"`
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

var Conf *AppConfig

func LoadConfigFile(filepath string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath)
	return viper.Unmarshal(Conf)
}
