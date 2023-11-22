BINARY = bin/songvote

lint:
	golangci-lint run

build: lint
	go build -o ${BINARY}	

run: build
	./${BINARY}

test: lint
	go test
