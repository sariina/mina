// Mina is an HTTP reverse proxy server with cache.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	opts options
)

func main() {
	opts = optionsFromArgs()
	absPath, err := filepath.Abs(opts.CacheDir)

	if err != nil {
		log.Fatalln(err)
	}

	if !isFileExist(opts.CacheDir) {
		err = os.Mkdir(opts.CacheDir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("  Addr: http://localhost:%s\n", opts.Port)
	fmt.Printf("  Host: %s\n", opts.Host)
	fmt.Printf(" Cache: %s\n", absPath)

	http.HandleFunc("/", mina)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", opts.Port), nil))
}
