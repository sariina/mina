package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const usage string = `Usage:
  mina -addr=<addr> -target=<target> [-o=<dir>] [-H=<header>]...

Options:
  -addr    address to listen to
  -target  target to route to
  -H       custom header
  -o       [optional] cache dir

Example:
  mina -addr=:8080 -target=https://www.domain.com:9000`

type colonSeparatedFlags map[string]string

func (h colonSeparatedFlags) String() string {
	return "string representation"
}

func (h colonSeparatedFlags) Set(value string) error {
	keyVal := strings.SplitN(value, ":", 2)
	if len(keyVal) != 2 {
		return nil
	}
	h[keyVal[0]] = keyVal[1]
	return nil
}

func main() {
	var (
		flagListen   = flag.String("addr", "", "address to listen to")
		flagTarget   = flag.String("target", "", "taget to route to")
		flagCacheDir = flag.String("o", "", "path to cache dir")
		flagHeaders  = make(colonSeparatedFlags)
	)
	flag.Var(&flagHeaders, "H", "custom header, e.g. 'Key: Value'")

	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	if len(*flagTarget) == 0 || len(*flagListen) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// flagTarget
	targetURL, err := url.Parse(*flagTarget)
	if err != nil {
		log.Fatal(err)
	}

	// flagCacheDir
	if *flagCacheDir == "" {
		*flagCacheDir = fmt.Sprintf("%s", targetURL.Host)
	}
	*flagCacheDir, err = filepath.Abs(*flagCacheDir)
	if err != nil {
		log.Fatalln(err)
	}
	if !isFileExist(*flagCacheDir) {
		err = os.Mkdir(*flagCacheDir, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	ln, err := net.Listen("tcp", *flagListen)
	if err != nil {
		log.Fatal(err)
	}

	// info
	fmt.Printf("Address: %v\n", *flagListen)
	fmt.Printf(" Target: %s\n", targetURL.String())
	fmt.Printf("  Cache: %s\n", *flagCacheDir)

	// Serve
	m := Mina{
		Target:   targetURL,
		CacheDir: *flagCacheDir,
		Headers:  flagHeaders,
	}
	http.HandleFunc("/", m.ServeHTTP)
	log.Fatal(http.Serve(ln, nil))
}
