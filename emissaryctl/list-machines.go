package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func listMachinesCommand(patterns ...string) {
	machines, err := api.ListMachinesPattern(patterns...)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"Name", "Address"})
	for _, v := range machines {
		w.Append([]string{v.Name, v.Address})
	}
	w.SetColumnSeparator("")
	w.SetBorder(false)
	w.Render()
}
