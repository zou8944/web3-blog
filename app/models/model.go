package models

import "github.com/project5e/web3-blog/pkg/database"

func Migrate() error {
	return database.DB.AutoMigrate(&User{}, &Article{})
}
