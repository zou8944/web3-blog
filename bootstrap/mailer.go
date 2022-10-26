package bootstrap

import (
	"blog-web3/app/http/controller"
	"blog-web3/app/mail"
	"fmt"
)

func SetupMailer() {
	ac := controller.ArticleController{}

	mail.RegisterEmailHandler(ac.HandleEmail)
	mail.StartListenMailer()
	fmt.Println("Mailer started")
}
