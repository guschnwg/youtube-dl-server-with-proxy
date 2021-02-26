package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	// newLocation := resp.Header.Get("Location")
	// resp.Header.Set("Location", "OI")
	resp.Header.Set("OI", "HELLOW")
	return resp, nil
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
		FlushInterval: -1,
		// Transport:     &transport{http.DefaultTransport},
	}

	http.HandleFunc("/videoplayback", proxy.ServeHTTP)
}
