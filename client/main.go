package main

import (
	"../core"
	"log"
	"net"
	"strconv"
	"fmt"
	"time"
	"math/rand"
)

const (
	CLIENT_NAME = "Miroslav"
	SHOOT_SPEED_MS = 1500
)

func refreshZombiePosition(zombies map[string]core.Zombie, args []string) {
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
	zombies[zombieName] = core.Zombie {
		Name: zombieName,
		X: x,
		Y: y,
	}
}

func refreshZombieState(zombies map[string]core.Zombie, args []string) {
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
		delete(zombies, zombieName)
	}
}

func main()  {
	if core.LOG_INFO {
		log.Println("Connecting to localhost:", core.TCP_PORT)
	}

	conn, err := net.Dial("tcp", "localhost:" + strconv.Itoa(core.TCP_PORT))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	ioServer := core.StartCommandIO(conn, "CLIENT")
	defer ioServer.Close()

	zombies := make(map[string]core.Zombie)
	ioServer.SendCommand("START", CLIENT_NAME)

	gameOver := false
	go func () {
		for !gameOver {
			time.Sleep(SHOOT_SPEED_MS * time.Millisecond)

			zombieNames := make([]string, 0, len(zombies))
			for k := range zombies {
				zombieNames = append(zombieNames, k)
			}

			randomName := zombieNames[rand.Intn(len(zombieNames))]
			aimToZombie := zombies[randomName]

			ioServer.SendCommand("SHOOT", aimToZombie.X, aimToZombie.Y)
		}
	}()

	for !gameOver {
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
			refreshZombiePosition(zombies, args)
		case "BOOM":
			refreshZombieState(zombies, args)
		case "WIN":
			if args[1] == CLIENT_NAME {
				fmt.Println("You win")
			} else {
				fmt.Println("Your team wins")
			}
			gameOver = true
		case "LOSE":
			fmt.Println("You lose")
			gameOver = true
		}
	}
}
