package main

import (
	"../core"
	"testing"
	"sync"
)

const CLIENT_NAME = "john"

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

	var wg sync.WaitGroup
	wg.Add(1)
	// Client side mock
	go func() {
		defer wg.Done()
		ioServer.SendCommand("START", CLIENT_NAME)
		ioServer.Close()
	}()

	// Wait synchronously until game in server side ends
	game.Play(ioClient)

	// Wait synchronously until game in client side ends
	wg.Wait()

	clientName := game.GetClientName(*ioClient)
	if clientName != CLIENT_NAME {
		t.Error("Wrong client name", clientName, "instead of", CLIENT_NAME)
	}
}
