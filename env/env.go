package env

import "aixigo/x"

// Environment interface
type Environment interface {
	Perform(action x.Action) x.Percept
}
