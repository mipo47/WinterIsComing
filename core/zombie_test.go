package core

import "testing"

func TestCreateZombies(t *testing.T) {
	maxX, maxY, zombieCount := 20, 10, 5
	zombies := CreateZombies(maxX, maxY, zombieCount)
	if zombies == nil {
		t.Error("Zombies is not initialized")
	}
	if len(zombies) != zombieCount {
		t.Error("Wrong amount of zombies:", len(zombies), "instead of", zombieCount)
	}
	for _, zombie := range zombies {
		if zombie.X < 0 || zombie.X >= maxX || zombie.Y < 0 || zombie.Y >= maxY {
			t.Error("Zombie x is out of bound", zombie)
		}
	}
}
