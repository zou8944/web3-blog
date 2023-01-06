package arweave

import (
	"testing"
)

func TestUpload(t *testing.T) {
	Init()
	data := []byte("# title \n ## Just for test")
	txId, err := UploadData(data)
	if err != nil {
		t.Error(err)
	}
	t.Logf("upload success, tx id: %s", txId)
}

func TestList(t *testing.T) {
	Init()
	posts, err := ListAllPost()
	if err != nil {
		t.Error(err)
	}
	t.Logf("list all post success, post: %+v", posts)
}
