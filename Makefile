SRC=./...

default: test build

build:
	go build

test:
	go test -v $(SRC)
