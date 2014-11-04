package main

import (
	"fmt"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"io/ioutil"
	"time"
)

func SchedulerLoop() {
	ch := api.EventListener("emissary:pending-unit", time.Millisecond*250)
	for {
		e := <-ch
		unitName := string(e)
		fmt.Println("Attempting to schedule", unitName)
		err := LoadUnit(unitName)
		if err != nil && err != emissaryapi.ErrAlreadyDeployed {
			fmt.Println("Failed to schedule unit:", err)
		}
	}
}

func ScheduleAllLoop() {
	for {
		LoadAllUnits()
		time.Sleep(3 * time.Second)
	}
}

func QualifiedForUnit(unit *emissaryapi.UnitFile) bool {
	return true
}

func LoadAllUnits() {
	names, err := api.ScheduledUnits()
	if err != nil {
		return
	}

	for _, v := range names {
		e := LoadUnit(v)
		fmt.Println(e)
	}
}

func LoadUnit(unitName string) error {
	fmt.Println("load", unitName)
	s, err := api.ScheduledUnit(unitName)
	if err != nil {
		return err
	}
	err = api.AquireUnit(s, sessionId)
	if err != nil {
		return err
	}
	unit, err := api.GetUnit(s.Name, s.Version)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/tmp/"+unit.Name, unit.Serialize(), 0644)
	if err != nil {
		return err
	}
	return nil
}
