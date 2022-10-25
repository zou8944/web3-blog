package bootstrap

import (
	"blog-web3/app/http/controller"
	"blog-web3/app/mail"
	"fmt"
)

func SetupMailer() {
	ec := controller.NewEmailController()

	mail.RegisterEmailHandler(ec.HandleEmail)
	mail.StartListenMailer()
	fmt.Println("Mailer started")
}
