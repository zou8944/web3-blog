package helpers

import (
	"github.com/google/uuid"
	"github.com/project5e/web3-blog/pkg/logger"
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
