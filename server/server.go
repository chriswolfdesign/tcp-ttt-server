package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Server struct {
	Port     string
	Listener net.Listener
}

type Message struct {
	Name string
	Text string
}

func GenerateServer(port string) *Server {
	return &Server{
		Port: port,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		return err
	}

	s.Listener = ln
	return nil
}

func (s *Server) Accept() (net.Conn, error) {
	return s.Listener.Accept()
}

func (s *Server) HandleConnection(conn net.Conn) {
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

	response := Message{
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
