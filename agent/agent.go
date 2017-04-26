package agent

import "aixigo/x"

//Agent interface
type Agent interface {
	GetAction() x.Action
	Update(x.Action, x.Percept)
}
