package main

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
		if zombie.x < 0 || zombie.x >= maxX || zombie.y < 0 || zombie.y >= maxY {
			t.Error("Zombie x is out of bound", zombie)
		}
	}
}
