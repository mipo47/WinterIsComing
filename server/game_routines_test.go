package main

import (
	"../core"
	"testing"
	"time"
)

func TestGame_ShowZombies(t *testing.T) {
	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, core.ZOMBIE_COUNT)

	ioServer, ioClient := core.CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	game.isStarted = true

	game.wg.Add(1)
	go game.ShowZombies(*ioClient)

	for i := 0; i < len(game.zombies); i++ {
		command := <-ioServer.Input
		args := command.Split()
		if len(args) != 4 || args[0] != "WALK" {
			t.Error("Wrong WALK format:", command.Line)
		}
	}

	game.isStarted = false
	game.wg.Wait()
}

func TestGame_MoveZombies(t *testing.T) {
	game := CreateGame(core.BOARD_WIDTH, core.BOARD_HEIGHT, core.ZOMBIE_COUNT)

	ioServer, ioClient := core.CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	zombie := &game.zombies[0]
	x, y := zombie.X, zombie.Y

	game.isStarted = true

	game.wg.Add(1)
	go game.MoveZombies(*ioClient)

	zombieMoved := make(chan bool)
	go func() {
		for zombie.X == x && zombie.Y == y {
			time.Sleep(core.MOVE_ZOMBIES_MS * time.Millisecond)
		}
		zombieMoved <- true
	}()

	select {
	case <-zombieMoved:
		break
	case <-time.After(5 * core.MOVE_ZOMBIES_MS * time.Millisecond):
		t.Error("Zombie didn't moved")
	}

	game.isStarted = false
	game.wg.Wait()
}
