package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type app struct {
	AWS              aws      `json:"aws"`
	Database         database `json:"database"`
	Business         business `json:"business"`
	Web3StorageToken string   `json:"web3_storage_token"`
	JWT              jwt      `json:"jwt"`
}

type business struct {
	SupportEmail string `json:"support_email"`
}

type aws struct {
	Region    string `json:"region"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	SNS       awsSqs `json:"sns"`
}

type awsSqs struct {
	QueueName string `json:"queue_name"`
	Timeout   int32  `json:"timeout"`
}

type database struct {
	Driver         string `json:"driver"`
	SqliteFilePath string `json:"sqlite_file_path"`
}

type jwt struct {
	SignKey string `json:"sign_key"`
}

var AWS aws
var Database database
var Business business
var Web3StorageToken string
var JWT jwt

func Parse() {
	var config app
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("Parse Config failed. %v", err))
	}
	AWS = config.AWS
	Database = config.Database
	Business = config.Business
	Web3StorageToken = config.Web3StorageToken
	JWT = config.JWT
}
