package web3

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"log"
)

func VerifySignature(publicAddr, signature, message string) bool {
	sig := hexutil.MustDecode(signature)
	msg := accounts.TextHash([]byte(message))
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27
	}
	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		log.Println(errors.WithStack(err))
		return false
	}
	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	return publicAddr == recoveredAddr.Hex()
}
