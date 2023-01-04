package ipfs

import (
	"bytes"
	shell "github.com/ipfs/go-ipfs-api"
)

const IPFS_GATWAY = "https://gateway.ipfs.io/ipfs/"

var ipfsClient *shell.Shell

func Init() {
	ipfsClient = shell.NewShell("localhost:5001")
}

func UploadData(b []byte) (string, error) {
	cid, err := ipfsClient.Add(bytes.NewReader(b))
	if err != nil {
		return "", nil
	}
	return IPFS_GATWAY + cid, err
}
