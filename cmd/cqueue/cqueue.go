package main

import (
	"SlurmXCli/internal/squeuex"
	"SlurmXCli/internal/util"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Arg must > 1")
		os.Exit(1)
	}

	path := "/etc/crane/config.yaml"
	config := util.ParseConfig(path)

	serverAddr := fmt.Sprintf("%s:%s", config.ControlMachine, config.CraneCtldListenPort)

	squeuex.Query(serverAddr, os.Args[1])
}
