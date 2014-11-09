package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func listUnitsCommand(patterns ...string) {
	units, err := api.ScheduledUnits()
	if err != nil {
		fmt.Println("Could not list units:", err)
		os.Exit(2)
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.SetHeader([]string{"UNIT", "MACHINE", "ACTIVE", "SUB"})
	for _, v := range units {
		if matchAny(patterns, v.Name) {
			w.Append([]string{v.Name, v.MachineId, v.ActiveState, v.SubState})
		}
	}
	w.SetBorder(false)
	w.SetColumnSeparator("")
	w.SetAlignment(tablewriter.ALIGN_LEFT)
	w.SetCenterSeparator("")
	w.SetRowLine(false)
	w.Render()
}
