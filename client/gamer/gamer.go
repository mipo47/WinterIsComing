package gamer

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"strconv"
	"os"
	"log"
	"fmt"
	"time"
	"sync"
)

type Gamer struct {
	Name     string
	Zombies  map[string]core.Zombie
	GameOver bool
	ai       AI
	commands chan string
}

func CreateGamer() *Gamer {
	return CreateCustomGamer(new(AI_Closest), nil)
}

func CreateCustomGamer(ai AI, commands chan string) *Gamer {
	gamer := Gamer {
		Name:     "Player" + strconv.Itoa(os.Getpid()),
		Zombies:  make(map[string]core.Zombie),
		ai:       ai,
		GameOver: false,
		commands: commands,
	}
	return &gamer
}

func (g *Gamer) Play(ioServer core.CommandIO)  {
	g.GameOver = false
	ioServer.SendCommand("START", g.Name)

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		wg.Done()
		for !g.GameOver {
			time.Sleep(core.SHOOT_SPEED_MS * time.Millisecond)
			if g.ai != nil {
				x, y := g.ai.GetShootXY(g.Zombies)
				ioServer.SendCommand("SHOOT", x, y)
			}
		}
	}()

	for !g.GameOver {
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

		if g.commands != nil {
			g.commands <- connCommand.Line
		}

		switch args[0] {
		case "WALK":
			g.RefreshZombiePosition(args)
		case "BOOM":
			g.RefreshZombieState(args)
		case "WIN":
			if args[1] == g.Name {
				fmt.Println("You win")
			} else {
				fmt.Println("Your team wins")
			}
			g.GameOver = true
		case "LOSE":
			fmt.Println("You lose")
			g.GameOver = true
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
	g.Zombies[zombieName] = core.Zombie {
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
		delete(g.Zombies, zombieName)
	}
}
