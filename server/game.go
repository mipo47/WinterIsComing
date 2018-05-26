package main

import (
	"../core"
	"log"
	"sync"
)

type Game struct {
	width int
	height int
	zombies []core.Zombie

	isStarted   bool
	gameResult  string
	clientNames map[core.CommandIO]string
	broadcast   core.Broadcast
	commands    map[string]Command
	wg          sync.WaitGroup
}

func CreateGame(width, height, zombieCount int) *Game {
	g := Game {
		width: width,
		height: height,
	}

	g.clientNames = make(map[core.CommandIO]string)
	g.broadcast = *new(core.Broadcast)

	// Bind commands
	g.commands = make(map[string]Command)
	g.commands["START"] = CommandStart{}
	g.commands["SHOOT"] = CommandShoot{}

	g.Restart(zombieCount)

	return &g
}

func (g *Game) SetClientName(io core.CommandIO, name string) {
	g.clientNames[io] = name
	g.broadcast.AddOutput(io)
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
	maxX, maxY := g.width - 20, g.height
	if maxX < 1 { maxX = 1 }
	g.zombies = core.CreateZombies(maxX, maxY, zombieCount)
}

func (g *Game) Play(io *core.CommandIO) {
	defer func() {
		if r := recover(); r != nil {
			g.isStarted = false
			if core.LOG_ERROR {
				log.Println("Disconnect after error:", r)
			}
		}
	}()
	defer io.Close()
	defer g.wg.Wait()

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
			g.executeCommand(connCommand, *io)
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

	g.broadcast.RemoveOutput(*io)
	delete(g.clientNames, *io)
}

func (g *Game) executeCommand(connCommand core.ConnCommand, io core.CommandIO) {
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
