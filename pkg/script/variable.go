package script

// var Variables map[string]Value
var Vars *Variables

func init() {
	Vars = NewVariables()
}

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
