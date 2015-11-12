package main

import (
	// "bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	conn.Write([]byte("Hi there!\n"))
	conn.Close()
}

func main() {
	fmt.Println("Welcome to Fluffy Succotash Game Server!")

	ln, err := net.Listen("tcp", ":8008")
	if err != nil {
		panic(err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf(err.Error())
		}

		go handleConnection(conn)
	}
}
