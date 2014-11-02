package main

import (
	"encoding/json"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"time"
)

func SchedulerLoop() {
	ch := api.EventListener("emissary:pending-unit", time.Millisecond*250)
	for {
		e := <-ch
		unit := new(emissaryapi.UnitFile)
		err := json.Unmarshal(e, unit)
		if err != nil {
			continue
		}
		//TODO: check if unit is valid here

	}
}
