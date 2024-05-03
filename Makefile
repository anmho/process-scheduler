


all: linux-amd64 darwin-arm64 build

build:
	go build -o ./bin/project1-linux-amd64 ./cmd/

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ./bin/project1-linux-amd64 ./cmd/

run:
	go run ./cmd/

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/project1-darwin-arm64 ./cmd/

test:
	go test ./...

clean:
	rm ./bin/*

.PHONY: clean linux-amd64 darwin-arm64 build