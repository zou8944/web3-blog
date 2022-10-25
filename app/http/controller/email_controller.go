package controller

import (
	"blog-web3/app/mail"
	"blog-web3/pkg/logger"
)

type EmailController struct{}

func NewEmailController() *EmailController {
	return &EmailController{}
}

func (ec *EmailController) HandleEmail(b *mail.BlogMail) error {
	logger.Infof("%+v", b)
	return nil
}
