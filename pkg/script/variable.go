package script

// var Variables map[string]Value
var Vars *Variables

type Variables struct {
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
