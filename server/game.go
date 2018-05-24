package main

import (
	"../core"
	"time"
	"log"
)

type Game struct {
	width int
	height int
	zombies []core.Zombie

	isStarted bool
	gameResult string
	clientNames map[core.CommandIO]string
	commands map[string]GameCommand
}

func CreateGame(width, height, zombieCount int) *Game {
	g := Game {
		width: width,
		height: height,
	}

	g.clientNames = make(map[core.CommandIO]string)

	// Bind commands
	g.commands = make(map[string]GameCommand)
	g.commands["START"] = CommandStart{}
	g.commands["SHOOT"] = CommandShoot{}

	g.Restart(zombieCount)

	return &g
}

func (g *Game) SetClientName(io core.CommandIO, name string) {
	g.clientNames[io] = name
}

func (g *Game) GetClientName(io core.CommandIO) string {
	return g.clientNames[io]
}

func (g *Game) SetResult(result string, io core.CommandIO) {
	g.gameResult = result
	io.Unlock()
}

func (g *Game) Restart(zombieCount int) {
	if core.LOG_INFO {
		log.Printf("Create %d zombies", zombieCount)
	}
	g.isStarted = false
	maxX, maxY := g.width - 10, g.height
	if maxX < 1 { maxX = 1 }
	g.zombies = core.CreateZombies(maxX, maxY, zombieCount)
}

func (g *Game) Play(io *core.CommandIO) {
	defer io.Close()

	if core.LOG_INFO {
		log.Println("Accepted connection")
	}

	var connCommand core.ConnCommand
	clientName := g.GetClientName(*io)
	for g.gameResult == "" || !g.isStarted {
		connCommand = <-io.Input
		if connCommand.Error != nil {
			g.isStarted = false
			if connCommand.EOF {
				if core.LOG_INFO {
					log.Println("Client left game:", clientName)
				}
				break
			} else {
				if core.LOG_ERROR {
					log.Println("Read error:", connCommand.Error)
				}
				panic(connCommand.Error)
			}
		} else if connCommand.Line != "" {
			g.ExecuteCommand(connCommand, *io)
		}
	}

	if g.gameResult != "" {
		if core.TCP_SEND_RESULT {
			io.SendLine(g.gameResult)
		}

		delete(g.clientNames, *io)
		if len(g.clientNames) == 0 {
			g.Restart(len(g.zombies))
		}
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
			if zombie.IsDead { continue }
			io.SendCommand("WALK", zombie.Name, zombie.X, zombie.Y)
		}
		time.Sleep(core.SHOW_ZOMBIES_MS * time.Millisecond)
	}
}

func (g *Game) MoveZombies(io core.CommandIO)  {
	for g.isStarted && g.gameResult == "" {
		for i := 0 ; i < len(g.zombies); i++ {
			zombie := &g.zombies[i]
			if !zombie.IsDead {
				lose := zombie.Move(g.width, g.height)
				if lose {
					g.SetResult("LOSE " + zombie.Name, io)
					return
				}
			}
		}
		time.Sleep(core.MOVE_ZOMBIES_MS * time.Millisecond)
	}
}
