package x

//Reward is int (for now)
type Reward int

//Observation is array of bools (for now)
type Observation []bool

//ToInt kek leh
func (o Observation) ToInt() int {
	s := 0
	for _, b := range o {
		if b {
			s++
		}
	}
	return s
}

//Percept is composite
type Percept struct {
	O Observation
	R Reward
}

//Action is int (for now)
type Action int

//Meta object
type Meta struct {
	ObsBits    int
	NumActions Action
	MaxReward  Reward
	MinReward  Reward
}

// Model interface
type Model interface {
	Perform(a Action) Percept
	Update(a Action, e Percept)
	ConditionalDistribution(e Percept) float64
	SaveCheckpoint()
	LoadCheckpoint()
	Copy() Model
}

//Utility function type
type Utility func(e Percept, dfr int) float64

//Environment ...
type Environment interface {
	Perform(action Action) Percept
}

//Agent ...
type Agent interface {
	GetAction() Action
	Update(Action, Percept)
}
