package grid

import "aixigo/x"

const wallPenalty = x.Reward(-10)
const dispenserReward = x.Reward(10)
const emptyReward = x.Reward(0)

type tile interface {
	X() int
	Y() int
	Rew() x.Reward
	Obs() x.Observation
	GenerateObs()
	AddNeighbor(action x.Action, t tile)
	GetNeighbor(action x.Action) tile
}

type baseTile struct {
	x         int
	y         int
	o         x.Observation
	neighbors []tile
}

func newBaseTile(i, j int) baseTile {
	return baseTile{
		x:         i,
		y:         j,
		o:         x.Observation(make([]bool, Meta.ObsBits, Meta.ObsBits)),
		neighbors: make([]tile, Meta.NumActions, Meta.NumActions),
	}
}

func (t baseTile) X() int {
	return t.x
}

func (t baseTile) Y() int {
	return t.y
}

func (t baseTile) Obs() x.Observation {
	return t.o
}

func (t baseTile) GenerateObs() {
	for i := 0; i < Meta.ObsBits; i++ {
		if t.neighbors[i] != nil {
			// not a wall
			continue
		}

		t.o[i] = true
	}
}

func (t baseTile) AddNeighbor(action x.Action, n tile) {
	t.neighbors[int(action)] = n
}

func (t baseTile) GetNeighbor(action x.Action) tile {
	return t.neighbors[int(action)]
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
