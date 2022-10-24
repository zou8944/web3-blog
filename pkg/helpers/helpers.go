package helpers

import (
	"blog-web3/pkg/logger"
	"github.com/google/uuid"
)

func GenerateNonce() string {
	_uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Errorf("Generate nonce fail. %v", err)
		return ""
	} else {
		return _uuid.String()
	}
}
