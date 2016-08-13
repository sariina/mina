package main

import (
	"fmt"
	"log"
	neturl "net/url"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

type options struct {
	Port     string
	Host     string
	Verbose  bool
	CacheDir string
	Headers  map[string]string
}

const usage string = `Usage:
  mina --port=<port> --host=<host> [--output=<dir>] [--header=<header>]...

Options:
  -p --port=<port>     port to listen to
  -h --host=<host>     host to redirect to
  -H --header=<header> custom header
  -o --output=<dir>    [optional] Path to cache dir (default $(pwd)/<host>)

Example:
  mina -p 8080 -h https://www.domain.com:9000
`

func optionsFromArgs() (opts options) {
	args, err := docopt.Parse(usage, nil, true, "version 0.0.1", false, false)

	if err != nil || len(args) == 0 {
		os.Exit(1)
	}

	if args["--header"] != nil {
		opts.Headers = make(map[string]string)
		for _, header := range args["--header"].([]string) {
			keyVal := strings.SplitN(header, ":", 2)
			if len(keyVal) < 2 {
				continue
			}
			opts.Headers[keyVal[0]] = keyVal[1]
		}

	}

	if args["--port"] != nil {
		opts.Port = args["--port"].(string)
	}

	var url *neturl.URL

	if args["--host"] != nil {
		urlstr := args["--host"].(string)
		if !strings.HasPrefix(urlstr, "http") {
			urlstr = "http://" + urlstr
		}

		url, err = neturl.Parse(urlstr)
		if err != nil {
			log.Fatal(err)
		}

		if url.Host == "" {
			log.Fatal("Please provide a valid url e.g. http://yourdomain.com:1234")
		}
		opts.Host = url.String()
	}

	if args["--output"] != nil {
		opts.CacheDir = args["--output"].(string)
	} else {
		opts.CacheDir = fmt.Sprintf("%s", url.Host)
	}
	return
}
