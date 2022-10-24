package bootstrap

import "blog-web3/pkg/database"

func SetupDatabase() {
	database.Connect()
}
