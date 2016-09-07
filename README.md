# Mina

Mina is saves API server responses to disk and serves them with its own HTTP server. Some use cases include:

- working with the API server when you or the server is offline
- each request costs money
- API server forces rate limits 
- etc.

## Install

    go get -u github.com/sariina/mina

## Example

Start a mina server for Github API on port 8080:

    mina -addr=:8080 -target=https://api.github.com

In your client/broweser/app, instead of sending a request to

    https://api.github.com/users/sariina

send it to

    http://localhost:8080/users/sariina

Voila, the same response. The response is saved to your disk.
Your app will think that you are using Github API even when you are offline.

## Options

    mina --help

    Usage:
      mina -addr=<addr> -target=<target> [-o=<dir>] [-H=<header>]...
    
    Options:
      -addr    address to listen to
      -target  target to route to
      -H       custom header
      -o       [optional] cache dir
    
    Example:
      mina -addr=:8080 -target=https://www.domain.com:9000

## Why mina?

Mina is named after the
[Myna bird](https://en.wikipedia.org/wiki/Common_hill_myna)
(Persian: [مرغ مینا](https://fa.wikipedia.org/wiki/%D9%85%DB%8C%D9%86%D8%A7%DB%8C_%D9%85%D8%B9%D9%85%D9%88%D9%84%DB%8C)),
renowned for their ability to imitate speech.
