package database

import (
	"github.com/pkg/errors"
	"github.com/project5e/web3-blog/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func Connect() error {
	dbLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		LogLevel: logger.Error,
	})
	newDB, err := gorm.Open(sqlite.Open(config.Database.SqliteFilePath), &gorm.Config{Logger: dbLogger})
	DB = newDB
	return errors.WithStack(err)
}
