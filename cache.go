package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

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

func requestMD5(req *http.Request) string {
	h := md5.New()
	headers := headerToSortedString(req.Header)
	io.WriteString(h, fmt.Sprintf("%+v", req.Method))
	io.WriteString(h, fmt.Sprintf("%+v", req.URL))
	io.WriteString(h, fmt.Sprintf("%+v", headers))

	body, _ := ioutil.ReadAll(req.Body)
	io.WriteString(h, fmt.Sprintf("%+v", string(body)))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	md5 := fmt.Sprintf("%x", h.Sum(nil))

	return filepath.Join(opts.CacheDir, md5)
}

func headerToSortedString(header http.Header) (ret string) {
	arr := make([]string, len(header))
	i := 0
	for key, _ := range header {
		arr[i] = key
		i++
	}
	sort.Strings(arr)

	for _, key := range arr {
		ret += fmt.Sprintf("%s:%s\n", key, header[key])
	}
	return
}
