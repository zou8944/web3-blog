package arweave

import (
	"fmt"
	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"io"
	"log"
	"os"
)

func Init() {
	wallet, err := goar.NewWalletFromPath("./keyfile.json", "https://arweave.net")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	file, err := os.Open("./demo.txt")
	bytes, err := io.ReadAll(file)

	t, err := wallet.SendData(
		bytes,
		[]types.Tag{
			{
				Name:  "name",
				Value: "value",
			},
		})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", t)
	fmt.Printf("%s\n", t.ID)
	tx, err := wallet.Client.GetTransactionByID("1DgEGr9ufH4ANjXyBGg6R8oGOSCQkcAi6O0gvP1i8hY")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", tx)
	fmt.Printf("%s\n", tx.ID)
	fmt.Printf("")
}
