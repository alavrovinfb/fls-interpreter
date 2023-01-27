package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_getID(t *testing.T) {
	assert := assert.New(t)
	type command struct {
		CMD        string
		Params     map[string]interface{}
		ParentVars map[string]interface{}
	}
	tests := []struct {
		name     string
		testCmds command
		want     string
		wantErr  bool
	}{
		{
			name: "getID test positive",
			testCmds: command{
				Params:     map[string]interface{}{"id": "$id", "operand1": "$value1", "operand2": "$value2"},
				ParentVars: map[string]interface{}{"id": "var1", "value1": "#var1", "value2": "#var2"},
			},
			want: "var1",
		},
		{
			name: "incorrect id type",
			testCmds: command{
				Params:     map[string]interface{}{"id": "$id", "operand1": "$value1", "operand2": "$value2"},
				ParentVars: map[string]interface{}{"id": true, "value1": "#var1", "value2": "#var2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				CMD:        tt.testCmds.CMD,
				Params:     tt.testCmds.Params,
				ParentVars: tt.testCmds.ParentVars,
			}
			got, err := c.getID()
			if tt.wantErr {
				assert.Error(err)
			}

			assert.Equal(tt.want, got)
		})
	}
}

var tstVars = Variables{
	V: map[string]Value{
		"var1": 1, "var2": 2, "var5": 5,
	},
}

func TestCommand_processParams(t *testing.T) {
	assert := assert.New(t)
	Vars = &tstVars
	type command struct {
		CMD        string
		Params     map[string]interface{}
		ParentVars map[string]interface{}
	}
	type args struct {
		opType int
	}
	tests := []struct {
		name     string
		testCmds command
		args     args
		want     []Value
	}{
		{
			name: "test process unary params",
			testCmds: command{
				Params: map[string]interface{}{"id": "var1", "value": 3.5},
			},
			args: args{opType: UNARY},
			want: []Value{3.5},
		},
		{
			name: "test process binary direct params",
			testCmds: command{
				Params: map[string]interface{}{"id": "var1", "value1": float64(4), "value2": float64(3)},
			},
			args: args{opType: BINARY},
			want: []Value{4, 3},
		},
		{
			name: "test process binary nester params",
			testCmds: command{
				Params: map[string]interface{}{"id": "var1", "value1": "#var1", "value2": "#var2"},
			},
			args: args{opType: BINARY},
			want: []Value{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				CMD:        tt.testCmds.CMD,
				Params:     tt.testCmds.Params,
				ParentVars: tt.testCmds.ParentVars,
			}
			got := c.processParams(tt.args.opType)
			assert.Equal(tt.want, got)
		})
	}
}

func TestCommand_Do(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		CMD        string
		Params     map[string]interface{}
		ParentVars map[string]interface{}
	}
	type args struct {
		out *[]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "Do create test",
			fields: fields{
				CMD:    CREATE,
				Params: map[string]interface{}{"id": "var7", "value": float64(77)},
			},
			want: Value(77),
		},
		{
			name: "Do update test",
			fields: fields{
				CMD:    CREATE,
				Params: map[string]interface{}{"id": "var7", "value": float64(88)},
			},
			want: Value(88),
		},
		{
			name: "Do update test",
			fields: fields{
				CMD:    DELETE,
				Params: map[string]interface{}{"id": "var7"},
			},
			want: Value(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				CMD:        tt.fields.CMD,
				Params:     tt.fields.Params,
				ParentVars: tt.fields.ParentVars,
			}
			err := c.Do(tt.args.out)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.Equal(tt.want, Vars.V["var7"])
			}
		})
	}
}
