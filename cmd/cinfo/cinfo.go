package main

import (
	"SlurmXCli/generated/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ServerAddr struct {
	ControlMachine 			string	`yaml:"ControlMachine"`
	SlurmCtlXdListenPort	string	`yaml:"SlurmCtlXdListenPort"`
}

func main() {

	confFile, err := ioutil.ReadFile("/etc/slurmx/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	confTxt := ServerAddr{}

	err = yaml.Unmarshal(confFile, &confTxt)
	if err != nil {
		log.Fatal(err)
	}
	ip := confTxt.ControlMachine
	port := confTxt.SlurmCtlXdListenPort

	serverAddr := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic("Cannot connect to SlurmCtlXd: " + err.Error())
	}

	stub := protos.NewSlurmCtlXdClient(conn)
	req := &protos.QueryClusterInfoRequest{}

	reply, err := stub.QueryClusterInfo(context.Background(), req)
	if err != nil {
		panic("QueryClusterInfo failed: " + err.Error())
	}

	if len(reply.PartitionNode) == 0 {
		fmt.Printf("No partition is available.\n")
	} else {
		fmt.Printf("PARTITION   AVAIL  TIMELIMIT  NODES  STATE  NODELIST\n")
		for _, partitionNode := range reply.PartitionNode {
			for _, commonNodeStateList := range partitionNode.CommonNodeStateList {
				if commonNodeStateList.NodesNum > 0 {
					fmt.Printf("%9s%8s%11s%7d%7s  %v\n", partitionNode.Name, partitionNode.State.String(), "infinite", commonNodeStateList.NodesNum, commonNodeStateList.State, commonNodeStateList.NodesList)
				}
			}
		}
	}

}