package core

import (
	"math/rand"
	"strconv"
)

type Zombie struct {
	Name   string
	X      int
	Y      int
	IsDead bool // Dead zombie = zombie that can't move anymore
}

// Returns true if zombie passes the wall (client lose)
func (zombie *Zombie) Move(width, height int) bool {
	if zombie.IsDead {
		return false
	}
	// move to random direction
	rnd := rand.Float32()
	if rnd < 0.5 { // go right (to the wall)
		zombie.X++
		if zombie.X >= width {
			return true
		}
	} else if rnd < 0.7 { // go down
		zombie.Y++
		if zombie.Y >= height {
			zombie.Y--
		}
	} else if rnd < 0.9 { // go up
		zombie.Y--
		if zombie.Y < 0 {
			zombie.Y++
		}
	} else if zombie.X > 0 { // go left (step back)
		zombie.X--
	}
	return false
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
			for _, zombie := range (zombies) {
				if zombie.X == x && zombie.Y == y {
					unique = false
					break
				}
			}
		}

		zombies[i] = Zombie{
			Name:   "Zombie" + strconv.Itoa(i+1),
			X:      x,
			Y:      y,
			IsDead: false,
		}
	}
	return zombies
}
