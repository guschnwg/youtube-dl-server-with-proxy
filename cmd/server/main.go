package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/guschnwg/youtube-dl-server-with-proxy/pkg/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": "OI",
		})
	})
	http.HandleFunc("/api/info", server.YoutubeInfo)

	server.BindProxy()

	fmt.Printf("Listening in port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
