BINARY = bin/songvote
SOURCE = cmd/songvote

lint:
	@golangci-lint run

build: lint
	@go build -o ${BINARY} ${SOURCE}/*.go

run: build
	@./${BINARY}

test: lint
	@go test -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html
