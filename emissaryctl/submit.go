package main

import (
	"fmt"
	"os"

	"github.com/mastercactapus/go-emissary/emissary-api"
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

func SubmitUnitFromFile(unitFilePath string) (unit *emissaryapi.UnitFile, err error) {
	if !isPath(unitFilePath) {
		return nil, fmt.Errorf("Unit paths must be absolute or relative (start with '/' or '.')")
	}
	unit, err = emissaryapi.NewUnitFromFile(unitFilePath)
	if err != nil {
		return
	}
	exists, err := store.Exists(unit.Name)
	if err != nil {
		return
	}
	if exists && !*force && !confirmYN("Unit '%s' already exists, update?", unit.Name) {
		return nil, fmt.Errorf("Unit '%s' has already been submitted.", unit.Name)
	}

	err = store.SetLatest(unit)
	if err != nil {
		return
	}
	return
}
