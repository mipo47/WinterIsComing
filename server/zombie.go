package main

import (
	"math/rand"
	"strconv"
)

type Zombie struct {
	name string
	x int
	y int
	isDead bool // Dead zombie = zombie that can't move anymore
}

func (zombie *Zombie) Move(g *Game) {
	if zombie.isDead {
		return
	}
	// move to random direction
	rnd := rand.Float32()
	if rnd < 0.5 { // go right (to the wall)
		zombie.x++
		if zombie.x >= g.width {
			g.gameResult = "LOSE " + zombie.name
		}
	} else if rnd < 0.7 { // go down
		zombie.y++
		if zombie.y >= g.height {
			zombie.y--
		}
	} else if rnd < 0.9 { // go up
		zombie.y--
		if zombie.y < 0 {
			zombie.y++
		}
	} else if zombie.x > 0 { // go left (step back)
		zombie.x--
	}
}

func CreateZombies(maxX, maxY, zombieCount int) []Zombie {
	zombies := make([]Zombie, zombieCount, zombieCount)
	for i := 0; i < len(zombies); i++ {
		var x, y int

		// Find unique coordinates
		for unique := false; !unique; {
			x = rand.Intn(maxX)
			y = rand.Intn(maxY)
			unique = true
			for _, zombie := range(zombies) {
				if zombie.x == x && zombie.y == y {
					unique = false
					break
				}
			}
		}

		zombies[i] = Zombie {
			name: "Zombie" + strconv.Itoa(i+1),
			x: x,
			y: y,
			isDead: false,
		}
	}
	return zombies
}
