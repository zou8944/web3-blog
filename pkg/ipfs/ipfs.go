package ipfs

import (
	"bytes"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/project5e/web3-blog/config"
)

const IPFS_GATWAY = "https://gateway.ipfs.io/ipfs/"

var ipfsClient *shell.Shell

func Init() {
	if config.IPFS.Enable {
		ipfsClient = shell.NewShell(config.IPFS.URL)
	}
}

func UploadData(b []byte) (string, error) {
	if config.IPFS.Enable {
		cid, err := ipfsClient.Add(bytes.NewReader(b))
		if err != nil {
			return "", err
		}
		return IPFS_GATWAY + cid, err
	}
	return "", nil
}
