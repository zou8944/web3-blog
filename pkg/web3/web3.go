package web3

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/project5e/web3-blog/pkg/logger"
	"strings"
)

func VerifySignature(publicAddress, signature, message string) bool {
	messageBytes := accounts.TextHash([]byte(message))
	signatureBytes := hexutil.MustDecode(signature)

	if signatureBytes[crypto.RecoveryIDOffset] == 27 || signatureBytes[crypto.RecoveryIDOffset] == 28 {
		signatureBytes[crypto.RecoveryIDOffset] -= 27
	}
	recoveredPublicKey, err := crypto.SigToPub(messageBytes, signatureBytes)
	if err != nil {
		logger.Debugf("web3 signature verify fail. public_address:%s, message:%s, signature:%s", publicAddress, message, signature)
		return false
	}
	recoveredPublicAddress := crypto.PubkeyToAddress(*recoveredPublicKey).Hex()

	return strings.ToUpper(publicAddress) == strings.ToUpper(recoveredPublicAddress)
}
