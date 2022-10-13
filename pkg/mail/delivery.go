package mail

import (
	"blog-web3/internal/configs"
	"blog-web3/pkg/infra/ses"
	"blog-web3/pkg/infra/sqs"
	"github.com/pkg/errors"
	"log"
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

var emailBlogHandleFun func(*BlogEmail) error

func RegisterEmailHandler(handler func(*BlogEmail) error) {
	emailBlogHandleFun = handler
}

func StartListenSQS() {
	go func() {
		for {
			if err := receiveFromSQSAndHandle(); err != nil {
				log.Println("Receive and parse email error")
				time.Sleep(time.Second * 10)
			}
		}
	}()
}

func receiveFromSQSAndHandle() error {
	emails, err := sqs.ReceiveMessageAsString()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, email := range emails {
		if err = onReceive(email); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func onReceive(content string) error {
	r := strings.NewReader(content)
	// 解析成邮件 -> 判断是否符合格式 -> 如果不符合格式，发送退回邮件 -> 如果符合格式，解析成待发布的博客结构
	email, err := parseEmail(r)
	if err != nil {
		return errors.WithStack(err)
	}
	blog, err := convert2Blog(email)
	if err != nil {
		if errors.As(err, formatErrorType) {
			err = SendEmailTemplate(configs.Conf.Business.SupportEmail, email.From.Address, BadFormatTemplate)
			if err != nil {
				return errors.WithStack(err)
			}
		} else {
			return errors.WithStack(err)
		}
	}
	err = emailBlogHandleFun(blog)
	return errors.WithStack(err)
}

func SendEmailTemplate(from, to string, template ResponseTemplate) error {
	return SendEmail(from, to, template.Subject, template.Body)
}

func SendEmail(from, to, subject, content string) error {
	err := ses.SendMail(from, to, subject, content)
	return errors.WithStack(err)
}
