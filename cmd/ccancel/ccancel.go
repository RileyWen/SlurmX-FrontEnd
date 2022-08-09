package main

import (
	"SlurmXCli/generated/protos"
	"SlurmXCli/internal/util"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
)

func main() {
	taskId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Invalid task id!")
	}

	path := "/etc/slurmx/config.yaml"
	config := util.ParseConfig(path)

	serverAddr := fmt.Sprintf("%s:%s", config.ControlMachine, config.SlurmCtlXdListenPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic("Cannot connect to SlurmCtlXd: " + err.Error())
	}

	stub := protos.NewSlurmCtlXdClient(conn)
	req := &protos.TerminateTaskRequest{TaskId: uint32(taskId)}

	reply, err := stub.TerminateTask(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to send TerminateTask gRPC: %s", err.Error())
	}

	if reply.Ok {
		fmt.Printf("Task #%d is terminating...", taskId)
	} else {
		fmt.Printf("Failed to terminating task #%d: %s", taskId, reply.Reason)
	}
}
