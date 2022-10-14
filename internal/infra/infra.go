package infra

import (
	"blog-web3/internal/infra/ses"
	"blog-web3/internal/infra/sqs"
)

func Init() error {
	err1 := ses.Init()
	err2 := sqs.Init()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
