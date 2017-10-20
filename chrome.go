/*
* @Author: thngo
* @Date:   2015-12-02 14:32:06
* @Last Modified by:   thngo
* @Last Modified time: 2015-12-04 17:06:42
 */

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
)

var StdOut = os.Stdout

type Message struct {
	Text    string
	Action  string
	Payload string
}

type TunnelAction struct {
	ApiKey    string
	Address   string
	ProjectId int
}

func StartChromeNativeClient() {
	read()
}
func captureOutput() {
	r, w, _ := os.Pipe()
	os.Stdout = w
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
	}()
}

func releaseOutput() {
	os.Stdout = StdOut
}

func read() {
	captureOutput()
	for {
		s := bufio.NewReader(os.Stdin)
		length := make([]byte, 4)
		s.Read(length)
		lengthNum := readMessageLength(length)
		content := make([]byte, lengthNum)
		s.Read(content)
		parseMessage(content)
	}
}

func parseMessage(msg []byte) {
	var content Message
	json.Unmarshal(msg, &content)
	if content.Action == "Hi" {
		send(Message{Text: "Oh hello there!"})
	} else if content.Action == "tunnel" {
		var payload TunnelAction
		json.Unmarshal([]byte(content.Payload), &payload)
		log.Println(content.Payload)
		if payload.Address == "" {
			listener, err := GenerateProxyServer()
			logErrorOrDie(err)
			payload.Address = listener.Addr().String()
			log.Println(payload.Address)
		}

		go func() {
			err := TunnelByProjectID(payload.ApiKey, payload.ProjectId, payload.Address, true)
			if err != nil {
				log.Println(err)
				send(Message{Text: err.Error()})
			} else {
				send(Message{Text: "Success!"})
			}
		}()

	} else {
		send(content)
	}
}

func send(msg Message) {
	byteMsg := dataToBytes(msg)
	var msgBuf bytes.Buffer
	writeMessageLength(byteMsg)
	msgBuf.Write(byteMsg)
	msgBuf.WriteTo(os.Stdout)
}

func dataToBytes(msg Message) []byte {
	byteMsg, _ := json.Marshal(msg)
	return byteMsg
}

func writeMessageLength(msg []byte) {
	binary.Write(os.Stdout, binary.LittleEndian, uint32(len(msg)))
}

func readMessageLength(msg []byte) int {
	var length uint32
	buf := bytes.NewBuffer(msg)
	binary.Read(buf, binary.LittleEndian, &length)
	return int(length)
}
