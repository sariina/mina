package main

import (
	"os"
	"testing"
	"time"
)

func TestGCMemory(t *testing.T) {
	cacheDir := os.TempDir()
	m := newInMemoryCache(cacheDir)

	m.GCInterval = time.Second
	m.ExpirationTime = 2 * time.Second
	go m.GC()

	m.Store("sample-key", []byte("Hello"))

	time.Sleep(time.Second)

	if !m.Exists("sample-key") {
		t.Fatal("key not exists")
	}

	time.Sleep(2 * time.Second)

	if m.Exists("sample-key") {
		t.Fatal("key still exists")
	}
}

func TestStoreLoadMemory(t *testing.T) {
	cacheDir := os.TempDir()
	m := newInMemoryCache(cacheDir)
	want := []byte("Hello")

	m.Store("sample-key", want)

	bytes := m.Load("sample-key")
	if string(bytes) != string(want) {
		t.Fatalf("got %q want %q", string(bytes), string(want))
	}
}

func TestExistsMemory(t *testing.T) {
	cacheDir := os.TempDir()
	m := newInMemoryCache(cacheDir)
	want := []byte("Hello")

	if m.Exists("key") {
		t.Fatal("key has been detected :|")
	}

	m.Store("key", want)

	if !m.Exists("key") {
		t.Fatal("where is the key ! , key not founded")
	}
}
