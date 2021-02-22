.PHONY: server

server:
	go build -ldflags "-s -w" -o ./app.out cmd/server/main.go