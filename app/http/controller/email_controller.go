package controller

import "blog-web3/app/mail"

type EmailController struct{}

func NewEmailController() *EmailController {
	return &EmailController{}
}

func (ec *EmailController) HandleEmail(b *mail.BlogMail) error {
	return nil
}
