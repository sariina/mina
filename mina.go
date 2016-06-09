package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

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
