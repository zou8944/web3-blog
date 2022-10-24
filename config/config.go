package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type app struct {
	Server           server   `json:"server"`
	Logger           logger   `json:"logger"`
	AWS              aws      `json:"aws"`
	Database         database `json:"database"`
	Business         business `json:"business"`
	Web3StorageToken string   `json:"web3_storage_token"`
	JWT              jwt      `json:"jwt"`
}

type server struct {
	Port int `json:"port"`
}

type logger struct {
	Filename  string `json:"filename"`
	MaxSize   int    `json:"max_size"`
	MaxBackup int    `json:"max_backup"`
	MaxAge    int    `json:"max_age"`
	Compress  bool   `json:"compress"`
	LogType   string `json:"log_type"`
	Level     string `json:"level"`
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

var ENV string
var Server server
var Logger logger
var AWS aws
var Database database
var Business business
var Web3StorageToken string
var JWT jwt

func Parse() error {
	var config app
	if err := viper.Unmarshal(&config); err != nil {
		return errors.WithStack(err)
	}
	Server = config.Server
	Logger = config.Logger
	AWS = config.AWS
	Database = config.Database
	Business = config.Business
	Web3StorageToken = config.Web3StorageToken
	JWT = config.JWT
	ENV = viper.GetString("env")
	return nil
}
