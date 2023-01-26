package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alavrovinfb/fls-interpreter/pkg/script"
)

var testScript = `
{
  "var1":1,
  "var2":2,
  
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

func main() {
	in := strings.NewReader(testScript)
	jDec := json.NewDecoder(in)
	docMap := make(map[string]interface{})
	if err := jDec.Decode(&docMap); err != nil {
		log.Fatal(err)
	}
	script.Vars, script.Body = script.Parse(docMap)
	fmt.Println("get maps", script.Vars, script.Body)
	script.Body.Execute(script.InitFunc, nil)

	fmt.Println("Print scrip vars", script.Vars)
}
