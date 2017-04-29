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
	Perform(Action) (Observation, Reward)                // Implements the Environment interface
	Update(Action, Observation, Reward)                  // Must be updateable
	ConditionalDistribution(Observation, Reward) float64 // Must be probabilistic
	SaveCheckpoint()                                     // Save and Load are needed to reset MCTS simulations
	LoadCheckpoint()                                     //
	Copy() Model                                         // Need to be easy to copy
}

//Utility function signature
type Utility func(o Observation, r Reward, dfr int) float64

//Environment interface
type Environment interface {
	Perform(action Action) (Observation, Reward)
}

//Agent interface
type Agent interface {
	GetAction() Action
	Update(Action, Observation, Reward)
}
