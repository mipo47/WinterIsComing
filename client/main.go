package main

import (
	"../core"
	"./gamer"
	"log"
	"net"
	"strconv"
)

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

	gamer := gamer.CreateGamer()
	gamer.Play(*ioServer)
}
