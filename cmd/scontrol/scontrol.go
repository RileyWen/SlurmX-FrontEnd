package main

import (
	"SlurmXCli/generated/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type ServerAddr struct {
	ControlMachine 			string	`yaml:"ControlMachine"`
	SlurmCtlXdListenPort	string	`yaml:"SlurmCtlXdListenPort"`
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Arg must > 1")
		os.Exit(1)
	}

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

	if os.Args[1] == "show" {
		if os.Args[2] == "node" {
			var req *protos.QueryNodeInfoRequest
			queryAll := false
			nodeName := ""

			if len(os.Args) <= 3 {
				req = &protos.QueryNodeInfoRequest{NodeName: ""}
				queryAll = true
			} else {
				nodeName = os.Args[3]
				req = &protos.QueryNodeInfoRequest{NodeName: nodeName}
			}

			reply, err := stub.QueryNodeInfo(context.Background(), req)
			if err != nil {
				panic("QueryNodeInfo failed: " + err.Error())
			}

			if queryAll {
				if len(reply.NodeInfoList) == 0 {
					fmt.Printf("No node is avalable.\n")
				} else {
					for _, nodeInfo := range reply.NodeInfoList {
						fmt.Printf("NodeName=%v State=%v CPUs=%d AllocCpus=%d FreeCpus=%d\n\tRealMemory=%d AllocMem=%d FreeMem=%d\n\tPatition=%s RunningTask=%d\n\n", nodeInfo.Hostname, nodeInfo.State.String(), nodeInfo.Cpus, nodeInfo.AllocCpus, nodeInfo.FreeCpus, nodeInfo.RealMem, nodeInfo.AllocMem, nodeInfo.FreeMem, nodeInfo.PartitionName, nodeInfo.RunningTaskNum)
					}
				}
			} else {
				if len(reply.NodeInfoList) == 0 {
					fmt.Printf("Node %s not found.\n", nodeName)
				} else {
					for _, nodeInfo := range reply.NodeInfoList {
						fmt.Printf("NodeName=%v State=%v CPUs=%d AllocCpus=%d FreeCpus=%d\n\tRealMemory=%d AllocMem=%d FreeMem=%d\n\tPatition=%s RunningTask=%d\n\n", nodeInfo.Hostname, nodeInfo.State.String(), nodeInfo.Cpus, nodeInfo.AllocCpus, nodeInfo.FreeCpus, nodeInfo.RealMem, nodeInfo.AllocMem, nodeInfo.FreeMem, nodeInfo.PartitionName, nodeInfo.RunningTaskNum)
					}
				}
			}
		} else if os.Args[2] == "partition" {
			var req *protos.QueryPartitionInfoRequest
			queryAll := false
			partitionName := ""

			if len(os.Args) <= 3 {
				req = &protos.QueryPartitionInfoRequest{PartitionName: ""}
				queryAll = true
			} else {
				partitionName = os.Args[3]
				req = &protos.QueryPartitionInfoRequest{PartitionName: partitionName}
			}

			reply, err := stub.QueryPartitionInfo(context.Background(), req)
			if err != nil {
				panic("QueryPartitionInfo failed: " + err.Error())
			}

			if queryAll {
				if len(reply.PartitionInfo) == 0 {
					fmt.Printf("No node is avalable.\n")
				} else {
					for _, partitionInfo := range reply.PartitionInfo {
						fmt.Printf("PartitionName=%v State=%v\n\tTotalNodes=%d AliveNodes=%d\n\tTotalCpus=%d AvailCpus=%d AllocCpus=%d FreeCpus=%d\n\tTotalMem=%d AvailMem=%d AllocMem=%d FreeMem=%d\n\tHostList=%v\n", partitionInfo.Name, partitionInfo.State.String(), partitionInfo.TotalNodes, partitionInfo.AliveNodes, partitionInfo.TotalCpus, partitionInfo.AvailCpus, partitionInfo.AllocCpus, partitionInfo.FreeCpus, partitionInfo.TotalMem, partitionInfo.AvailMem, partitionInfo.AllocMem, partitionInfo.FreeMem, partitionInfo.Hostlist)
					}
				}
			} else {
				if len(reply.PartitionInfo) == 0 {
					fmt.Printf("Partition %s not found.\n", partitionName)
				} else {
					for _, partitionInfo := range reply.PartitionInfo {
						fmt.Printf("PartitionName=%v State=%v\n\tTotalNodes=%d AliveNodes=%d\n\tTotalCpus=%d AvailCpus=%d AllocCpus=%d FreeCpus=%d\n\tTotalMem=%d AvailMem=%d AllocMem=%d FreeMem=%d\n\tHostList=%v\n", partitionInfo.Name, partitionInfo.State.String(), partitionInfo.TotalNodes, partitionInfo.AliveNodes, partitionInfo.TotalCpus, partitionInfo.AvailCpus, partitionInfo.AllocCpus, partitionInfo.FreeCpus, partitionInfo.TotalMem, partitionInfo.AvailMem, partitionInfo.AllocMem, partitionInfo.FreeMem, partitionInfo.Hostlist)
					}
				}
			}
		}
	}
}
