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
	var hit = false

	md5, reqDump := requestMD5(req)
	reqFilename := filepath.Join(opts.CacheDir, fmt.Sprintf("%s.req", md5))
	resFilename := filepath.Join(opts.CacheDir, fmt.Sprintf("%s.res", md5))

	hit = isCacheExist(resFilename)

	if hit {
		dump, err = cacheRead(resFilename)
	} else {
		dump, err = getResponseDump(w, req, resFilename)
	}

	if hit {
		log.Printf("%s [HIT] %s %s", filepath.Base(resFilename)[:8], req.Method, req.URL)
	} else {
		log.Printf("%s [MISS] %s %s", filepath.Base(resFilename)[:8], req.Method, req.URL)
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
	defer resp.Body.Close()

	if !hit {
		go cacheWrite(opts.CacheDir, resFilename, dump)
		go cacheWrite(opts.CacheDir, reqFilename, reqDump)
	}

	for name, _ := range resp.Header {
		if _, ok := opts.Headers[name]; ok {
			continue
		}
		w.Header().Add(name, resp.Header.Get(name))
	}
	for name, value := range opts.Headers {
		w.Header().Add(name, value)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	w.Write(body)
}
