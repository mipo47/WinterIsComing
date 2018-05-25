package main

import "../core"

type Command interface {
	Execute(g *Game, connCommand core.ConnCommand, commandIO core.CommandIO)
}
