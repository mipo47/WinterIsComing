package gamer

import (
	"github.com/mysteriumnetwork/winter-server/core"
	"sort"
)

type AI interface {
	GetShootXY(zombies map[string]core.Zombie) (int, int)
}

type AI_Closest struct { AI }

// Implement zombie sorting by distance to the wall (from closest)
type zombieByDistance []core.Zombie
func (d zombieByDistance) Len() int {
	return len(d)
}
func (d zombieByDistance) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
func (d zombieByDistance) Less(i, j int) bool {
	return d[i].X > d[j].X
}

func (ai *AI_Closest) GetShootXY(zombies map[string]core.Zombie) (int, int) {
	zombieList := make([]core.Zombie, 0, len(zombies))
	for _, zombie := range zombies {
		zombieList = append(zombieList, zombie)
	}
	if len(zombieList) == 0 {
		return 0, 0
	}

	sort.Sort(zombieByDistance(zombieList))
	closestZombie := zombieList[0]

	return closestZombie.X, closestZombie.Y
}
