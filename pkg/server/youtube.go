package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os/exec"
)

// YoutubeInfo ...
func YoutubeInfo(w http.ResponseWriter, r *http.Request) {
	songURL := r.URL.Query().Get("url")

	if songURL == "" {
		songURL = "https://www.youtube.com/watch?v=lgjY-lVtJZA"
	}

	w.Header().Set("Content-Type", "application/json")

	cmd := exec.Command("youtube-dl", songURL, "--skip-download", "--dump-json", "-4")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Error running youtuble-dl query",
		})
		return
	}

	var songData struct {
		Formats []struct {
			URL string `json:"url"`
		} `json: "formats"`
	}
	err = json.Unmarshal(out.Bytes(), &songData)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Error unmarshaling",
		})
		return
	}

	lengthFormats := len(songData.Formats)
	bestFormat := songData.Formats[lengthFormats-1]

	URL, _ := url.Parse(bestFormat.URL)
	URL.RawQuery += "&host=https://" + URL.Host
	URL.Host = "localhost:8000"
	URL.Scheme = "http"

	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": map[string]string{
			"proxy": URL.String(),
			"url":   bestFormat.URL,
		},
	})
}
