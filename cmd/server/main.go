package main

import (
	"log"
	"net/http"
	"os"

	"github.com/guschnwg/youtube-dl-server-with-proxy/pkg/server"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/api/info", server.YoutubeInfo)

	server.BindProxy()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
