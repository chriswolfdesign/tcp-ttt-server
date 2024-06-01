package main

import (
	"fmt"
	"tcp-ttt-server/server"
	"time"

	"github.com/chriswolfdesign/tcp-ttt-common/strings"
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

	fmt.Println("Server started")

	for serve.Game.PlayerOneName == "" {
		fmt.Println("Waiting for player one to join.")

		conn, err := serve.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		serve.ListenForPlayerOne(conn)
	}

	fmt.Println("Player One:", serve.Game.PlayerOneName)

	for serve.Game.PlayerTwoName == "" {
		fmt.Println("Waiting for player two to join.")

		conn, err := serve.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		serve.ListenForPlayerTwo(conn)
	}

	fmt.Println("Player Two:", serve.Game.PlayerTwoName)

	serve.InformGameStarted()

	fmt.Println("Game has begun")

	fmt.Println("Current board state")
	serve.Game.Board.PrintBoard()

	for serve.Game.Winner == strings.NOT_OVER {
		// FIXME: must wait to give player two time to prepare to receive
		time.Sleep(time.Millisecond * 250)
		serve.SendPlayerTurn()

		fmt.Println("sent messages")

		// TODO: Implement the rest of the game
		return
	}
}
