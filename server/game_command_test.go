package main

import (
	"../core"
	"testing"
	"fmt"
	"sync"
)

func TestCommandStart_Execute(t *testing.T) {
	ioServer, ioClient := core.CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, core.ZOMBIE_COUNT)
	if len(game.zombies) != core.ZOMBIE_COUNT {
		t.Error("Wrong amount of zombies", len(game.zombies), "instead of", core.ZOMBIE_COUNT)
	}

	commandStart := CommandStart{}
	command := core.ConnCommand {
		Line: fmt.Sprintf("START %v", CLIENT_NAME),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		commandStart.Execute(game, command, *ioClient)

		clientName := game.GetClientName(*ioClient)
		if clientName != CLIENT_NAME {
			t.Error("Client's name is not set properly:", clientName, "instead of", CLIENT_NAME)
		}
	}()

	for _, zombie := range game.zombies {
		command := <-ioServer.Input
		args := command.Split()
		if len(args) != 4 || args[0] != "WALK" {
			t.Error("Wrong WALK format:", command.Line)
		}
		if args[1] != zombie.name {
			t.Error("Wrong zombie name", args[1], "instead of", zombie.name)
		}
	}
}

func TestCommandShoot_Execute(t *testing.T) {
	ioServer, ioClient := core.CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, 1)
	game.SetClientName(*ioClient, CLIENT_NAME)
	zombie := game.zombies[0]

	commandShoot := CommandShoot{}
	command := core.ConnCommand {
		Line: fmt.Sprintf("SHOOT %d %d", zombie.x, zombie.y),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		commandShoot.Execute(game, command, *ioClient)
		if game.gameResult != "WIN "+CLIENT_NAME {
			t.Error("Game doesn't shows WIN result: " + game.gameResult)
		}
	}()

	reply := (<-ioServer.Input).Line
	expected := "BOOM " + CLIENT_NAME + " 1 " + zombie.name
	if reply != expected {
		t.Error("Wrong reply to SHOOT command:", reply)
		t.Error("Expected reply:", expected)
	}

	zombie = game.zombies[0]
	if !zombie.isDead {
		t.Error("Zombie should be dead after shoot")
	}
}
