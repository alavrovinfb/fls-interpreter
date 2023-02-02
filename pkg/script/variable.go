package script

import "sync"

// var Variables map[string]Value
var Vars *Variables

func init() {
	Vars = NewVariables()
}

type Variables struct {
	sync.RWMutex
	V map[string]Value
}

func NewVariables() *Variables {
	return &Variables{
		V: make(map[string]Value),
	}
}

func (v *Variables) GetVars() map[string]Value {
	return v.V
}

func (v *Variables) Reset() {
	v.RLock()
	defer v.RUnlock()
	v.V = make(map[string]Value)
}
