package script

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testScript = `
{
  "var1":1,
  "var2":2,
  "var5":5,  
  
  "init": [
    {"cmd" : "#setup" }
  ],
  
  "setup": [
    {"cmd":"update", "id": "var1", "value":3.5},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"#sum", "id": "var1", "value1":"#var1", "value2":"#var2"},
    {"cmd":"print", "value": "#var1"},
    {"cmd":"create", "id": "var3", "value":5},
    {"cmd":"delete", "id": "var1"},
    {"cmd":"#printAll"}
  ],
  
  "sum": [
      {"cmd":"add", "id": "$id", "operand1":"$value1", "operand2":"$value2"}
  ],

  "printAll":
  [
    {"cmd":"print", "value": "#var1"},
    {"cmd":"print", "value": "#var2"},
    {"cmd":"print", "value": "#var3"}
  ]
}
`
var resVars = Variables{
	V: map[string]Value{
		"var1": 1, "var2": 2, "var5": 5,
	},
}
var resFuncs = Script{
	Funcs: map[string]*Func{
		"init": {
			Name: "init",
			Commands: []Command{
				{
					CMD:        "#setup",
					Params:     map[string]interface{}{},
					ParentVars: nil,
				},
			},
		},
		"setup": {
			Name: "setup",
			Commands: []Command{
				{
					CMD:        "update",
					Params:     map[string]interface{}{"id": "var1", "value": 3.5},
					ParentVars: nil,
				},
				{
					CMD:        "print",
					Params:     map[string]interface{}{"value": "#var1"},
					ParentVars: nil,
				},
				{
					CMD:        "#sum",
					Params:     map[string]interface{}{"id": "var1", "value1": "#var1", "value2": "#var2"},
					ParentVars: nil,
				},
				{
					CMD:        "print",
					Params:     map[string]interface{}{"value": "#var1"},
					ParentVars: nil,
				},
				{
					CMD:        "create",
					Params:     map[string]interface{}{"id": "var3", "value": float64(5)},
					ParentVars: nil,
				},
				{
					CMD:        "delete",
					Params:     map[string]interface{}{"id": "var1"},
					ParentVars: nil,
				},
				{
					CMD:        "#printAll",
					Params:     map[string]interface{}{},
					ParentVars: nil,
				},
			},
		},
		"sum": {
			Name: "sum",
			Commands: []Command{
				{
					CMD:        "add",
					Params:     map[string]interface{}{"id": "$id", "operand1": "$value1", "operand2": "$value2"},
					ParentVars: nil,
				},
			},
		},
		"printAll": {
			Name: "printAll",
			Commands: []Command{
				{
					CMD:        "print",
					Params:     map[string]interface{}{"value": "#var1"},
					ParentVars: nil,
				},
				{
					CMD:        "print",
					Params:     map[string]interface{}{"value": "#var2"},
					ParentVars: nil,
				},
				{
					CMD:        "print",
					Params:     map[string]interface{}{"value": "#var3"},
					ParentVars: nil,
				},
			},
		},
	},
	Out: nil,
}

var emptyBody = `
{
  "var1":1,
  "var2":2,
  "var5":5,  

  "init": [
    {"cmd" : "#setup" }
  ],
  "setup": []
}
`

var cmdMissed = `
{
  "var1":1,
  "var2":2,
  "var5":5,  

  "init": [
    {"c" : "#setup" }
  ],
  "setup": []
}
`

func unmarshalDoc(assert *assert.Assertions, rawDoc string) map[string]interface{} {
	doc := make(map[string]interface{})
	dr := strings.NewReader(rawDoc)
	jd := json.NewDecoder(dr)
	err := jd.Decode(&doc)
	assert.NoError(err)

	return doc
}

func TestParse(t *testing.T) {
	assert := assert.New(t)
	v := NewVariables()
	f := NewScript()
	type args struct {
		doc  map[string]interface{}
		vars *Variables
		fns  *Script
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "script parse parse positive",
			args: args{
				doc:  unmarshalDoc(assert, testScript),
				vars: v,
				fns:  f,
			},
		},
		{
			name: "init missed",
			args: args{
				doc: func() map[string]interface{} {
					d := unmarshalDoc(assert, testScript)
					delete(d, InitFunc)
					return d
				}(),
				vars: v,
				fns:  f,
			},
			wantErr: true,
		},
		{
			name: "function body is empty",
			args: args{
				doc:  unmarshalDoc(assert, emptyBody),
				vars: v,
				fns:  f,
			},
			wantErr: true,
		},
		{
			name: "cmd key is missed",
			args: args{
				doc:  unmarshalDoc(assert, cmdMissed),
				vars: v,
				fns:  f,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("doc %v", tt.args.doc)
			err := Parse(tt.args.doc, tt.args.vars, tt.args.fns)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(&resVars, tt.args.vars)
				assert.Equal(&resFuncs, tt.args.fns)
			}
		})
	}
}

func TestScriptExecute(t *testing.T) {
	assert := assert.New(t)
	Vars = &resVars
	type testScript struct {
		Funcs map[string]*Func
		Out   *[]interface{}
	}
	type args struct {
		fName     string
		localVars map[string]interface{}
	}
	tests := []struct {
		name    string
		script  testScript
		args    args
		wantErr bool
	}{
		{
			name: "execute positive",
			script: testScript{
				Funcs: resFuncs.Funcs,
				Out:   &[]interface{}{3.5, 5.5, -1, 2, 5},
			},
			args: args{
				fName:     InitFunc,
				localVars: nil,
			},
		},
		{
			name: "unknown function",
			script: testScript{
				Funcs: resFuncs.Funcs,
				Out:   &[]interface{}{},
			},
			args: args{
				fName:     "wrong function name",
				localVars: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &Script{
				Funcs: tt.script.Funcs,
				Out:   tt.script.Out,
			}
			err := sc.Execute(tt.args.fName, tt.args.localVars)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tt.script.Out, sc.Out)
			}
		})
	}
}
