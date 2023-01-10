package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// viper use mapstructure to unmarshal config to struct, instead json or yaml

type app struct {
	Server   server   `mapstructure:"server"`
	Logger   logger   `mapstructure:"logger"`
	AWS      aws      `mapstructure:"aws"`
	Database database `mapstructure:"database"`
	Business business `mapstructure:"business"`
	ArWeave  arWeave  `mapstructure:"arweave"`
	IPFS     ipfs     `mapstructure:"ipfs"`
	JWT      jwt      `mapstructure:"jwt"`
}

type server struct {
	Port int `mapstructure:"port"`
}

type logger struct {
	Filename  string `mapstructure:"filename"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
	MaxAge    int    `mapstructure:"max_age"`
	Compress  bool   `mapstructure:"compress"`
	LogType   string `mapstructure:"log_type"`
	Level     string `mapstructure:"level"`
}

type business struct {
	SupportEmail string `mapstructure:"support_email"`
}

type aws struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	SQS       awsSqs `mapstructure:"sqs"`
}

type awsSqs struct {
	QueueName string `mapstructure:"queue_name"`
	Timeout   int32  `mapstructure:"timeout"`
}

type database struct {
	Driver         string `mapstructure:"driver"`
	SqliteFilePath string `mapstructure:"sqlite_file_path"`
}

type jwt struct {
	SignKey string `mapstructure:"sign_key"`
}

type arWeave struct {
	Enable         bool   `mapstructure:"enable"`
	WalletKeyFile  string `mapstructure:"wallet_key_file"`
	Endpoint       string `mapstructure:"endpoint"`
	BundlrEndpoint string `mapstructure:"bundlr_endpoint"`
	AppName        string `mapstructure:"app_name"`
}

type ipfs struct {
	Enable bool   `mapstructure:"enable"`
	URL    string `mapstructure:"url"`
}

var ENV string
var Server server
var Logger logger
var AWS aws
var Database database
var Business business
var ArWeave arWeave
var JWT jwt
var IPFS ipfs

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
	ArWeave = config.ArWeave
	IPFS = config.IPFS
	JWT = config.JWT
	ENV = viper.GetString("env")
	return nil
}
