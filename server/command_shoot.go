package main

import (
	"../core"
	"strconv"
	"fmt"
)

type CommandShoot struct {
	GameCommand
}

func (CommandShoot) Execute(g *Game, connCommand core.ConnCommand, io core.CommandIO)  {
	var x, y int
	var err error

	args := connCommand.Split()
	if len(args) != 3 {
		panic("Wrong SHOOT command format")
	}

	x, err = strconv.Atoi(args[1])
	if err != nil {
		panic("Cannot parse SHOOT x coordinate: " + args[1])
	}
	y, err = strconv.Atoi(args[2])
	if err != nil {
		panic("Cannot parse SHOOT y coordinate: " + args[2])
	}

	aliveZombieCount := 0
	hitZombies := make([]Zombie, 0, len(g.zombies))
	for i := 0 ; i < len(g.zombies); i++ {
		zombie := &g.zombies[i]
		if !zombie.isDead {
			if zombie.x == x && zombie.y == y {
				zombie.isDead = true
				hitZombies = append(hitZombies, *zombie)
			} else {
				aliveZombieCount++
			}
		}
	}

	clientName := g.GetClientName(io)
	reply := fmt.Sprintf("BOOM %v %d", clientName, len(hitZombies))
	for _, zombie := range hitZombies {
		reply += " " + zombie.name
	}
	io.SendLine(reply)

	// All zombies are dead
	if aliveZombieCount == 0 {
		g.gameResult = "WIN " + clientName
	}
}
