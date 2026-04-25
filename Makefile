.PHONY: run test lint vet fmt build tidy

run:
	air

test:
	go test ./... -cover

lint:
	golangci-lint run

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

build:
	go build -o ./tmp/server.exe ./cmd/server
