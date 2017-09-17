package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
)

func TestGenerateProxy(t *testing.T) {
	listener, err := GenerateProxyServer()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(listener.Addr().String())
	req, err := http.NewRequest("GET", "http://"+listener.Addr().String()+"/status/418", nil)
	req.Header.Add("X-Forwarded-Host", "httpbin.org")
	req.Header.Add("X-Forwarded-Port", "80")
	req.Header.Add("X-Forwarded-Proto", "http")
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 418 {
		t.Error("Failed to proxy hostname")
	}
	addr, err := GetRandomAddress()
	go http.ListenAndServe(addr.String(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	}))
	if err != nil {
		t.Error(err)
	}
	req, err = http.NewRequest("GET", "http://"+listener.Addr().String(), nil)
	req.Header.Add("X-Forwarded-Host", "127.0.0.1")
	req.Header.Add("X-Forwarded-Port", strconv.Itoa(addr.Port))
	req.Header.Add("X-Forwarded-Proto", "http")
	resp, err = httpClient.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	if string(body) != "Hello world" {
		t.Error("Failed to proxy port")
	}
}
