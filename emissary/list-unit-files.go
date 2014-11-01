package main

import (
	"fmt"
	"os"
)

func listUnitFilesCommand(filter string) {
	list, err := store.List(filter)
	if err != nil {
		fmt.Println("Failed to list units:", err)
		os.Exit(2)
	}

	for _, v := range list {
		fmt.Println(v)
	}
}
