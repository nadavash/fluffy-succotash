package main

import (
	"fmt"
)

func main() {
	serv := NewGameServer("Fluffy Succotash")
	err := serv.Run()
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}
}
