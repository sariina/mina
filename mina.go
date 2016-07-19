package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"time"

	"github.com/fatih/color"
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

func getResponseDump(w http.ResponseWriter, req *http.Request, filename string) (dump []byte, err error) {
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

	dump, err = httputil.DumpResponse(res, true)
	// body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	return
}

func mina(w http.ResponseWriter, req *http.Request) {
	var dump []byte
	var err error
	var resp *http.Response
	log.Printf("%s %s\n", req.Method, req.URL)

	filename := reqToMd5Filename(req)
	path := filepath.Dir(filename)

	if isCacheExist(filename) {
		dump, err = cacheRead(filename)
	} else {
		dump, err = getResponseDump(w, req, filename)
	}

	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	dumpio := bufio.NewReader(bytes.NewBuffer(dump))
	resp, err = http.ReadResponse(dumpio, req)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	go cacheWrite(path, filename, dump)

	for name, _ := range resp.Header {
		w.Header().Add(name, resp.Header.Get(name))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	w.Write(body)
}
