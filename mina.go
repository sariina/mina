package main

import (
	"crypto/md5"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type options struct {
	Port     string
	Host     string
	Verbose  bool
	CacheDir string
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func isCacheExist(filename string) bool {
	return isFileExist(filename)
}

func reqToMd5(req *http.Request) string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%+v", req))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func urlToFilename(relURL string) (filename string) {
	if relURL == "/" || relURL == "" {
		relURL = "/index.html"
	}

	base := filepath.Base(relURL)
	if len(base) > 255 {
		base = base[:255]
		log.Println("Warning: filename is truncated to 255 bytes")
	}
	relPath := filepath.Dir(relURL)

	filename = filepath.Join(opts.CacheDir, relPath, base)
	return
}

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

func requestHost(w http.ResponseWriter, req *http.Request, filename string) (body []byte, err error) {
	url := fmt.Sprintf("%s%s", opts.Host, req.URL.String())
	newreq, err := http.NewRequest(req.Method, url, req.Body)

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	res, err := client.Do(newreq)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ = ioutil.ReadAll(res.Body)
	return
}

func mina(w http.ResponseWriter, req *http.Request) {
	var body []byte
	var err error
	log.Printf("%s %s\n", req.Method, req.URL)

	filename := urlToFilename(req.URL.String())
	path := filepath.Dir(filename)

	if isCacheExist(filename) {
		body, err = cacheRead(filename)
	} else {
		body, err = requestHost(w, req, filename)
		go cacheWrite(path, filename, body)
	}

	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	w.Write(body)
}
