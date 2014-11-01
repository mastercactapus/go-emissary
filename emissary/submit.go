package main

import (
	"fmt"
	"os"
	"path"
)

func submitUnitsCommand(units []string) {
	for _, v := range units {
		_, err := SubmitUnitFromFile(v)
		if err != nil {
			fmt.Println("Submit failed:", err)
			os.Exit(2)
		}
	}
}

func SubmitUnitFromFile(unitPath string) (unit *UnitFile, err error) {
	if !isPath(unitPath) {
		return nil, fmt.Errorf("Unit paths must be absolute or relative (start with '/' or '.')")
	}
	name := path.Base(unitPath)
	exists, err := store.Exists(name)
	if err != nil {
		return
	}
	if exists && !*force && !confirmYN("Unit '%s' already exists, update?", name) {
		return nil, fmt.Errorf("Unit '%s' has already been submitted.", name)
	}
	unit, err = NewUnitFromFile(unitPath)
	if err != nil {
		return
	}
	err = store.Set(name, unit)
	if err != nil {
		return
	}
	if *verbose {
		var actionStr string
		if exists {
			actionStr = "Updated"
		} else {
			actionStr = "Submitted"
		}
		fmt.Printf("%s unit '%s'\n", actionStr, name)
	}
	return
}
