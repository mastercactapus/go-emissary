package main

import (
	"fmt"
	"os"
)

func listUnitFilesCommand(patterns ...string) {
	list, err := store.List(patterns...)
	if err != nil {
		fmt.Println("Failed to list units:", err)
		os.Exit(2)
	}

	for _, v := range list {
		fmt.Println(v)
	}
}
