package main

import (
	"fmt"
	"tcp-ttt-server/server"
)

type Message struct {
	Name string
	Text string
}

func main() {
	serve := server.GenerateServer(":8080")
	err := serve.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := serve.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go serve.HandleConnection(conn)
	}
}
