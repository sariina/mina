package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
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
		color.Red("Error: %s", err)
		return
	}

	newreq.Header = req.Header

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	res, err := client.Do(newreq)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	return
}

func mina(w http.ResponseWriter, req *http.Request) {
	var body []byte
	var err error
	log.Printf("%s %s\n", req.Method, req.URL)

	filename := reqToMd5Filename(req)
	path := filepath.Dir(filename)

	if isCacheExist(filename) {
		body, err = cacheRead(filename)
		if err != nil {
			color.Red("Error: %s", err)
			return
		}
	} else {
		body, err = requestHost(w, req, filename)
		if err != nil {
			color.Red("Error: %s", err)
			return
		}
		go cacheWrite(path, filename, body)
	}

	w.Write(body)
}
