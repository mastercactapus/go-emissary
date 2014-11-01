package main

import (
	"fmt"
	"os"
)

func ListUnitFilesCommand(filter string) {
	list, _, err := consul.KV().Keys("emissary/unit-files/", "/", nil)

	if err != nil {
		fmt.Println("Failed to list units:", err)
		os.Exit(2)
	}

	for _, v := range list {
		name := UnitNameFromKey(v)
		fmt.Println(name)
	}
}
