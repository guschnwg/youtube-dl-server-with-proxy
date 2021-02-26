.PHONY: server

server:
	go build -ldflags "-s -w" -o ./app.out cmd/server/main.go

run:
	go run -ldflags "-s -w" cmd/server/main.go

docker-build:
	docker build -t youtube-dl-server .
docker-run:
	docker run -p 8000:8000 -e PORT=8000 -it youtube-dl-server