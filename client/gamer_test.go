package main

import (
	"../core"
	"testing"
	"time"
	"sync"
	"strings"
)

const ZOMBIE_NAME = "zombie1"

func TestCreateGamer(t *testing.T) {
	g := CreateGamer()
	if g == nil {
		t.Error("Gamer instance wasn't created (nil)")
	}

	if g.name == "" {
		t.Error("Gamer name cannot be empty")
	}

	if g.ai == nil {
		t.Error("AI for gamer wasn't selected")
	}

	if g.gameOver == true {
		t.Error("Game is not started yet, but gameOver is already set to true")
	}
}

func TestGamer_Play(t *testing.T) {
	ioServer, ioClient := core.CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	g := CreateGamer()
	g.gameOver = true

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Play(*ioServer)
	}()

	time.Sleep(20 * time.Millisecond) // give time to update gameOver field
	if g.gameOver == true {
		t.Error("gameOver field wasn't updated after game started")
	}

	startCommand := <-ioClient.Input
	if startCommand.Line != "START " + g.name {
		t.Error("START command must be sent to server")
	}

	ioClient.SendCommand("WIN", g.name)
	wg.Wait()
	if g.gameOver != true {
		t.Error("Game should be over now")
	}
}

func TestGamer_RefreshZombiePosition(t *testing.T) {
	g := CreateGamer()
	if len(g.zombies) != 0 {
		t.Error("Zombie should be added in later stage")
	}

	g.RefreshZombiePosition(strings.Split("WALK " + ZOMBIE_NAME + " 1 2", " "))
	if len(g.zombies) != 1 {
		t.Error("Exactly one zombie should be added, not", len(g.zombies))
	}

	zombie, found := g.zombies[ZOMBIE_NAME]
	if !found {
		t.Errorf("Added zombie '%v' not found", ZOMBIE_NAME)
	}
	if zombie.X != 1 || zombie.Y != 2 {
		t.Error("Incorrect zombie position:", zombie)
	}
	if zombie.Name != ZOMBIE_NAME {
		t.Errorf("Zombie should be called '%v' not", zombie.Name)
	}

	g.RefreshZombiePosition(strings.Split("WALK " + ZOMBIE_NAME + " 2 2", " "))
	if len(g.zombies) != 1 {
		t.Error("New zombie added instead of refresh")
	}
	zombie = g.zombies[ZOMBIE_NAME]
	if zombie.X != 2 || zombie.Y != 2 {
		t.Error("Zombie position wasn't updated:", zombie)
	}
}

func TestGamer_RefreshZombieState(t *testing.T) {
	g := CreateGamer()
	g.zombies[ZOMBIE_NAME] = core.Zombie {
		Name: ZOMBIE_NAME,
		X: 3,
		Y: 4,
	}
	if len(g.zombies) != 1 {
		t.Error("Only 1 zombie was added, not", len(g.zombies))
	}

	g.RefreshZombieState(strings.Split("BOOM GAMER_NAME 0", " "))
	if len(g.zombies) != 1 {
		t.Error("Zombie should stay alive after missed shoot")
	}

	// zombie was killed
	g.RefreshZombieState(strings.Split("BOOM GAMER_NAME 1 " + ZOMBIE_NAME, " "))
	if len(g.zombies) != 0 {
		t.Error("All zombies should be deleted")
	}
}
