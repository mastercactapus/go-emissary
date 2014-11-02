package main

import (
	"errors"
	"fmt"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"os"
	"path"
)

var ErrAlreadyActive = errors.New("Unit is already active")
var ErrLockFailed = errors.New("Failed to aquire lock")

func loadUnitsCommand(units []string) {
	if len(units) == 0 {
		fmt.Println("You must specify at least one unit to load.")
		os.Exit(2)
	}
	for _, v := range units {
		_, _, err := LoadUnit(v)
		if err != nil {
			fmt.Printf("Failed to load %s: %s\n", v, err.Error())
			os.Exit(2)
		}
	}
}

func LoadUnit(unitPath string) (unit *emissaryapi.UnitFile, version string, err error) {
	name := path.Base(unitPath)
	if isPath(unitPath) {
		unit, err = SubmitUnitFromFile(unitPath)
	}
	if err != nil {
		return
	}

	if !containsString(emissaryapi.ValidUnitTypes, emissaryapi.UnitTypeFromName(name)) {
		name += "." + emissaryapi.ValidUnitTypes[0]
	}

	unit, version, err = api.Store.Find(name)
	if err != nil {
		return
	}

	err = api.ScheduleUnit(unit, version)
	if err != nil {
		return
	}

	return
}
