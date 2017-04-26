package model

import "aixigo/x"

// Model interface
type Model interface {
	Perform(a x.Action) x.Percept
	Update(a x.Action, e x.Percept)
	ConditionalDistribution(e x.Percept) float64
	SaveCheckpoint()
	LoadCheckpoint()
}
