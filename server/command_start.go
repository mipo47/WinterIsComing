package main

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"log"
)

type CommandStart struct {
	Command
}

func (CommandStart) Execute(g *Game, connCommand core.ConnCommand, io core.CommandIO)  {
	g.broadcast.SendLine(connCommand.Line)
	args := connCommand.Split()
	clientName := args[1]
	g.SetClientName(io, clientName)
	log.Println("Set client name:", clientName)
	if !g.isStarted {
		if core.LOG_INFO {
			log.Println("Starting game")
		}
		g.gameResult = ""
		g.isStarted = true
	}

	g.wg.Add(2)
	go g.ShowZombies(io)
	go g.MoveZombies(io)
}
