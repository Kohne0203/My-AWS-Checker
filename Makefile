.PHONY: build test clean

BINARY_NAME = awscheck

build:
    go build -o $(BINARY_NAME)

test:
    go test -v ./...

clean:
    rm -f $(BINARY_NAME)

install:
    go install

.DEFAULT_GOAL := build
