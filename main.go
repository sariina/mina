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

	fmt.Printf("Listening on localhost:%s\n", opts.Port)
	fmt.Printf("Redirect requests to %s\n", opts.Host)
	fmt.Printf("Cache responses to %s\n", absPath)

	http.HandleFunc("/", mina)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", opts.Port), nil))
}
