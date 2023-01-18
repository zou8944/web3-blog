package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(bs []byte) string {
	digest := md5.New()
	digest.Write(bs)
	digestData := digest.Sum([]byte(nil))
	return hex.EncodeToString(digestData)
}

func MD5String(s string) string {
	digest := md5.New()
	digest.Write([]byte(s))
	digestData := digest.Sum([]byte(nil))
	return hex.EncodeToString(digestData)
}
