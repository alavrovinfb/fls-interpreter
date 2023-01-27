package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alavrovinfb/fls-interpreter/pkg/script"
)

func cli(files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("file list is not provided")
	}
	for _, fName := range files {
		log.Printf("executing script %s", fName)
		f, err := os.Open(strings.TrimSpace(fName))
		if err != nil {
			return err
		}
		defer f.Close()
		jDec := json.NewDecoder(f)
		docMap := make(map[string]interface{})
		if err := jDec.Decode(&docMap); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
		if err := script.Parse(docMap, script.Vars, script.Body); err != nil {
			return err
		}
		script.Body.RestOut()
		if err := script.Body.Execute(script.InitFunc, nil); err != nil {
			return err
		}
		log.Printf("script %s done.", fName)
	}

	return nil
}
