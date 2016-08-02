# Mina

Mina is a single binary HTTP reverse proxy that repeates your HTTP requests to another host and caches the response.

## Install

    go get github.com/sariina/mina

## Example

To start a server on port 8080 and redirects all requests to www.bing.com run this command:

    mina -p 8080 -h http://www.bing.com

make a request to a resource in bing.com e.g. http://www.bing.com/?scope=news in your browser of choice.
now make a request to http://localhost:8080/?scope=news

voila, the same response.


## Options

    mina --help

    Usage:
      mina --port=<port> --host=<host> [--output=<dir>]
    
    Options:
      -p --port=<port>  Port to listen to.
      -h --host=<host>  Host to redirect to.
      -o --output=<dir> [optional] Path to cache dir.
                        Default: $(pwd)/<host>
    
    Example:
      mina -p 8080 -h https://www.domain.com:9000

## Why mina?

Mina is named after the
[Myna bird](https://en.wikipedia.org/wiki/Common_hill_myna)
(Persian: [مرغ مینا](https://fa.wikipedia.org/wiki/%D9%85%DB%8C%D9%86%D8%A7%DB%8C_%D9%85%D8%B9%D9%85%D9%88%D9%84%DB%8C)),
renowned for their ability to imitate speech.
