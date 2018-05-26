package core

type CommandOutput interface {
	SendLine(line string)
	SendCommand(commandName string, args... interface{})
}

// Used to send data to multiple clients
type Broadcast struct {
	outputs []CommandIO
}

func (b *Broadcast) SendLine(line string) {
	if b.outputs == nil {
		return // nothing to broadcast
	}
	for _, commandIO := range b.outputs {
		commandIO.SendLine(line)
	}
}

func (b *Broadcast) SendCommand(commandName string, args... interface{}) {
	line := GetCommandLine(commandName, args...)
	b.SendLine(line)
}

func (b *Broadcast) AddOutput(output CommandIO) {
	if b.outputs == nil {
		b.outputs = make([]CommandIO, 0, 10)
	} else {
		// I know it's not effective search, but array will be always small
		for _, o := range b.outputs {
			if o == output {
				return // skip duplicate
			}
		}
	}
	b.outputs = append(b.outputs, output)
}

func (b *Broadcast) RemoveOutput(output CommandIO) bool {
	if b.outputs == nil {
		return false
	}
	for i, o := range b.outputs {
		if o == output {
			b.outputs = append(b.outputs[:i], b.outputs[i+1:]...)
			return true
		}
	}
	return false
}
