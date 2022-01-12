package main

import (
	"SlurmXCli/internal/squeuex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("Arg must = 2")
		os.Exit(1)
	}
	squeuex.Query(os.Args[1], os.Args[2])
}
