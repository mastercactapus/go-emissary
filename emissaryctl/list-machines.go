package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

func listMachinesCommand(patterns ...string) {
	machines, err := api.Machines(*dc, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Machine", "Name", "Address", "Metadata"})
	for _, v := range machines {
		w.Append([]string{v.FQDN(), v.Name, v.Address, strings.Join(v.Metadata, ",")})
	}
	w.SetColumnSeparator("")
	w.SetCenterSeparator("")
	w.SetBorder(false)
	w.Render()
}
