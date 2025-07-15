.PHONY: fmt test

fmt:
	go fmt ./...

test:
	go test ./...

run:
	go run ./cmd/api

