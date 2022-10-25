package bootstrap

import (
	"blog-web3/pkg/database"
	"fmt"
)

func SetupDatabase() {
	if err := database.Connect(); err != nil {
		panic(fmt.Sprintf("Database Connect fail. %v", err))
	}
	fmt.Println("Database load success")
}
