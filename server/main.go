package main

import (
	"../core"
	"log"
	"net"
	"strconv"
)

func main() {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(core.TCP_PORT))
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	log.Println("Wait for connection to port", core.TCP_PORT)

	// Allows to play multiple archers with same board
	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, core.ZOMBIE_COUNT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		// Start async reading of commands
		io := core.StartCommandIO(conn, "SERVER")
		go game.Play(io)
	}
}
