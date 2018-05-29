package main

import "github.com/mysteriumnetwork/winter-server/core"

type Command interface {
	Execute(g *Game, connCommand core.ConnCommand, io core.CommandIO)
}
