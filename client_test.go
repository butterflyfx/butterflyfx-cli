package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func getRandomAddress() (*net.TCPAddr, error) {
	laddr, err := net.ResolveTCPAddr("tcp", "localhost:0") // Port == 0 - free
	if err != nil {
		return laddr, err
	}
	l, err := net.ListenTCP("tcp", laddr)
	l.Close()
	return l.Addr().(*net.TCPAddr), err
}

func setupTestServer() (addr *net.TCPAddr, err error) {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Success!")
	})
	addr, err = getRandomAddress()
	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		return
	}
	go http.Serve(listener, nil)
	return
}

func TestTunnelByProjectID(t *testing.T) {
	apiKey := os.Getenv("BUTTERFLYFX_API_KEY")
	projectID := 0
	tunnelAddress, _ := setupTestServer()
	err := TunnelByProjectID(apiKey, projectID, tunnelAddress.String(), true)
	time.Sleep(time.Second)
	if err != nil {
		t.Error(err)
	}

}
