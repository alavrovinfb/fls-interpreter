package main

import "github.com/spf13/pflag"

const (
	cliMode     = "cli"
	serverMode  = "server"
	defaultMode = cliMode
)

var (
	defaultFiles = []string{}
	// define flag overrides
	flagMode       = pflag.String("mode", defaultMode, "FLS mode cli or server default cli")
	flagServerPort = pflag.StringSlice("script.files", defaultFiles, "comma separated list of files")
)
