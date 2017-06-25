package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"errors"

	"github.com/koding/tunnel"
)

type HostRegistration struct {
	Identifier string `json:"identifier"`
	Subdomain  string `json:"subdomain"`
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
}

func TunnelByProjectID(apiKey string, projectId int, localAddr string, daemon bool) (err error) {
	req, err := http.NewRequest("GET", "https://www.butterflyfx.io/api/projects/"+strconv.Itoa(projectId)+"/tunnel", nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	httpClient := http.Client{
		CheckRedirect: func(redirRequest *http.Request, via []*http.Request) error {
			// Go's http.DefaultClient does not forward headers when a redirect 3xx
			// response is recieved. Thus, the header (which in this case contains the
			// Authorization token) needs to be passed forward to the redirect
			// destinations.
			redirRequest.Header = req.Header

			// Go's http.DefaultClient allows 10 redirects before returning an
			// an error. We have mimicked this default behavior.
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(resp.Status)
		}
		return errors.New(string(body))
	}
	registration := HostRegistration{}
	json.NewDecoder(resp.Body).Decode(&registration)
	if err != nil {
		return
	}
	cfg := &tunnel.ClientConfig{
		Identifier: registration.Identifier,
		ServerAddr: "tunnel.butterflyfx.io:80",
		Debug:      true,
		LocalAddr:  localAddr,
	}
	client, err := tunnel.NewClient(cfg)
	if err != nil {
		return
	}
	fmt.Println(registration.Subdomain)
	if daemon {
		go client.Start()
	} else {
		client.Start()
	}

	return
}
