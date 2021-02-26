package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

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
		FlushInterval: -1,
	}

	http.HandleFunc("/videoplayback", proxy.ServeHTTP)
}
