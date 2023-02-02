package main

import (
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
		if err := script.Body.Run(f); err != nil {
			return err
		}
		log.Printf("script %s done.", fName)
	}

	return nil
}
