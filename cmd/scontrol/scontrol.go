package main

import (
	"SlurmXCli/generated/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("Arg must > 2")
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic("Cannot connect to SlurmCtlXd: " + err.Error())
	}

	stub := protos.NewSlurmCtlXdClient(conn)

	if os.Args[2] == "show" {
		if os.Args[3] == "node" {
			var req *protos.QueryNodeInfoRequest
			queryAll := false
			nodeName := os.Args[4]

			if len(os.Args) <= 4 {
				req = &protos.QueryNodeInfoRequest{NodeName: ""}
				queryAll = true
			} else {
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
						fmt.Printf("NodeName=%v State=%v\n", nodeInfo.Hostname, nodeInfo.State.String())
					}
				}
			} else {
				if len(reply.NodeInfoList) == 0 {
					fmt.Printf("Node %s not found.\n", nodeName)
				} else {
					for _, nodeInfo := range reply.NodeInfoList {
						fmt.Printf("NodeName=%v State=%v\n", nodeInfo.Hostname, nodeInfo.State.String())
					}
				}
			}
		} else if os.Args[3] == "partition" {
			var req *protos.QueryPartitionInfoRequest
			queryAll := false
			partitionName := os.Args[4]

			if len(os.Args) <= 4 {
				req = &protos.QueryPartitionInfoRequest{PartitionName: ""}
				queryAll = true
			} else {
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
						fmt.Printf("PartitionName=%v State=%v HostList=%v\n", partitionInfo.Name, partitionInfo.State.String(), partitionInfo.Hostlist)
					}
				}
			} else {
				if len(reply.PartitionInfo) == 0 {
					fmt.Printf("Partition %s not found.\n", partitionName)
				} else {
					for _, partitionInfo := range reply.PartitionInfo {
						fmt.Printf("PartitionName=%v State=%v HostList=%v\n", partitionInfo.Name, partitionInfo.State.String(), partitionInfo.Hostlist)
					}
				}
			}
		}
	}

}
