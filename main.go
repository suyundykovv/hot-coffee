package main

import (
	"fmt"
	"hot-coffee/config"
	"hot-coffee/flags"
	"hot-coffee/logging"
	"os"

	server "hot-coffee/server"
)

func main() {
	flags.Setup()
	//   "error": "Insufficient inventory for ingredient 'Milk'. Required: 200ml, Available: 150ml."
	if err := logging.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	server.Start(config.Port)
}
