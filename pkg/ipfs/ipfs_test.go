package ipfs

import "testing"

func TestUpload(t *testing.T) {
	Init()
	_, err := UploadData([]byte("Hello"))
	if err != nil {
		t.Error(err)
	}
}
