.PHONY: fmt test run tidy

all: tidy fmt test

fmt:
	go fmt ./...

test:
	go test ./...

run:
	go run ./cmd/api

tidy:
	go mod tidy
