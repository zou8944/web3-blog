package arweave

import (
	"encoding/json"
	"fmt"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const walletKeyFile = "./keyfile.json"
const arweaveNode = "https://arseed.web3infra.dev"
const bundlrNode = "https://node2.bundlr.network"

const appName = "4tune-web3-blog"

const graphQLFmt = `
query {
  transactions(
    first: %d,
    after: "%s",
    tags: [
        {
            name: "App-Name",
            values: ["4tune-web3-blog"]
        }
    ],
    owners: [
        "%s"
    ]
  ) 
  {
    edges {
      cursor
      node {
        id
        anchor
        owner {
            address
            key
        }
        data {
            size
            type
        }
        tags {
            name
            value
        }
        parent {
            id
        }
      }
    }
  }
}
`

var wallet *goar.Wallet
var itemSigner *goar.ItemSigner

func Init() {
	if _wallet, err := goar.NewWalletFromPath(walletKeyFile, arweaveNode); err != nil {
		log.Fatalf("%+v", err)
	} else {
		wallet = _wallet
	}
	if _itemSigner, err := goar.NewItemSigner(wallet.Signer); err != nil {
		log.Fatalf("%+v", err)
	} else {
		itemSigner = _itemSigner
	}
}

func createTags() []types.Tag {
	return []types.Tag{
		{"Content-Type", "text/plain"},
		{"App-Name", appName},
		{"App-Version", "v1.0"},
		{"Unix-Time", strconv.FormatInt(time.Now().Unix(), 10)},
	}
}

// UploadPost  upload post, then return transaction id of arweave
func UploadPost(post string) (string, error) {
	return UploadData([]byte(post))
}

// UploadData upload data with bytes
func UploadData(b []byte) (string, error) {
	tags := createTags()
	item, err := itemSigner.CreateAndSignItem(b, "", "", tags)
	if err != nil {
		return "", err
	}
	resp, err := utils.SubmitItemToBundlr(item, bundlrNode)
	if err != nil {
		return "", err
	}
	fmt.Printf("%+v", resp)
	return resp.Id, nil
}

// GetPost get post from arweave by transaction id
func GetPost(txId string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", arweaveNode, txId))
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", errors.New(fmt.Sprintf("response error when get post. txid: %s, resp: %+v", txId, resp))
	}
	postBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(postBytes), nil
}

// ListAllPost list all posts of specific wallet address
func ListAllPost() ([]string, error) {
	txIds, err := listPostTxId()
	if err != nil {
		return nil, err
	}
	var posts []string
	for _, txId := range txIds {
		post, err := GetPost(txId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func listPostTxId() ([]string, error) {
	type graphQLRes struct {
		Transactions struct {
			Edges []struct {
				Cursor string `json:"cursor"`
				Node   struct {
					ID    string `json:"id"`
					Owner struct {
						Address string `json:"address"`
						Key     string `json:"key"`
					} `json:"owner"`
					Tags []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"tags"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"transactions"`
	}

	var txIds []string

	limit := 100
	cursor := ""
	owner := wallet.Signer.Address
	for true {
		resBytes, err := wallet.Client.GraphQL(fmt.Sprintf(graphQLFmt, limit, cursor, owner))
		if err != nil {
			return nil, err
		}
		var res graphQLRes
		err = json.Unmarshal(resBytes, &res)
		if err != nil {
			return nil, err
		}
		edges := res.Transactions.Edges
		for _, edge := range edges {
			txIds = append(txIds, edge.Node.ID)
		}
		if len(edges) < limit {
			break
		}
		cursor = edges[len(edges)-1].Cursor
	}

	return txIds, nil
}
