package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func checkResponse(body string, wantBody string, header string, wantHeader string) error {
	if body != wantBody {
		return fmt.Errorf("body got %q; want %q", body, wantBody)
	}
	if header != wantHeader {
		return fmt.Errorf("header got %q; want %q", header, wantHeader)
	}
	return nil
}

func TestMina(t *testing.T) {
	want := []byte("tweet")
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(want)
	}))
	defer backend.Close()

	cacheDir := os.TempDir()
	url, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	m := &Mina{
		Target:   url,
		CacheDir: cacheDir,
		Headers:  map[string]string{},
	}

	frontend := httptest.NewServer(m)
	defer frontend.Close()

	// first time
	gotBody, gotHeader := get(frontend.URL)
	err = checkResponse(gotBody, string(want), gotHeader, XHeaderValueMiss)
	if err != nil {
		t.Fatal(err)
	}

	// second time
	gotBody, gotHeader = get(frontend.URL)
	err = checkResponse(gotBody, string(want), gotHeader, XHeaderValueHit)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIgnoreHeader(t *testing.T) {
	want := []byte("tweet")
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(want)
	}))
	defer backend.Close()

	cacheDir := os.TempDir()
	url, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	m := &Mina{
		Target:   url,
		CacheDir: cacheDir,
		Headers:  map[string]string{},
	}

	frontend := httptest.NewServer(m)
	defer frontend.Close()

	req, err := http.NewRequest("GET", frontend.URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add(RequestOptionsHeaderName, XHeaderValueIgnore)

	client := http.Client{}
	firstResponse, err := client.Do(req) // first time to call
	if err != nil {
		log.Fatal(err)
	}
	defer firstResponse.Body.Close()
	buf, err := ioutil.ReadAll(firstResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	gotBody, gotHeader := string(buf), firstResponse.Header.Get(XHeaderName)
	err = checkResponse(gotBody, string(want), gotHeader, XHeaderValueMiss)
	if err != nil {
		t.Fatal(err)
	}

	secondResponse, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer secondResponse.Body.Close()
	buf, err = ioutil.ReadAll(secondResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	gotBody, gotHeader = string(buf), secondResponse.Header.Get(XHeaderName)
	err = checkResponse(gotBody, string(want), gotHeader, XHeaderValueMiss)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotModifiedStatusCode(t *testing.T) {
	firstTime := true

	want := []byte("tweet")
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if firstTime {
			// for example
			w.Header().Set("If-Modified-Since", "Sun, 22 Nov 2015 07:10:45 GMT")
			w.Header().Set("If-None-Match", "W/\"6a6248d0dfd45c6aac94b6ad02c856bb")
			firstTime = false
		} else {
			w.WriteHeader(http.StatusNotModified)
		}
		w.Write(want)
	}))
	defer backend.Close()

	cacheDir := os.TempDir()
	url, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	m := &Mina{
		Target:   url,
		CacheDir: cacheDir,
		Headers:  map[string]string{},
	}

	frontend := httptest.NewServer(m)
	defer frontend.Close()

	// first time with If-Modified-Since header
	gotBody, gotHeader := get(frontend.URL)
	err = checkResponse(gotBody, string(want), gotHeader, XHeaderValueMiss)
	if err != nil {
		t.Fatal(err)
	}

	// second time
	gotBody, gotHeader = get(frontend.URL)
	err = checkResponse(gotBody, string(want), gotHeader, XHeaderValueHit)
	if err != nil {
		t.Fatal(err)
	}
}

func get(url string) (body string, xHeaderValue string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf), res.Header.Get(XHeaderName)
}
