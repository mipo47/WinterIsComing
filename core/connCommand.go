package core

import (
	"strings"
)

type ConnCommand struct {
	Line  string
	Error error
	EOF   bool
}

func (c *ConnCommand) Split() []string {
	if c.Line == "" {
		return nil
	}
	parts := strings.Split(c.Line, " ")
	parts[0] = strings.ToUpper(parts[0])
	return parts
}
