package script

const InitFunc = "init"

type Func struct {
	Name     string
	Commands []Command
}

func NewFunc(name string) *Func {
	return &Func{
		Name: name,
	}
}

func (f *Func) Add(commands []Command) *Func {
	f.Commands = commands
	return f
}
