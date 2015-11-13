package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// GameServer for handling client connections and movement into rooms.
type GameServer struct {
	name     string
	clients  []*client
	incoming chan string
	outgoing chan string
}

// NewGameServer makes a new server.
func NewGameServer(name string) *GameServer {
	gs := new(GameServer)
	gs.name = name
	gs.clients = make([]*client, 0)
	gs.incoming = make(chan string)
	gs.outgoing = make(chan string)
	return gs
}

// Broadcast sends the given data to all clients in the GameServer.
func (gs *GameServer) Broadcast(data string) {
	for _, c := range gs.clients {
		c.outgoing <- data
	}
}

// Name of the current server.
func (gs *GameServer) Name() string {
	return gs.name
}

// Run opens the GameServer to accept connections.
func (gs *GameServer) Run() error {
	fmt.Printf("Welcome to '%s' game server!\n", gs.name)

	ln, err := net.Listen("tcp", ":8008")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}

	go func() {
		for data := range gs.incoming {
			gs.Broadcast(data)
		}
	}()

	for {
		fmt.Println("Awaiting connection...")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("%s\n", err.Error())
			continue
		}

		fmt.Println("Accepting connection from remote: ", conn.RemoteAddr())
		go gs.handleConnection(conn)
	}

	fmt.Println("Goodbye!")
	return nil
}

func (gs *GameServer) handleConnection(conn net.Conn) {
	buf := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	buf.WriteString("Say my name...\n")
	err := buf.Flush()
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	pass, err := buf.ReadString('\n')
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	if strings.TrimSpace(pass) != "heisenberg" {
		conn.Close()
		return
	}

	buf.WriteString("And what's yours?\n")
	err = buf.Flush()

	if err != nil {
		fmt.Errorf(err.Error())
		conn.Close()
		return
	}

	name, err := buf.ReadString('\n')
	if err != nil {
		fmt.Errorf(err.Error())
		conn.Close()
		return
	}

	c := newClient(strings.TrimSpace(name), conn)
	gs.clients = append(gs.clients, c)
	go func() {
		for !c.closed {
			msg := <-c.incoming
			gs.incoming <- fmt.Sprintf("%s: %s", c.name, msg)
		}
	}()
	c.Listen()
}

type client struct {
	name     string
	conn     net.Conn
	closed   bool
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func newClient(name string, conn net.Conn) *client {
	return &client{
		name:     name,
		conn:     conn,
		closed:   false,
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
	}
}

func (c *client) Read() {
	for !c.closed {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(c.name, " closed the connection!")
			} else {
				fmt.Errorf("%s", err.Error())
			}
			return
		}
		c.incoming <- line
	}
}

func (c *client) Write() {
	for data := range c.outgoing {
		if c.closed {
			return
		}

		_, err := c.writer.WriteString(data)
		if err != nil {
			fmt.Errorf("%s", err.Error())
			return
		}

		c.writer.Flush()
	}
}

func (c *client) Listen() {
	go c.Read()
	go c.Write()
}
