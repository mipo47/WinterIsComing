package main

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"time"
)

func (g *Game) ShowZombies(io core.CommandIO)  {
	defer g.wg.Done()
	for g.isStarted && g.gameResult == "" {
		for _, zombie := range g.zombies {
			if zombie.IsDead { continue }
			io.SendCommand("WALK", zombie.Name, zombie.X, zombie.Y)
		}
		time.Sleep(core.SHOW_ZOMBIES_MS * time.Millisecond)
	}
}

func (g *Game) MoveZombies(io core.CommandIO)  {
	defer g.wg.Done()
	for g.isStarted && g.gameResult == "" {
		for i := 0 ; i < len(g.zombies); i++ {
			zombie := &g.zombies[i]
			if !zombie.IsDead {
				lose := zombie.Move(g.width, g.height)
				if lose {
					g.SetResult("LOSE " + zombie.Name, io)
					return
				}
			}
		}
		time.Sleep(core.MOVE_ZOMBIES_MS * time.Millisecond)
	}
}
