package main

import (
	"fmt"
	"github.com/armon/consul-api"
	"github.com/olekukonko/tablewriter"
	"os"
)

func listMachinesCommand(patterns ...string) {
	nodes, _, err := consul.Catalog().Nodes(&consulapi.QueryOptions{Datacenter: *dc})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Name", "Address"})
	for _, v := range nodes {
		w.Append([]string{v.Node, v.Address})
	}
	w.Render()
}
