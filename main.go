package main

import (
	"flag"
	"log"
	"strings"
)

func logErrorOrDie(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	apiKey := flag.String("api-key", "foo", "The api key for your account")
	projectID := flag.Int("project", 0, "The id of the project you want to tunnel")
	flag.Parse()
	args := flag.Args()
	command := args[0]
	if command == "tunnel" {
		localAddr := ""
		if len(args) > 1 {
			localAddr = args[1]
		}
		if localAddr == "" {
			listener, err := GenerateProxyServer()
			logErrorOrDie(err)
			localAddr = listener.Addr().String()
		}
		err := TunnelByProjectID(*apiKey, *projectID, localAddr, false)
		logErrorOrDie(err)
	} else if command == "install" {
		InstallChromeManifest()
	} else if strings.HasPrefix(command, "chrome") {
		StartChromeNativeClient()
	}

}
