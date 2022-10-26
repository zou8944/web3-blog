package bootstrap

import (
	"fmt"
	"github.com/project5e/web3-blog/app/http/controller"
	"github.com/project5e/web3-blog/app/mail"
)

func SetupMailer() {
	ac := controller.ArticleController{}

	mail.RegisterEmailHandler(ac.HandleEmail)
	mail.StartListenMailer()
	fmt.Println("Mailer started")
}
