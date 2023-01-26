package script

import (
	"fmt"
	"log"
	"strings"
)

var Body *Script

type Script struct {
	Funcs map[string]*Func
}

func NewScript() *Script {
	return &Script{
		Funcs: make(map[string]*Func, 0),
	}
}

func (sc *Script) Execute(fName string, localVars map[string]interface{}) {
	f := sc.Funcs[fName]
	for _, c := range f.Commands {
		if strings.HasPrefix(c.CMD, HASH) {
			cRef := strings.TrimPrefix(c.CMD, HASH)
			cmd := sc.Funcs[cRef]
			fmt.Println("Referred cmd", cmd)
			localVars = c.Params
			sc.Execute(cRef, localVars)
		} else {
			fmt.Println("Normal cmd", c.CMD)
			c.ParentVars = localVars
			if err := c.Do(); err != nil {
				log.Println("error command executing")
			}
		}
	}
}

func Parse(doc map[string]interface{}) (*Variables, *Script) {
	vars := NewVariables() //make(map[string]Value)
	fns := NewScript()     //make(map[string]Func, 0)
	for n, v := range doc {
		if v == nil {
			log.Fatal("function body is nil")
		}
		switch typedVal := v.(type) {
		case float64:
			vars.V[n] = Value(typedVal)
		case []interface{}:
			c, err := ProcessFunc(n, typedVal)
			if err != nil {
				log.Println(err)
				continue
			}
			fns.Funcs[n] = NewFunc(n).Add(c)
		default:
			log.Printf("unknown type %s", n)
		}
	}

	return vars, fns
}

func ProcessFunc(fnName string, rawCmds []interface{}) ([]Command, error) {
	if len(rawCmds) == 0 {
		return nil, fmt.Errorf("%s function body is empty", fnName)
	}
	cmds := make([]Command, 0)
	for i, c := range rawCmds {
		pMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s can't cast cmd params", fnName)
		}
		if _, ok := pMap[CMD]; !ok {
			return nil, fmt.Errorf("comand %d in function %s has incorrect format 'cmd' key is missed", i, fnName)
		}
		cmd := Command{
			Params: map[string]interface{}{},
		}

		for pk, pv := range pMap {
			switch pk {
			case CMD:
				cmd.CMD = pv.(string)
				//if strings.HasPrefix(cmd.CMD, "#") {
				//
				//}
			//case "id":
			//	cmd.ID = pv.(string)
			default:
				cmd.Params[pk] = pv
			}
		}
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}
