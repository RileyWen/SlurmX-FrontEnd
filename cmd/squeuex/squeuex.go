package main

import (
	"SlurmXCli/internal/squeuex"
	"SlurmXCli/internal/util"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Arg must = 1")
		os.Exit(1)
	}

	path := "/etc/slurmx/config.yaml"
	config := util.ParseConfig(path)

	serverAddr := fmt.Sprintf("%s:%s", config.ControlMachine, config.SlurmCtlXdListenPort)

	squeuex.Query(serverAddr, os.Args[1])
}
