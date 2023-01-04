package bootstrap

import "github.com/project5e/web3-blog/pkg/ipfs"

func SetupIPFS() {
	ipfs.Init()
}
