package bootstrap

import (
	"blog-web3/app/models"
	"blog-web3/pkg/database"
	"fmt"
)

func SetupDatabase() {
	if err := database.Connect(); err != nil {
		panic(fmt.Sprintf("Database Connect fail. %v", err))
	}
	fmt.Println("Database connect success")
	if err := models.Migrate(); err != nil {
		panic(fmt.Sprintf("Database migrate fail. %v", err))
	}
	fmt.Println("Database migrate success")
}
