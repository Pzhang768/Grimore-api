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

build:
	go build -o ./tmp/server.exe ./cmd/server
