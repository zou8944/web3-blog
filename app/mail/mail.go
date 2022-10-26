package mail

import (
	"github.com/pkg/errors"
	"github.com/project5e/web3-blog/config"
	"github.com/project5e/web3-blog/pkg/logger"
	"github.com/project5e/web3-blog/pkg/mail"
	"github.com/project5e/web3-blog/pkg/mail/eml"
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
	messages, err := mail.DefaultMailer().Receive()
	if err != nil {
		logger.Errorf("Receive message fail. %v", err)
		return false
	}
	for _, message := range messages {
		onReceive(message)
	}
	return true
}

func onReceive(message mail.ReceivedMail) bool {
	r := strings.NewReader(message.Content)
	// convert string to eml.Email -> check format -> send backoff email if illegal format -> convert eml.Email to BlogMail if legal format
	email, err := eml.Parse(r)
	if err != nil {
		logger.Errorf("Convert message to eml.Email fail. message: %v, err: %v", message, err)
		return false
	}
	blog, err := convert2Blog(email)
	if err != nil {
		if errors.As(err, &formatErrorType) {
			logger.Infof("Convert email to blog fail. %v", err)
			SendEmailTemplate(config.Business.SupportEmail, email.From.Address, BadFormatTemplate)
			// notify success, even if send mail fail
			message.Notifier.Notify()
		}
		return false
	}
	if err := emailBlogHandleFun(blog); err == nil {
		message.Notifier.Notify()
		return false
	} else {
		return true
	}
}

func SendEmailTemplate(from, to string, template ResponseTemplate) bool {
	return mail.DefaultMailer().Send(from, to, template.Subject, template.Body)
}
