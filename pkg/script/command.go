package script

import (
	"fmt"
	"strings"

	"github.com/alavrovinfb/fls-interpreter/pkg/messages"
)

// Operations constants
const (
	CREATE  = "create"
	DELETE  = "delete"
	UPDATE  = "update"
	ADD     = "add"
	PRINT   = "print"
	SUB     = "subtract"
	MUL     = "multiply"
	DIV     = "divide"
	VALUE   = "value"
	OPERAND = "operand"
)

// operation type either unary or binary
const (
	UNARY = iota
	BINARY
)

// tokens constants
const (
	DOLLAR = "$"
	HASH   = "#"
)

const (
	CMD = "cmd"
	ID  = "id"
)

type Command struct {
	CMD        string
	Params     map[string]interface{}
	ParentVars map[string]interface{}
}

func (c *Command) Do(out *[]interface{}) error {
	vars := Vars.GetVars()
	switch c.CMD {
	case UPDATE:
		v := c.processParams(UNARY)
		id, err := c.getID()
		if err != nil {
			return err
		}
		vars[id] = v[0]
	case ADD:
		v := c.processParams(BINARY)
		id, err := c.getID()
		if err != nil {
			return err
		}
		vars[id] = v[0] + v[1]
	case CREATE:
		v := c.processParams(UNARY)
		id, err := c.getID()
		if err != nil {
			return err
		}
		vars[id] = v[0]
	case PRINT:
		v := c.processParams(UNARY)
		fmt.Println(v[0])
		*out = append(*out, v[0])
	case DELETE:
		id, err := c.getID()
		if err != nil {
			return err
		}
		delete(vars, id)
	case SUB:
		id, err := c.getID()
		if err != nil {
			return err
		}
		v := c.processParams(BINARY)
		vars[id] = v[0] - v[1]
	case MUL:
		id, err := c.getID()
		if err != nil {
			return err
		}
		v := c.processParams(BINARY)
		vars[id] = v[0] * v[1]
	case DIV:
		id, err := c.getID()
		if err != nil {
			return err
		}
		v := c.processParams(BINARY)
		vars[id] = v[0] / v[1]
	}
	// subtract, multiply, divide
	return nil
}

func (c *Command) getID() (string, error) {
	id := c.Params[ID]
	switch typedID := id.(type) {
	case string:
		if strings.HasPrefix(typedID, DOLLAR) {
			refID := strings.TrimPrefix(typedID, DOLLAR)
			switch typedRefID := c.ParentVars[refID].(type) {
			case string:
				return typedRefID, nil
			default:
				return "", fmt.Errorf(messages.ErrIncorrectType, typedRefID, typedRefID)
			}
		} else {
			return typedID, nil
		}
	default:
		return "", fmt.Errorf("incorrect id type %T", typedID)
	}
}

func (c *Command) processParams(opType int) []Value {
	var p []Value
	if opType == BINARY {
		p = make([]Value, 2)
		for i := 0; i < 2; i++ {
			oName := fmt.Sprintf("%s%d", OPERAND, i+1)
			op, ok := c.Params[oName]
			if ok {
				switch typedOp := op.(type) {
				case float64:
					p[i] = Value(typedOp)
				case string:
					p[i] = resolveValue(c.ParentVars[strings.TrimPrefix(typedOp, DOLLAR)])
				}
			} else {
				p[i] = resolveValue(c.Params[fmt.Sprintf("%s%d", VALUE, i+1)])
			}
		}
	} else {
		p = make([]Value, 1)
		p[0] = resolveValue(c.Params[VALUE])
	}

	return p
}

func resolveValue(v interface{}) Value {
	vars := Vars.GetVars()
	switch typedVal := v.(type) {
	case float64:
		return Value(typedVal)
	case string:
		if strings.HasPrefix(typedVal, HASH) || strings.HasPrefix(typedVal, DOLLAR) {
			refVal := strings.TrimPrefix(typedVal, HASH)
			refVal = strings.TrimPrefix(refVal, DOLLAR)
			v, ok := vars[refVal]
			if !ok {
				return -1
			}

			return v
		}
	default:
		return -1
	}

	return -1
}
