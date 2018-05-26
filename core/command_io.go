package core

import (
	"net"
	"io"
	"bufio"
	"log"
	"strings"
	"fmt"
)

type CommandIO struct {
	Conn net.Conn
	Name  string
	Input chan ConnCommand
}

func StartCommandIO(conn net.Conn, name string) *CommandIO {
	commandIO := CommandIO {
		Conn: conn,
		Name:  name,
		Input: make(chan ConnCommand),
	}

	go func () {
		bufReader := bufio.NewReader(conn)
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT_SEC * time.Second))
		for {
			command := ConnCommand{}
			buf, err := bufReader.ReadBytes('\n')
			if err != nil {
				command.Error = err
				if err == io.EOF {
					command.EOF = true
				}
				commandIO.Input <- command
				break
			}
			if len(buf) == 0 {
				continue
			}

			// Skip last \n and convert to upper case
			command.Line = string(buf[:len(buf)-1])
			commandIO.Input <- command

			if LOG_TCP_RECEIVE {
				log.Println(name+" received:", command.Line)
			}
		}
	}()

	return &commandIO
}

func (c *CommandIO) SendLine(line string) {
	c.Conn.Write([]byte(line + "\n"))
	if LOG_TCP_SEND {
		log.Println(c.Name+" sent:", line)
	}
}

func (c *CommandIO) SendCommand(commandName string, args... interface{}) {
	line := GetCommandLine(commandName, args...)
	c.SendLine(line)
}

func (c *CommandIO) Unlock() {
	c.Input <- ConnCommand {}
}

func (c *CommandIO) Close() {
	c.Conn.Close()
}

func CreatePipeIO() (*CommandIO, *CommandIO) {
	server, client := net.Pipe()
	ioServer := StartCommandIO(server, "CLIENT")
	ioClient := StartCommandIO(client, "SERVER")
	return ioServer, ioClient
}

func GetCommandLine(commandName string, args... interface{}) string {
	argsStr := make([]string, 0, len(args)+1)
	argsStr = append(argsStr, strings.ToUpper(commandName))
	for _, v := range args {
		argsStr = append(argsStr, fmt.Sprintf("%v", v))
	}
	line := strings.Join(argsStr, " ")
	return line
}