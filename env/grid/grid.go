package grid

import (
	"aixigo/x"
	"fmt"
)

//Meta gives the Gridworld Metadata
var Meta x.Meta

func init() {
	Meta = x.Meta{
		ObsBits:    4,
		NumActions: 5,
		MaxReward:  dispenserReward,
		MinReward:  wallPenalty}
}

//Gridworld in the motherfucking streets
type Gridworld struct {
	Tiles [][]tile
	n     int
	pos   tile
}

//Perform an action and get a Percept back
func (gw *Gridworld) Perform(action x.Action) x.Percept {
	if action < 0 || action > Meta.NumActions {
		panic("at the Disco!")
	}
	var r x.Reward
	n, found := gw.pos.GetNeighbor(action)
	if found {
		gw.pos = n
		r = gw.pos.Rew()
	} else {
		r = wallPenalty
	}
	o := gw.pos.Obs()

	return x.Percept{O: o, R: r}
}

//Print does some bullshit Holy shit dude
func (gw *Gridworld) Print() {
	for _, row := range gw.Tiles {
		s := ""
		for _, t := range row {
			if gw.pos == t {
				s += "X"
			} else {
				switch t.(type) {
				case *empty:
					s += "0"
				case *wall:
					s += "1"
				case *dispenser:
					s += "2"
				}
			}
		}
		fmt.Println(s)
	}
}

//New Gangsta shit
func New(spec [][]int) *Gridworld {
	n := len(spec)
	tiles := [][]tile{}
	for i := 0; i < n; i++ {
		row := []tile{}
		for j := 0; j < n; j++ {
			var t tile
			switch spec[j][i] {
			case 0:
				t = &empty{baseTile{i, j, make(map[x.Action]tile)}}
			case 1:
				t = &wall{baseTile{i, j, make(map[x.Action]tile)}}
			case 2:
				t = &dispenser{baseTile{i, j, make(map[x.Action]tile)}}
			default:
				panic("Unsupported tile type")
			}
			row = append(row, t)
		}
		tiles = append(tiles, row)
	}
	// add connections
	for _, row := range tiles {
		for _, t := range row {
			x := t.X()
			y := t.Y()
			if x != 0 {
				n := tiles[x-1][y]
				switch v := n.(type) {
				case *wall:
				default:
					t.AddNeighbor(0, v) // left
				}
			}
			if x != n-1 {
				n := tiles[x+1][y]
				switch v := n.(type) {
				case *wall:
				default:
					t.AddNeighbor(1, v) // right
				}
			}
			if y != 0 {
				n := tiles[x][y-1]
				switch v := n.(type) {
				case *wall:
				default:
					t.AddNeighbor(2, v) // up
				}
			}
			if y != n-1 {
				n := tiles[x][y+1]
				switch v := n.(type) {
				case *wall:
				default:
					t.AddNeighbor(3, v) // down
				}
			}
			t.AddNeighbor(4, t)
		}
	}

	return &Gridworld{tiles, n, tiles[0][0]}
}
