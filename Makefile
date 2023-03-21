mod:
	go mod download
	go mod tidy

build: mod
	go build -o ./bin/main

tests:
	go test -coverprofile=coverage.out -cover -race ./... -count=1

testcache:
	go test -coverprofile=coverage.out -v ./cache
	go tool cover -html=cache/coverage.out

linters:
	golangci-lint run --config=.golangci.yml

start:
	docker-compose up
