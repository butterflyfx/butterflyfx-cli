package main

import "net"

func GetRandomAddress() (*net.TCPAddr, error) {
	laddr, err := net.ResolveTCPAddr("tcp", "localhost:0") // Port == 0 - free
	if err != nil {
		return laddr, err
	}
	l, err := net.ListenTCP("tcp", laddr)
	l.Close()
	return l.Addr().(*net.TCPAddr), err
}
