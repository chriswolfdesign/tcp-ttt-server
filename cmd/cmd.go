package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Message struct {
	Name string
	Text string
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmp := bytes.NewBuffer(buf)

	message := &Message{}
	dec := gob.NewDecoder(tmp)

	if err := dec.Decode(message); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Received: %+v\n", message)

	var responseBuffer bytes.Buffer
	enc := gob.NewEncoder(&responseBuffer)

	response := Message {
		Name: message.Name,
		Text: "Message received",
	}

	if err = enc.Encode(response); err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(responseBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
}
