package grid

import "aixigo/x"

const wallPenalty = x.Reward(-10)
const dispenserReward = x.Reward(10)
const emptyReward = x.Reward(0)

type tile interface {
	Rew() x.Reward
	Obs() x.Observation
	X() int
	Y() int
	AddNeighbor(action x.Action, t tile)
	GetNeighbor(action x.Action) (tile, bool)
}

type baseTile struct {
	x         int
	y         int
	neighbors map[x.Action]tile
}

func (t baseTile) Obs() x.Observation {
	o := x.Observation{false, false, false, false}
	for i := 0; i < Meta.ObsBits; i++ {
		_, found := t.neighbors[x.Action(i)] // TODO: ow
		if found {
			// not a wall
			continue
		}

		o[i] = true
	}
	return o
}

func (t baseTile) X() int {
	return t.x
}

func (t baseTile) Y() int {
	return t.y
}

func (t baseTile) AddNeighbor(action x.Action, n tile) {
	t.neighbors[action] = n
}

func (t baseTile) GetNeighbor(action x.Action) (tile, bool) {
	n, found := t.neighbors[action]
	return n, found
}

type empty struct {
	baseTile
}

func (e empty) Rew() x.Reward {
	return emptyReward
}

type dispenser struct {
	baseTile
}

func (d dispenser) Rew() x.Reward {
	return dispenserReward
}

type wall struct {
	baseTile
}

func (w wall) Rew() x.Reward {
	return wallPenalty
}
