package server

import (
	"bytes"
	"encoding/json"
	"net/http"
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
			"error": err.Error(),
		})
		return
	}

	var songData map[string]json.RawMessage
	err = json.Unmarshal(out.Bytes(), &songData)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": songData,
	})
}
