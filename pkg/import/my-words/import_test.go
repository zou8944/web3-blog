package my_words

import (
	"fmt"
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	testRepoAddress := "https://github.com/zou8944/my-words"
	zipFileName, err := download(testRepoAddress)
	defer os.Remove(zipFileName)
	if err != nil {
		t.Errorf("%+v", err)
	} else {
		fmt.Printf("文件下载成功，路径为: %s", zipFileName)
	}
}

func TestExtract(t *testing.T) {
	testRepoAddress := "https://github.com/zou8944/my-words"
	zipFileName, err := download(testRepoAddress)
	defer os.Remove(zipFileName)
	if err != nil {
		t.Errorf("%+v", err)
	}
	dir, err := extract(zipFileName)
	defer os.RemoveAll(dir)
	if err != nil {
		t.Errorf("%+v", err)
	} else {
		fmt.Printf("文件解压成功，路径为: %s", dir)
	}
}
