package main

import (
	"../core"
	"time"
	"log"
)

type Game struct {
	isStarted bool
	width int
	height int
	zombies []Zombie
	gameResult string
	clientNames map[core.CommandIO]string
	commands map[string]GameCommand
}

func CreateGame(width, height, zombieCount int) *Game {
	maxX, maxY := width - 10, height
	if maxX < 1 { maxX = 1 }
	g := Game {
		width: width,
		height: height,
		zombies: CreateZombies(maxX, maxY, zombieCount),
	}

	g.clientNames = make(map[core.CommandIO]string)

	// Bind commands
	g.commands = make(map[string]GameCommand)
	g.commands["START"] = CommandStart{}
	g.commands["SHOOT"] = CommandShoot{}

	return &g
}

func (g *Game) SetClientName(io core.CommandIO, name string) {
	g.clientNames[io] = name
}

func (g *Game) GetClientName(io core.CommandIO) string {
	return g.clientNames[io]
}

func (g *Game) Play(io *core.CommandIO) {
	defer io.Close()

	if core.LOG_INFO {
		log.Println("Accepted connection")
	}

	var connCommand core.ConnCommand
	for g.gameResult == "" {
		connCommand = <-io.Input
		if connCommand.Error != nil {
			g.isStarted = false
			if connCommand.EOF {
				if core.LOG_INFO {
					log.Println("Client left game:", g.GetClientName(*io))
				}
				break
			} else {
				if core.LOG_ERROR {
					log.Println("Read error:", connCommand.Error)
				}
				panic(connCommand.Error)
			}
		} else {
			g.ExecuteCommand(connCommand, *io)
		}
	}

	if core.TCP_SEND_RESULT && g.gameResult != "" {
		io.SendLine(g.gameResult)
	}
}

func (g *Game) ExecuteCommand(connCommand core.ConnCommand, io core.CommandIO) {
	defer func() {
		if r := recover(); r != nil {
			if core.LOG_ERROR {
				log.Println("Recovered after unsuccessful command:", r)
			}
			if core.TCP_SEND_ERRORS {
				io.SendCommand("ERROR ", r)
			}
		}
	}()
	args := connCommand.Split()
	if len(args) > 0 {
		gameCommand := g.commands[args[0]]
		if gameCommand != nil {
			gameCommand.Execute(g, connCommand, io)
		} else {
			if core.LOG_ERROR {
				log.Println("Unknown command:", connCommand.Line)
			}
			if core.TCP_SEND_ERRORS {
				io.SendLine("ERROR Unknown command")
			}
		}
	}
}

func (g *Game) ShowZombies(io core.CommandIO)  {
	for g.isStarted && g.gameResult == "" {
		for _, zombie := range g.zombies {
			if zombie.isDead { continue }
			io.SendCommand("WALK", zombie.name, zombie.x, zombie.y)
		}
		time.Sleep(core.SHOW_ZOMBIES_MS * time.Millisecond)
	}
}

func (g *Game) MoveZombies(io core.CommandIO)  {
	for g.isStarted && g.gameResult == "" {
		for i := 0 ; i < len(g.zombies); i++ {
			zombie := &g.zombies[i]
			if !zombie.isDead {
				zombie.Move(g)
			}
		}
		time.Sleep(core.MOVE_ZOMBIES_MS * time.Millisecond)
	}
}
