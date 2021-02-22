package server

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type override struct {
	Header string
	Match  string
	Host   string
	Path   string
}

type config struct {
	Path     string
	Host     string
	Override override
}

// BindProxy ...
func BindProxy() {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			URL := r.URL.Query().Get("host")
			if URL == "" {
				URL = "r4---sn-8p8v-bg0sl.googlevideo.com"
			}
			parsedURL, _ := url.Parse(URL)

			originHost := parsedURL.Host
			r.Header.Add("X-Forwarded-Host", r.Host)
			r.Header.Add("X-Origin-Host", originHost)
			r.Host = originHost
			r.URL.Host = originHost
			r.URL.Scheme = "https"
		},
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 50 * time.Second,
			}).Dial,
		},
	}

	http.HandleFunc("/videoplayback", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}
