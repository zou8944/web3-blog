package mail

import (
	appConfig "blog-web3/config"
	"blog-web3/pkg/logger"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
)

type AWSMailer struct {
	Context context.Context
	SES     *sesv2.Client
	SQS     *sqs.Client
}

type SQSNotifier struct {
	ReceiptHandle *string
	NotifyFun     func(*string)
}

func NewAWSMailer() *AWSMailer {
	_config, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(appConfig.AWS.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     appConfig.AWS.AccessKey,
				SecretAccessKey: appConfig.AWS.SecretKey,
			},
		}),
	)
	if err != nil {
		logger.Errorf("Create config for ses fail. %v", err)
		return nil
	}
	return &AWSMailer{
		Context: context.Background(),
		SES:     sesv2.NewFromConfig(_config),
		SQS:     sqs.NewFromConfig(_config),
	}
}

func (m *AWSMailer) Send(sender, recipient, subject, content string) bool {
	in := &sesv2.SendEmailInput{
		Content: &sesTypes.EmailContent{
			Simple: &sesTypes.Message{
				Body: &sesTypes.Body{
					Text: &sesTypes.Content{
						Data:    aws.String(content),
						Charset: aws.String("UTF-8"),
					},
				},
				Subject: &sesTypes.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
			},
		},
		Destination: &sesTypes.Destination{
			ToAddresses: []string{recipient},
		},
		FromEmailAddress: aws.String(sender),
		ReplyToAddresses: []string{sender},
	}
	out, err := m.SES.SendEmail(m.Context, in)
	if err != nil {
		logger.Errorf("Send Driver fail. %v", err)
		return false
	}
	logger.Infof("Email sent, message id: %s", *out.MessageId)
	return true
}

func (m *AWSMailer) getQueueUrl() *string {
	req := &sqs.GetQueueUrlInput{
		QueueName: aws.String(appConfig.AWS.SQS.QueueName),
	}
	res, err := m.SQS.GetQueueUrl(m.Context, req)
	if err != nil {
		logger.Errorf("Get queue url fail. %v", err)
		return nil
	}
	return res.QueueUrl
}

// Receive message from sqs, and convert to ReceivedMail
func (m *AWSMailer) Receive() ([]ReceivedMail, error) {
	messages, err := m._receiveMessages()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var rMessages []ReceivedMail
	for _, message := range messages {
		body := *message.Body

		var snsMessage map[string]string
		var mailMessage map[string]interface{}
		_ = json.Unmarshal([]byte(body), &snsMessage)
		_ = json.Unmarshal([]byte(snsMessage["Message"]), &mailMessage)

		contentB64 := mailMessage["content"].(string)
		contentBytes, err := base64.StdEncoding.DecodeString(contentB64)
		if err != nil {
			logger.Warnf("Message content base64 decode fail. using raw message")
			contentBytes = []byte(contentB64)
		}
		rMessage := ReceivedMail{
			Content:  string(contentBytes),
			Notifier: m.NewSQSNotifier(message.ReceiptHandle),
		}
		rMessages = append(rMessages, rMessage)
	}
	return rMessages, nil
}

// get sqs message
func (m *AWSMailer) _receiveMessages() ([]sqsTypes.Message, error) {
	queueUrl := m.getQueueUrl()
	if queueUrl == nil {
		return nil, nil
	}
	req := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(sqsTypes.QueueAttributeNameAll),
		},
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     appConfig.AWS.SQS.Timeout,
	}
	logger.Infof("Try receive message from sqs...")
	res, err := m.SQS.ReceiveMessage(m.Context, req)
	if err != nil {
		return nil, errors.WithStack(err)
	} else {
		messages := res.Messages
		if len(messages) > 0 {
			logger.Debugf("received sqs message: %+v", messages)
		}
		return messages, nil
	}
}

func (m *AWSMailer) deleteMessage(receiptHandle *string) {
	queueUrl := m.getQueueUrl()
	if queueUrl == nil {
		return
	}
	req := &sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: receiptHandle,
	}
	_, err := m.SQS.DeleteMessage(m.Context, req)
	if err != nil {
		logger.Errorf("Delete message fail. receipt handle: %v, err: %v", receiptHandle, err)
		return
	}
	logger.Infof("Message deleted: %v", receiptHandle)
}

func (m *AWSMailer) NewSQSNotifier(receiptHandle *string) *SQSNotifier {
	return &SQSNotifier{
		ReceiptHandle: receiptHandle,
		NotifyFun:     m.deleteMessage,
	}
}

func (n *SQSNotifier) Notify() {
	n.NotifyFun(n.ReceiptHandle)
}
