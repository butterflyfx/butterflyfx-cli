package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func newDirector(r *http.Request) func(*http.Request) {
	return func(req *http.Request) {
		req.Host = r.Header.Get("X-Forwarded-Host")
		req.URL.Host = req.Host + ":" + r.Header.Get("X-Forwarded-Port")
		req.URL.Scheme = r.Header.Get("X-Forwarded-Proto")

		reqLog, err := httputil.DumpRequestOut(req, false)
		if err != nil {
			log.Printf("Got error %s\n %+v\n", err.Error(), req)
		}

		log.Println(string(reqLog))
	}
}

func GenerateProxyServer() (net.Listener, error) {
	addr, err := GetRandomAddress()
	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req.Host = r.Header.Get("X-Forwarded-Host")
			req.URL.Host = req.Host + ":" + r.Header.Get("X-Forwarded-Port")
			req.URL.Scheme = r.Header.Get("X-Forwarded-Proto")

			reqLog, err := httputil.DumpRequestOut(req, false)
			if err != nil {
				log.Printf("Got error %s\n %+v\n", err.Error(), req)
			}

			log.Println(string(reqLog))
		}
		proxy := &httputil.ReverseProxy{
			Transport: &http.Transport{},
			Director:  director,
		}
		proxy.ServeHTTP(w, r)
	})
	go http.Serve(listener, nil)
	return listener, err
}
