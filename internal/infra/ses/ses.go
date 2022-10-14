package ses

import (
	"blog-web3/internal/configs"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/pkg/errors"
	"log"
)

var client *sesv2.Client

func Init() error {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(configs.Conf.AWS.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     configs.Conf.AWS.AccessKey,
				SecretAccessKey: configs.Conf.AWS.SecretKey,
			},
		}),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	client = sesv2.NewFromConfig(cfg)
	return nil
}

func SendMail(sender, recipient, subject, content string) error {
	in := &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data:    aws.String(content),
						Charset: aws.String("UTF-8"),
					},
				},
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		FromEmailAddress: aws.String(sender),
		ReplyToAddresses: []string{sender},
	}
	out, err := client.SendEmail(context.Background(), in)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf("Email sent, message id: %s", *out.MessageId)
	return nil
}
