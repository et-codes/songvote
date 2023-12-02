SONGVOTE_BINARY = bin/songvote
SONGVOTE_SOURCE = cmd/songvote
SEED_BINARY = bin/seed
SEED_SOURCE = cmd/seed

lint:
	@golangci-lint run

build: lint
	@go build -o ${SONGVOTE_BINARY} ${SONGVOTE_SOURCE}/*.go
	@go build -o ${SEED_BINARY} ${SEED_SOURCE}/*.go

run: build
	@./${SONGVOTE_BINARY}

test: lint
	@go test -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html

commit: test
	@git add .
	@git commit
	@git push
