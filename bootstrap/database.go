package bootstrap

import (
	"fmt"
	"github.com/project5e/web3-blog/app/models"
	"github.com/project5e/web3-blog/pkg/database"
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
