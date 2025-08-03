.PHONY: fmt test run tidy

all: tidy fmt test

fmt:
	go fmt ./...

test:
	go test ./... -v -count=1 2>&1 | tee test.out


run:
	go run ./cmd/api

tidy:
	go mod tidy
