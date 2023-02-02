package main

import (
	"log"

	"github.com/alavrovinfb/fls-interpreter/internal/server"
)

const (
	cliMode    = "cli"
	serverMode = "server"
)

func main() {
	switch *flagMode {
	case cliMode:
		if err := cli(*flagScriptFiles); err != nil {
			log.Fatal(err)
		}
	case serverMode:
		l := server.NewLogger(*flagLoggingLevel)
		i := server.NewInterpreter(l, flagServerAddress, flagServerPort,
			flagGatewayAddress, flagGatewayPort, flagGatewayURL,
			flagInternalAddress, flagInternalPort, flagInternalHealth)
		if err := i.Run(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown mod, could be either cli or server")
	}
}
