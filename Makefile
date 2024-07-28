all: build

run:
	go run ./cmd/main.go

request:
	sh create.sh
