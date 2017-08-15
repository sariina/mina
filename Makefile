build:
	go build -i -o mina
test:
	go test -v .
install:
	mv mina /usr/local/bin