package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Mina struct {
	Target   *url.URL
	CacheDir string
	Headers  map[string]string
}

// newSingleHostReverseProxy is copied from stdlib, except we are change
// req.Host here to req.URL.Host
func newSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.URL.Host = target.Host
		req.Host = req.URL.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

// singleJoiningSlash is coped from stdlib, because it was called from
// newSingleHostReverseProxy.
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func writeRespToWR(wr http.ResponseWriter, resp *http.Response, headers map[string]string) {
	// write headers
	for name := range resp.Header {
		// overwrite custom headers
		if _, ok := headers[name]; !ok {
			wr.Header().Add(name, resp.Header.Get(name))
		}
	}
	for name, value := range headers {
		wr.Header().Add(name, value)
	}

	// write body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	wr.Write(body)
}

func cacheWrite(path string, filename string, body []byte) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Printf("Error while mkdir: %s", err)
		return
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		log.Printf("Error while writing: %s", err)
		return
	}
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

func (m *Mina) NewHandler() func(http.ResponseWriter, *http.Request) {
	p := newSingleHostReverseProxy(m.Target)
	return func(wr http.ResponseWriter, req *http.Request) {
		md5, reqDump := requestMD5(req)
		reqFilename := filepath.Join(m.CacheDir, fmt.Sprintf("%s.req", md5))
		resFilename := filepath.Join(m.CacheDir, fmt.Sprintf("%s.res", md5))

		if isFileExist(resFilename) {
			log.Printf("%s [HIT] %s %s", filepath.Base(resFilename)[:8], req.Method, req.URL)
			resDump, err := ioutil.ReadFile(resFilename)
			if err != nil {
				log.Println(err)
				return
			}

			dumpio := bufio.NewReader(bytes.NewBuffer(resDump))
			resp, err := http.ReadResponse(dumpio, req)
			if err != nil {
				log.Printf("Error: %s", err)
				return
			}
			defer resp.Body.Close()
			writeRespToWR(wr, resp, m.Headers)
		} else {
			log.Printf("%s [MISS] %s %s", filepath.Base(resFilename)[:8], req.Method, req.URL)

			wrRecorder := httptest.NewRecorder()
			p.ServeHTTP(wrRecorder, req)

			resp := wrRecorder.Result()
			defer resp.Body.Close()

			writeRespToWR(wr, resp, m.Headers)

			resDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				log.Printf("Error: %s", err)
				return
			}

			go cacheWrite(m.CacheDir, resFilename, resDump)
			go cacheWrite(m.CacheDir, reqFilename, reqDump)
		}
	}
}
