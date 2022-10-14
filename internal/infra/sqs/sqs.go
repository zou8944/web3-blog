package sqs

import (
	"blog-web3/internal/configs"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
	"log"
)

var client *sqs.Client

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
	client = sqs.NewFromConfig(cfg)
	return nil
}

func ReceiveMessage() ([]types.Message, error) {
	getQueueUrlInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(configs.Conf.AWS.SNS.QueueName),
	}
	getQueueUrlOutput, err := client.GetQueueUrl(context.TODO(), getQueueUrlInput)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	receiveInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            getQueueUrlOutput.QueueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     configs.Conf.AWS.SNS.Timeout,
	}
	log.Println("Try receive message from sqs...")
	receiveOutput, err := client.ReceiveMessage(context.Background(), receiveInput)
	return receiveOutput.Messages, errors.WithStack(err)
}

func ReceiveMessageAsString() ([]string, error) {
	messages, err := ReceiveMessage()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var contents []string
	for _, message := range messages {
		contents = append(contents, *message.Body)
	}
	return contents, nil
}
