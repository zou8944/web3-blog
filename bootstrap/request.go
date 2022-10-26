package bootstrap

import (
	"fmt"
	"github.com/project5e/web3-blog/app/http/requests"
)

func SetupRequest() {
	requests.Init()
	fmt.Println("validator load success")
}
