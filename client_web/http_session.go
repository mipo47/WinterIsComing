package main

import (
	"../core"
	"../client/gamer"
	"time"
)

type HttpSession struct {
	gamer *gamer.Gamer
	commands chan string
	newCommands []string
	io core.CommandIO
}

func CreateHttpSession(ai gamer.AI, io core.CommandIO) *HttpSession {
	commands := make(chan string)
	sess := HttpSession {
		gamer: gamer.CreateCustomGamer(ai, commands),
		commands: commands,
		newCommands: make([]string, 0, 1000),
		io: io,
	}
	return &sess
}

func (s *HttpSession) TrackCommands()  {
	for !s.gamer.GameOver {
		select {
		case command := <-s.commands:
			s.newCommands = append(s.newCommands, command)
		case <-time.After(1 * time.Second):
		}
	}
}
