package main

import "../core"

type GameCommand interface {
	Execute(g *Game, connCommand core.ConnCommand, commandIO core.CommandIO)
}
