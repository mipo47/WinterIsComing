package main

import (
	"../core"
	"testing"
	"time"
)

const CLIENT_NAME = "john"
const PROCESSING_TIME = time.Millisecond * 100

func TestCreateGame(t *testing.T) {
	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, core.ZOMBIE_COUNT)
	if game.width != core.BOARD_WIDTH {
		t.Error("Wrong width", game.width, "instead of", core.BOARD_WIDTH)
	}
	if game.height != core.BOARD_HEIGHT {
		t.Error("Wrong height", game.height, "instead of", core.BOARD_HEIGHT)
	}
}

func TestPlayGame(t *testing.T) {
	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, core.ZOMBIE_COUNT)

	ioServer, ioClient := core.CreatePipeIO()
	defer ioClient.Close()

	go game.Play(ioClient)

	ioServer.SendCommand("START", CLIENT_NAME)

	time.Sleep(PROCESSING_TIME) // Server need few moments to set client name

	clientName := game.GetClientName(*ioClient)
	if clientName != CLIENT_NAME {
		t.Error("Wrong client name", clientName, "instead of", CLIENT_NAME)
	}

	// Client disconnects
	ioServer.Close()

	time.Sleep(PROCESSING_TIME) // Server need few moments to erase client name

	clientName = game.GetClientName(*ioClient)
	if clientName != "" {
		t.Error("Client name not erased after game:", clientName)
	}
}
