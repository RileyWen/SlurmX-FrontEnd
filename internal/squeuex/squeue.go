package squeuex

import (
	"SlurmXCli/generated/protos"
	"context"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

func Query(serverAddr string, partition string) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic("Cannot connect to SlurmCtlXd: " + err.Error())
	}

	stub := protos.NewSlurmCtlXdClient(conn)

	request := protos.QueryJobsInPartitionRequest{
		Partition: partition,
	}
	reply, err := stub.QueryJobsInPartition(context.Background(), &request)
	if err != nil {
		panic("QueryJobsInPartition failed: " + err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetHeaderLine(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetNoWhiteSpace(true)

	table.SetHeader([]string{"TaskId", "Type", "Status", "NodeIndex"})

	tableData := make([][]string, len(reply.TaskMetas))
	for _, taskMeta := range reply.TaskMetas {
		tableData = append(tableData, []string{
			strconv.FormatUint(uint64(taskMeta.TaskId), 10),
			taskMeta.Type.String(),
			taskMeta.Status.String(),
			strconv.FormatUint(uint64(taskMeta.NodeIndex), 10)})
	}

	table.AppendBulk(tableData)
	table.Render()
}
