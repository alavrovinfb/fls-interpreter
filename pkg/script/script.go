package script

import (
	"fmt"
	"strings"

	"github.com/alavrovinfb/fls-interpreter/pkg/messages"
)

var Body *Script

func init() {
	Body = NewScript().WithOutPut()
}

type Script struct {
	Funcs map[string]*Func
	Out   *[]interface{}
}

func NewScript() *Script {
	return &Script{
		Funcs: make(map[string]*Func, 0),
	}
}

func (s *Script) WithOutPut() *Script {
	tmp := make([]interface{}, 0)
	s.Out = &tmp
	return s
}

func (s *Script) RestOut() {
	s.WithOutPut()
}

func (sc *Script) Execute(fName string, localVars map[string]interface{}) error {
	f, ok := sc.Funcs[fName]
	if !ok {
		return fmt.Errorf(messages.ErrFuncMissed, fName)
	}
	for _, c := range f.Commands {
		if strings.HasPrefix(c.CMD, HASH) {
			cRef := strings.TrimPrefix(c.CMD, HASH)
			localVars = c.Params
			if err := sc.Execute(cRef, localVars); err != nil {
				return err
			}
		} else {
			c.ParentVars = localVars
			if err := c.Do(sc.Out); err != nil {
				return fmt.Errorf(messages.ErrExecution, err)
			}
		}
	}

	return nil
}

func Parse(doc map[string]interface{}, vars *Variables, fns *Script) error {
	if vars == nil {
		vars = NewVariables()
	}
	if fns == nil {
		fns = NewScript().WithOutPut()
	}
	if _, ok := doc[InitFunc]; !ok {
		return fmt.Errorf(messages.ErrInitMissed)
	}
	for n, v := range doc {
		if v == nil {
			return fmt.Errorf(messages.ErrBodyEmpty, v)
		}
		switch typedVal := v.(type) {
		case float64:
			vars.V[n] = Value(typedVal)
		case []interface{}:
			if len(typedVal) == 0 {
				return fmt.Errorf(messages.ErrBodyEmpty, v)
			}
			c, err := ProcessFunc(n, typedVal)
			if err != nil {
				return err
			}
			fns.Funcs[n] = NewFunc(n).Add(c)
		default:
			return fmt.Errorf(messages.ErrIncorrectType, typedVal, typedVal)
		}
	}

	return nil
}

func ProcessFunc(fnName string, rawCmds []interface{}) ([]Command, error) {
	if len(rawCmds) == 0 {
		return nil, fmt.Errorf(messages.ErrBodyEmpty, fnName)
	}
	cmds := make([]Command, 0)
	for i, c := range rawCmds {
		pMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(messages.ErrCast, fnName)
		}
		if _, ok := pMap[CMD]; !ok {
			return nil, fmt.Errorf(messages.ErrCmdMissed, i, fnName)
		}
		cmd := Command{
			Params: map[string]interface{}{},
		}

		for pk, pv := range pMap {
			switch pk {
			case CMD:
				cmd.CMD = pv.(string)
			default:
				cmd.Params[pk] = pv
			}
		}
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}
