package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	apiKey := flag.String("api-key", "foo", "The api key for your account")
	projectID := flag.Int("project", 0, "The id of the project you want to tunnel")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Not enough arguments")
	}
	command := args[0]
	fmt.Println(command)
	if command == "tunnel" {
		if len(args) > 1 {
			localAddr := args[1]
			err := TunnelByProjectID(*apiKey, *projectID, localAddr, false)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Not enough arguments")
		}
	}

}
