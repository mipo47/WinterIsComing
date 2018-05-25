package main

import (
	"../core"
	"strconv"
	"os"
	"log"
	"fmt"
	"time"
	"sync"
)

type Gamer struct {
	name string
	zombies map[string]core.Zombie
	gameOver bool
	ai AI
}

func CreateGamer() *Gamer {
	gamer := Gamer {
		name: "Player" + strconv.Itoa(os.Getpid()),
		zombies : make(map[string]core.Zombie),
		ai: new(AI_Closest),
	}
	return &gamer
}

func (g *Gamer) Play(ioServer core.CommandIO)  {
	g.gameOver = false
	ioServer.SendCommand("START", g.name)

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		wg.Done()
		for !g.gameOver {
			time.Sleep(core.SHOOT_SPEED_MS * time.Millisecond)
			x, y := g.ai.GetShootXY(g.zombies)
			ioServer.SendCommand("SHOOT", x, y)
		}
	}()

	for !g.gameOver {
		connCommand := <-ioServer.Input
		if connCommand.Error != nil {
			if !connCommand.EOF {
				log.Fatalln("Connection to server is broked", connCommand.Error)
			} else {
				log.Println("Server closed connection")
			}
			break
		}
		args := connCommand.Split()
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "WALK":
			g.RefreshZombiePosition(args)
		case "BOOM":
			g.RefreshZombieState(args)
		case "WIN":
			if args[1] == g.name {
				fmt.Println("You win")
			} else {
				fmt.Println("Your team wins")
			}
			g.gameOver = true
		case "LOSE":
			fmt.Println("You lose")
			g.gameOver = true
		}
	}
}

func  (g *Gamer) RefreshZombiePosition(args []string) {
	var x, y int
	var err error

	zombieName := args[1]
	x, err = strconv.Atoi(args[2])
	if err != nil {
		panic("Cannot parse WALK x coordinate: " + args[1])
	}
	y, err = strconv.Atoi(args[3])
	if err != nil {
		panic("Cannot parse WALK y coordinate: " + args[2])
	}
	g.zombies[zombieName] = core.Zombie {
		Name: zombieName,
		X: x,
		Y: y,
	}
}

func  (g *Gamer) RefreshZombieState(args []string) {
	var hitCount int
	var err error
	hitCount, err = strconv.Atoi(args[2])
	if err != nil {
		panic("Cannot parse BOOM hit count: " + args[2])
	}
	for i := 0; i < hitCount; i++ {
		zombieName := args[3+i]
		if core.LOG_INFO {
			log.Println("Deleting killed zombie:", zombieName)
		}
		delete(g.zombies, zombieName)
	}
}
