package mail

import (
	"blog-web3/config"
	"blog-web3/pkg/logger"
	"blog-web3/pkg/mail"
	"blog-web3/pkg/mail/eml"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type ResponseTemplate struct {
	Subject string
	Body    string
}

var (
	BadFormatTemplate = ResponseTemplate{
		Subject: "Bad format email",
		Body: `
The format of your email is incorrect, the correct format should be:

Subject: [action] #tag1 #tag2 #tagn title
	action should be one of [create, update, delete]
	there must be at least one tag: #private or #public, it will be treated as visibility. other tag will be treated as real tag.
	title will be treated as your blog title.
Body: raw markdown text, or plain text.
	body will be treated as your blog content.
`,
	}
)

var emailBlogHandleFun func(*BlogMail) error

func RegisterEmailHandler(handler func(*BlogMail) error) {
	emailBlogHandleFun = handler
}

func StartListenMailer() {
	go func() {
		for {
			if ok := receiveAndHandle(); !ok {
				time.Sleep(time.Second * 10)
			}
		}
	}()
}

func receiveAndHandle() bool {
	messages, err := mail.DefaultMailer().ReceiveMessageAsString()
	if err != nil {
		logger.Errorf("Receive Mail message fail. %v", err)
		return false
	}
	for _, message := range messages {
		onReceive(message)
	}
	return true
}

func onReceive(message string) bool {
	r := strings.NewReader(message)
	// convert string to eml.Email -> check format -> send backoff email if illegal format -> convert eml.Email to BlogMail if legal format
	email, err := eml.Parse(r)
	if err != nil {
		logger.Errorf("Convert message to eml.Email fail. %v", err)
		return false
	}
	blog, err := convert2Blog(email)
	if err != nil {
		if errors.As(err, formatErrorType) {
			SendEmailTemplate(config.Business.SupportEmail, email.From.Address, BadFormatTemplate)
		}
		return false
	}
	err = emailBlogHandleFun(blog)
	return err == nil
}

func SendEmailTemplate(from, to string, template ResponseTemplate) bool {
	return mail.DefaultMailer().Send(from, to, template.Subject, template.Body)
}
