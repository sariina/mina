package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

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
	if gotBody != string(want) {
		t.Fatalf("got %q; want %q", gotBody, string(want))
	}
	if gotHeader != XHeaderValueMiss {
		t.Fatalf("got %q; want %q", gotBody, XHeaderValueMiss)
	}

	// second time
	gotBody, gotHeader = get(frontend.URL)
	if gotBody != string(want) {
		t.Fatalf("got %q; want %q", gotBody, string(want))
	}
	if gotHeader != XHeaderValueHit {
		t.Fatalf("got %q; want %q", gotBody, XHeaderValueHit)
	}
}

func TestNotModifiedStatusCode(t *testing.T) {
	// Ignoring sync.Mutex for two requests that are not concurrent
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
	if gotBody != string(want) {
		t.Fatalf("got %q; want %q", gotBody, string(want))
	}
	if gotHeader != XHeaderValueMiss {
		t.Fatalf("got %q; want %q", gotBody, XHeaderValueMiss)
	}

	// second time
	gotBody, gotHeader = get(frontend.URL)
	if gotBody != string(want) {
		t.Fatalf("got %q; want %q", gotBody, string(want))
	}
	if gotHeader != XHeaderValueHit {
		t.Fatalf("got %q; want %q", gotBody, XHeaderValueHit)
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
