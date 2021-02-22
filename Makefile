.PHONY: test test-with-race

test:
	go test -i .
	go test .

test-with-race:
	go test -race .

build:
	go build -i -o mina .
