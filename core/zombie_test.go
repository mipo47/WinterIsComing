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

func TestZombie_Move(t *testing.T) {
	maxX, maxY, zombieCount := 10, 5, 1
	zombie := CreateZombies(maxX, maxY, zombieCount)[0]

	for i := 0; i < 10; i++ {
		x, y := zombie.X, zombie.Y
		reachesWall := zombie.Move(maxX, maxY)
		if reachesWall {
			if zombie.X < maxX {
				t.Error("Zombie doesn't rich wall yet", zombie)
			}
			break
		} else {
			if zombie.X < 0 || zombie.X >= maxX || zombie.Y < 0 || zombie.Y >= maxY {
				t.Error("Zombie moves out of bound", zombie)
			}

			// zombie stays in same place
			if zombie.X == x && zombie.Y == y {
				// is it board border
				if x != 0 && y != 0 && y != maxY-1 {
					t.Error("Zombie should move from current position", zombie)
				}
			}
		}
	}
}