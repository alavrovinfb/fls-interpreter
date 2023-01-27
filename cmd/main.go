package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	switch viper.GetString("mode") {
	case cliMode:
		if err := cli(viper.GetStringSlice("script.files")); err != nil {
			log.Fatal(err)
		}
	case serverMode:
		fmt.Println("Server mode")
	default:
		log.Fatal("unknown mod, could be either cli or server")
	}
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
