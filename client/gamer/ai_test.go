package gamer

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"testing"
)

func TestAI_Closest_GetShootXY(t *testing.T) {
	zombies := make(map[string]core.Zombie)
	ai := new(AI_Closest)
	var x, y int

	x, y = ai.GetShootXY(zombies)
	if x != 0 || y != 0 {
		t.Error("Nothing to shoot at. (0, 0) should be returned as default value")
	}

	zombies["z1"] = core.Zombie { X: 3, Y: 2 }
	x, y = ai.GetShootXY(zombies)
	if x != 3 || y != 2 {
		t.Error("Wrong closest position: ", x, y)
	}

	zombies["z2"] = core.Zombie { X: 2, Y: 1 }
	x, y = ai.GetShootXY(zombies)
	if x != 3 || y != 2 {
		t.Error("Wrong closest position: ", x, y)
	}

	zombies["z3"] = core.Zombie { X: 4, Y: 1 }
	x, y = ai.GetShootXY(zombies)
	if x != 4 || y != 1 {
		t.Error("Wrong closest position: ", x, y)
	}
}
