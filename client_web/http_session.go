package main

import (
	"../client/gamer"
	"time"
)

type HttpSession struct {
	gamer *gamer.Gamer
	commands chan string
	newCommands []string
}

func CreateHttpSession(ai gamer.AI) *HttpSession {
	commands := make(chan string)
	sess := HttpSession {
		gamer: gamer.CreateCustomGamer(new(gamer.AI_Closest), commands),
		commands: commands,
		newCommands: make([]string, 0, 1000),
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
