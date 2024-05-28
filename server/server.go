package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/chriswolfdesign/tcp-ttt-common/model"
	"github.com/chriswolfdesign/tcp-ttt-common/strings"
	"github.com/chriswolfdesign/tcp-ttt-common/tcp_payloads"
)

type Server struct {
	Port     string
	Listener net.Listener
	Game     *model.Game
}

type Message struct {
	Name string
	Text string
}

func GenerateServer(port string) *Server {
	return &Server{
		Port: port,
		Game: model.GenerateGame(),
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

func (s *Server) ListenForPlayerOne(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(conn)
		return
	}

	tmp := bytes.NewBuffer(buf)

	onboardingRequest := &tcp_payloads.PlayerOnboardingRequest{}
	dec := gob.NewDecoder(tmp)

	if err := dec.Decode(onboardingRequest); err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(conn)
		return
	}

	fmt.Printf("Received: %+v\n", onboardingRequest)

	if onboardingRequest.PayloadType != strings.TYPE_ONBOARDING_REQUEST {
		s.sendOnboardingFailure(conn)
		return
	}

	s.Game.SetPlayerOne(onboardingRequest.Name)

	var responseBuffer bytes.Buffer
	enc := gob.NewEncoder(&responseBuffer)

	response := tcp_payloads.GeneratePlayerOnboardingResponse(strings.ONBOARD_SUCCESS)

	if err = enc.Encode(response); err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(conn)
		return
	}

	_, err = conn.Write(responseBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(conn)
		return
	}
}

func (s *Server) sendOnboardingFailure(conn net.Conn) {
	var responseBuffer bytes.Buffer
	enc := gob.NewEncoder(&responseBuffer)

	response := tcp_payloads.GeneratePlayerOnboardingResponse(strings.ONBOARD_FAILURE)

	if err := enc.Encode(response); err != nil {
		fmt.Println(err)
		return
	}

	_, err := conn.Write(responseBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
}
