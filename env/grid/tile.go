package grid

import (
	"aixigo/x"
	"math/rand"
)

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
	oBits     x.ObservationBits
	o         x.Observation
	neighbors []tile
}

func newBaseTile(i, j int) baseTile {
	return baseTile{
		x:         i,
		y:         j,
		oBits:     x.ObservationBits(make([]bool, Meta.ObsBits, Meta.ObsBits)),
		o:         -1,
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

		t.oBits[i] = true
	}
	t.o = x.ToInt(t.oBits)
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

func newDispenser(i, j int) tile {
	return &dispenser{newBaseTile(i, j), x.NewPRN(), rand.Float64()}
}

type dispenser struct {
	baseTile
	prn   *rand.Rand
	theta float64
}

func (d dispenser) Rew() x.Reward {
	if d.prn.Float64() <= d.theta {
		return dispenserReward
	}
	return emptyReward
}

type wall struct {
	baseTile
}

func (w wall) Rew() x.Reward {
	return wallPenalty
}
