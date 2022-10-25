package mail

import (
	appConfig "blog-web3/config"
	"blog-web3/pkg/logger"
	"context"
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

func (m *AWSMailer) ReceiveMessage() ([]sqsTypes.Message, error) {
	getQueueUrlInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(appConfig.AWS.SQS.QueueName),
	}
	getQueueUrlOutput, err := m.SQS.GetQueueUrl(context.TODO(), getQueueUrlInput)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	receiveInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(sqsTypes.QueueAttributeNameAll),
		},
		QueueUrl:            getQueueUrlOutput.QueueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     appConfig.AWS.SQS.Timeout,
	}
	logger.Infof("Try receive message from sqs...")
	receiveOutput, err := m.SQS.ReceiveMessage(context.Background(), receiveInput)
	if err != nil {
		return nil, errors.WithStack(err)
	} else {
		return receiveOutput.Messages, nil
	}
}

func (m *AWSMailer) ReceiveMessageAsString() ([]string, error) {
	messages, err := m.ReceiveMessage()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var contents []string
	for _, message := range messages {
		contents = append(contents, *message.Body)
	}
	return contents, nil
}
