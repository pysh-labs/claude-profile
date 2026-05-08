.PHONY: build test lint fmt vet cover

build:
	go build -o claude-profile ./cmd/claude-profile

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

cover:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
