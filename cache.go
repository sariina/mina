package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/fatih/color"
)

func cacheWrite(path string, filename string, body []byte) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		color.Red("Error while mkdir: %s", err)
		return
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		color.Red("Error while writing: %s", err)
		return
	}
}

func cacheRead(filename string) (body []byte, err error) {
	body, err = ioutil.ReadFile(filename)
	return
}

func isCacheExist(filename string) bool {
	return isFileExist(filename)
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func requestMD5(req *http.Request) (string, []byte) {
	h := md5.New()
	body, _ := httputil.DumpRequest(req, true)
	io.WriteString(h, fmt.Sprintf("%+v", string(body)))

	return fmt.Sprintf("%x", h.Sum(nil)), body
}
