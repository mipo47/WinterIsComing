package core

import (
	"testing"
	"time"
)

const BROADCAST_COUNT = 5

func TestBroadcast_SendCommand(t *testing.T) {
	broadcast := Broadcast {
		outputs: make([]CommandIO, BROADCAST_COUNT, BROADCAST_COUNT),
	}
	recepients := make([]CommandIO, BROADCAST_COUNT, BROADCAST_COUNT)

	for i := 0; i < BROADCAST_COUNT; i++ {
		ioServer, ioClient := CreatePipeIO()
		broadcast.outputs[i] = *ioClient
		recepients[i] = *ioServer
		defer ioServer.Close()
		defer ioClient.Close()
	}

	broadcast.SendCommand("Hello", 2, "test")

	for _, recepient := range recepients {
		com := <-recepient.Input
		if com.Error != nil {
			t.Error("Error while receiving:", com.Error)
		}
		if com.Line != "HELLO 2 test" {
			t.Error("Wrong data received:", com.Line)
		}
	}

	select {
	case com := <-recepients[0].Input:
		t.Error("Received unexpected data", com.Line)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestBroadcast_AddOutput(t *testing.T) {
	broadcast := new(Broadcast)
	if broadcast.outputs != nil || len(broadcast.outputs) != 0 {
		t.Error("Unknown outputs detected")
	}

	output := new(CommandIO)
	broadcast.AddOutput(*output)
	if broadcast.outputs == nil {
		t.Error("Broadcast outputs was not initialized")
	}

	if len(broadcast.outputs) != 1 {
		t.Error("Output was not added")
	}

	if broadcast.outputs[0] != *output {
		t.Error("Wrong output was added")
	}

	broadcast.AddOutput(*output)
	if len(broadcast.outputs) != 1 {
		t.Error("Duplicate output was added")
	}
}

func TestBroadcast_RemoveOutput(t *testing.T) {
	broadcast := new(Broadcast)
	output := new(CommandIO)
	broadcast.AddOutput(*output)

	broadcast.RemoveOutput(*output)
	if len(broadcast.outputs) != 0 {
		t.Error("Output was not removed")
	}
}
