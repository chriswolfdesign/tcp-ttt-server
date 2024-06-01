package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/chriswolfdesign/tcp-ttt-common/enums"
	"github.com/chriswolfdesign/tcp-ttt-common/model"
	"github.com/chriswolfdesign/tcp-ttt-common/strings"
	"github.com/chriswolfdesign/tcp-ttt-common/tcp_payloads"
)

type Server struct {
	Port          string
	Listener      net.Listener
	Game          *model.Game
	PlayerOneConn net.Conn
	PlayerTwoConn net.Conn
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

	s.PlayerOneConn = conn
	s.Game.SetPlayerOne(onboardingRequest.Name)

	var responseBuffer bytes.Buffer
	enc := gob.NewEncoder(&responseBuffer)

	response := tcp_payloads.GeneratePlayerOnboardingResponse(strings.ONBOARD_SUCCESS, enums.PLAYER_ONE)

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

	response := tcp_payloads.GeneratePlayerOnboardingResponse(strings.ONBOARD_FAILURE, "")

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

func (s *Server) ListenForPlayerTwo(conn net.Conn) {
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

	s.PlayerTwoConn = conn
	s.Game.SetPlayerTwo(onboardingRequest.Name)

	var responseBuffer bytes.Buffer
	enc := gob.NewEncoder(&responseBuffer)

	response := tcp_payloads.GeneratePlayerOnboardingResponse(strings.ONBOARD_SUCCESS, enums.PLAYER_TWO)

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

func (s *Server) InformGameStarted() {
	var gameStartedBuffer bytes.Buffer
	enc := gob.NewEncoder(&gameStartedBuffer)

	gameStartedMessage := tcp_payloads.GameStartingMessage{
		Message:     strings.GAME_STARTING_MESSAGE,
		PayloadType: strings.TYPE_GAME_STARTING_MESSAGE,
		Game:        *s.Game,
	}

	if err := enc.Encode(gameStartedMessage); err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(s.PlayerOneConn)
		s.sendOnboardingFailure(s.PlayerTwoConn)
		return
	}

	_, err := s.PlayerOneConn.Write(gameStartedBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(s.PlayerOneConn)
		return
	}

	_, err = s.PlayerTwoConn.Write(gameStartedBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(s.PlayerTwoConn)
		return
	}
}

func (s *Server) SendPlayerTurn() {
	var playerTurnBuffer bytes.Buffer
	enc := gob.NewEncoder(&playerTurnBuffer)

	playerTurnMessage := tcp_payloads.PlayerTurnMessage{
		Player: s.Game.CurrentPlayer,
		PayloadType: strings.TYPE_PLAYER_TURN_MESSAGE,
	}

	if err := enc.Encode(playerTurnMessage); err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(s.PlayerOneConn)
		s.sendOnboardingFailure(s.PlayerTwoConn)
		return
	}

	_, err := s.PlayerOneConn.Write(playerTurnBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(s.PlayerOneConn)
		return
	}

	fmt.Println("Sent player turn message to player one")

	_, err = s.PlayerTwoConn.Write(playerTurnBuffer.Bytes())
	if err != nil {
		fmt.Println(err)
		s.sendOnboardingFailure(s.PlayerTwoConn)
	}

	fmt.Println("Sent player turn message to player two")
}