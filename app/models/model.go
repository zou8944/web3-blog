package models

import "blog-web3/pkg/database"

func Migrate() error {
	return database.DB.AutoMigrate(&User{})
}
