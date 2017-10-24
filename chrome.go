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
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
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

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

type ManifestData struct {
	Path string
}

// func installChromeManifestRegistry(location string) (err error) {
// 	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Google\Chrome\NativeMessagingHosts\io.butterflyfx.client`, registry.QUERY_VALUE|registry.SET_VALUE)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := k.SetStringValue("xyz", "blahblah"); err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := k.Close(); err != nil {
// 		log.Fatal(err)
// 	}
// }

func InstallChromeManifest() (err error) {
	hasBeenInstalled := false
	templateString, err := Asset("data/io.butterflyfx.client.template.json")
	if err != nil {
		return err
	}
	filePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}
	data := ManifestData{Path: filePath}
	t, err := template.New("chrome-manifest").Parse(string(templateString))
	var tpl bytes.Buffer
	t.Execute(&tpl, data)
	manifestContent := tpl.Bytes()
	usr, err := user.Current()
	if err != nil {
		return err
	}
	home := usr.HomeDir
	paths := []string{
		path.Join(home, ".config", "google-chrome"),
		path.Join(home, ".config", "google-chrome-beta"),
		path.Join(home, ".config", "google-chrome-unstable"),
		path.Join(home, "Library", "Application Support", "Google", "Chrome"),
		path.Join(home, "Library", "Application Support", "Google", "Chromium"),
	}
	for _, dir := range paths {
		if fileExists(dir) {
			nativeDirectory := path.Join(dir, "NativeMessagingHosts")
			manifestFile := path.Join(nativeDirectory, "io.butterflyfx.client.json")
			if !fileExists(nativeDirectory) {
				err = os.Mkdir(nativeDirectory, 0755)
			}
			err = ioutil.WriteFile(manifestFile, manifestContent, 0755)
			if err != nil {
				return err
			}
			hasBeenInstalled = true
			log.Println(manifestFile)
		}
	}

	if !hasBeenInstalled {
		err = errors.New("Chrome manifest was not installed")
	}
	return
}
