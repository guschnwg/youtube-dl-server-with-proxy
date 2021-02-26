package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
)

func YoutubePlay(w http.ResponseWriter, r *http.Request) {
	songURL := r.URL.Query().Get("url")
	fmt.Println(songURL)
	dec, _ := base64.URLEncoding.DecodeString(songURL)
	songURL = string(dec)

	fmt.Println(songURL)
	resp, err := http.Get(songURL)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("VAI SE FODER"))
		return
	}
	defer resp.Body.Close()

	fmt.Println("OI")
	fmt.Println(resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)

	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(body)
}

// YoutubeInfo ...
func YoutubeInfo(w http.ResponseWriter, r *http.Request) {
	songURL := r.URL.Query().Get("url")

	if songURL == "" {
		songURL = "https://www.youtube.com/watch?v=lgjY-lVtJZA"
	}

	w.Header().Set("Content-Type", "application/json")

	cmd := exec.Command("youtube-dl", "--skip-download", "--dump-json", "--extract-audio", "--audio-format", "mp3", "-4", songURL)
	fmt.Println(cmd.String())
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
			Ext string `json:"ext"`
			URL string `json:"url"`
		} `json:"formats"`
	}
	err = json.Unmarshal(out.Bytes(), &songData)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Error unmarshaling",
		})
		return
	}
	fmt.Println(songData)
	// lengthFormats := len(songData.Formats)
	// bestFormat := songData.Formats[lengthFormats-1]
	bestFormat := songData.Formats[0]
	for _, format := range songData.Formats {
		if format.Ext == "m4a" {
			bestFormat = format
		}
	}

	URL, _ := url.Parse(bestFormat.URL)
	URL.RawQuery += "&host=https://" + URL.Host
	URL.Host = "youtube-dl-alexa.herokuapp.com"
	// URL.Host = "75ef8d448436.ngrok.io"
	URL.Scheme = "https"

	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": map[string]interface{}{
			"proxy":        URL.String(),
			"url":          bestFormat.URL,
			"proxy_base64": base64.URLEncoding.EncodeToString([]byte(URL.String())),
			"url_base64":   base64.URLEncoding.EncodeToString([]byte(bestFormat.URL)),
		},
	})
}
