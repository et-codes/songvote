SONGVOTE_BINARY = bin/songvote
SONGVOTE_SOURCE = .

lint:
	@golangci-lint run

build: lint
	@go build -o ${SONGVOTE_BINARY} ${SONGVOTE_SOURCE}/*.go

run: build
	@./${SONGVOTE_BINARY}

test: lint
	@go test -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

commit: test
	@git add .
	@git commit
	@git push
