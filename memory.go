package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"
)

var (
	inMemoryCache *Memory
)

// Bytes is array of bytes for requests or responses
type Bytes []byte

// Cache , for storing time and data
type Cache struct {
	CreatedAt time.Time
	Data      Bytes
}

// Memory is main struct for in-memory cache
type Memory struct {
	Table                map[string]Cache
	CacheDir             string
	ExpirationTime       time.Duration
	GCInterval           time.Duration
	MaximumContentLength int64
	sync.RWMutex
}

func newInMemoryCache(cacheDir string) *Memory {
	return &Memory{
		CacheDir: cacheDir,
		Table:    make(map[string]Cache),

		// fill default values.
		// TODO: fill these using flags.
		MaximumContentLength: 256 * 1024,
		ExpirationTime:       600 * time.Second,
		GCInterval:           60 * time.Second,
	}
}

// Store will save content in map
func (m *Memory) Store(hash string, cache Bytes) {
	m.Lock()
	defer m.Unlock()

	// use m.Exists(hash) before storing to memory

	m.Table[hash] = Cache{
		CreatedAt: time.Now(),
		Data:      cache,
	}
}

// Update will save content in map
func (m *Memory) Update(hash string) {
	m.Lock()
	defer m.Unlock()

	// use m.Exists(hash) before storing to memory

	item := m.Table[hash]
	item.CreatedAt = time.Now()
	m.Table[hash] = item
}

// Exists return true if cache is available on memory
func (m *Memory) Exists(hash string) bool {
	m.RLock()
	defer m.RUnlock()

	_, ok := m.Table[hash]
	return ok
}

// Load return response bytes
func (m *Memory) Load(hash string) Bytes {
	m.RLock()
	defer m.RUnlock()

	// are you sure that data exists ?

	cache := m.Table[hash]
	return cache.Data
}

// GC is garbage collector that will remove expired caches and store them in file.
func (m *Memory) GC() {
	m.Lock()
	defer m.Unlock()

	for key, value := range m.Table {
		if time.Since(value.CreatedAt) > m.ExpirationTime {
			resFilename := filepath.Join(m.CacheDir, fmt.Sprintf("%s.res", key))
			err := m.writeCache(resFilename, value.Data)
			if err != nil {
				log.Printf("\033[0;31mError on saving cache : %s\033[0m", err)
				continue
			}
			delete(m.Table, key)
		}
	}
	time.AfterFunc(m.GCInterval, m.GC)
}

// Shutdown will erase all caches
func (m *Memory) Shutdown() {
	m.Lock()
	defer m.Unlock()

	for key, value := range m.Table {
		resFilename := filepath.Join(m.CacheDir, fmt.Sprintf("%s.res", key))
		err := m.writeCache(resFilename, value.Data)
		if err != nil {
			log.Printf("\033[0;31mError on saving cache : %s\033[0m", err)
			continue
		}
		delete(m.Table, key)
	}
}

func (m *Memory) writeCache(filename string, body Bytes) (err error) {
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		log.Printf("Error while writing: %s", err)
		return
	}
	return nil
}
