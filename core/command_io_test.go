package core

import (
	"testing"
	"time"
)

func TestCommandIO_SendByParts(t *testing.T) {
	ioServer, ioClient := CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	ioClient.Conn.Write([]byte("HELLO 1 "))
	ioClient.Conn.Write([]byte("test\n"))

	com := <-ioServer.Input
	if com.Error != nil {
		t.Error("Error while receiving:", com.Error)
	}
	if com.Line != "HELLO 1 test" {
		t.Error("Wrong data received:", com.Line)
	}
}

func TestCommandIO_SendLine(t *testing.T) {
	ioServer, ioClient := CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	ioClient.SendLine("HELLO 1 test")

	com := <-ioServer.Input
	if com.Error != nil {
		t.Error("Error while receiving:", com.Error)
	}
	if com.Line != "HELLO 1 test" {
		t.Error("Wrong data received:", com.Line)
	}
}

func TestCommandIO_SendCommand(t *testing.T) {
	ioServer, ioClient := CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	ioClient.SendCommand("Hello", 1, "test")

	com := <-ioServer.Input
	if com.Error != nil {
		t.Error("Error while receiving:", com.Error)
	}
	if com.Line != "HELLO 1 test" {
		t.Error("Wrong data received:", com.Line)
	}

	select {
	case com = <-ioServer.Input:
		t.Error("Received unexpected data", com.Line)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestCommandIO_Close(t *testing.T) {
	ioServer, ioClient := CreatePipeIO()
	defer ioServer.Close()
	defer ioClient.Close()

	ioClient.Close()

	com := <-ioServer.Input
	if !com.EOF {
		t.Error("End of file not received:", com.Error)
	}
}