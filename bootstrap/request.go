package bootstrap

import (
	"blog-web3/app/http/requests"
	"fmt"
)

func SetupRequest() {
	requests.Init()
	fmt.Println("validator load success")
}
