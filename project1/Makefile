


all: project1



run:
	go run ./cmd/

project1:
	go build -o ./bin/project1 ./cmd/ 

test:
	go test ./...

clean:
	rm ./bin/*

.PHONY: clean run test