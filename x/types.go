package x

//Reward is int (for now)
type Reward int

//Observation is int (for now)
type Observation int

//ObservationBits bit format
type ObservationBits []bool

//Percept is a simple (o,r) composite
type Percept struct {
	O Observation
	R Reward
}

//Action is int (for now)
type Action int

//Meta object contains environment metadata
type Meta struct {
	ObsBits     int
	NumActions  Action
	MaxReward   Reward
	MinReward   Reward
	RewardRange float64
}

// Model interface
type Model interface {
	Perform(a Action) *Percept                  // Implements the Environment interface
	Update(a Action, e *Percept)                // Must be updateable
	ConditionalDistribution(e *Percept) float64 // Must be probabilistic
	SaveCheckpoint()                            // Save and Load are needed to reset MCTS simulations
	LoadCheckpoint()                            //
	Copy() Model                                // Need to be easy to copy
}

//Utility function signature
type Utility func(e *Percept, dfr int) float64

//Environment interface
type Environment interface {
	Perform(action Action) *Percept
}

//Agent interface
type Agent interface {
	GetAction() Action
	Update(Action, *Percept)
}
