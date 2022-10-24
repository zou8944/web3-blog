package database

import (
	"blog-web3/config"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func Connect() error {
	dbLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		LogLevel: logger.Info,
	})
	newDB, err := gorm.Open(sqlite.Open(config.Database.SqliteFilePath), &gorm.Config{Logger: dbLogger})
	DB = newDB
	return errors.WithStack(err)
}
