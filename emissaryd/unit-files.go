package main

import (
	"fmt"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var syncUnitsCh = make(chan bool, 1)

func SyncUnits() {
	//only append if there isn't already a pending request
	if len(syncUnitsCh) < cap(syncUnitsCh) {
		syncUnitsCh <- true
	}
}

func SyncUnitLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	doSync, _ := api.EventListener("emissary:schedule-unit", time.Millisecond*250)
	_SyncUnits(true)
	for {
		var err error
		select {
		case <-doSync:
			err = _SyncUnits(false)
		case <-ticker.C:
			err = _SyncUnits(false)
		case <-syncUnitsCh:
			err = _SyncUnits(false)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

//To schedule a unit, a node must have the least amount of units running of qualified machines, and be qualified
//Then aquire a lock to emissary/scheduler-lock, finally a lock on the unit itself
//_SyncUnits will do the actual heavy lifting after
func ScheduleUnit(name string, units emissaryapi.ScheduledUnits) {
	fmt.Println("TODO: Get qualified machines for a unit, do checks from there")
	counts := make(map[string]int, len(units))
	for _, v := range units {
		if v.MachineId == "" {
			continue
		}
		counts[v.MachineId]++
	}
	min := -1
	for _, c := range counts {
		if min == -1 {
			min = c
			continue
		}
		if c < min {
			min = c
		}
	}
	if min == -1 {
		min = 0
	}
	if counts[name] == min {
		scheduled := false
		avail := 0
		for k, v := range units {
			if v.MachineId == "" && !scheduled {
				err := api.LockSchedule(v.Name, name)
				if err != nil {
					fmt.Printf("Attempted and failed to schedule '%s': %s\n", v.Name, err.Error())
				} else {
					v.MachineId = self.FQDN()
					units[k] = v
					scheduled = true
				}
			} else if v.MachineId == "" {
				avail++
			}
		}

		if avail > 0 {
			//immediately resync, after the current loop since there are more to go
			SyncUnits()
		}
	} else {
		avail := 0
		for _, v := range units {
			if v.MachineId == "" {
				avail++
			}
		}

		//while there are things to be scheduled, go into overdrive :D
		//TODO: tweak this
		if avail > 0 {
			time.AfterFunc(500*time.Millisecond, SyncUnits)
		}
	}
}

func CleanupUnits(units emissaryapi.ScheduledUnits) error {
	files, err := ioutil.ReadDir(*unitDir)
	if err != nil {
		return err
	}
	m := make(map[string]bool, len(units))
	for _, v := range units {
		m[v.Name] = true
	}
	for _, v := range files {
		if !m[v.Name()] {
			_, err := DestroyUnit(v.Name())
			if err != nil {
				fmt.Printf("Failed to cleanup unit '%s': %s\n", v.Name(), err.Error())
				continue
			}
		}
	}
	return nil
}

func _SyncUnits(firstRun bool) error {
	units, err := api.ScheduledUnits()
	if err != nil {
		return err
	}

	var fqdn = self.FQDN()
	ScheduleUnit(fqdn, units)

	toStart := make([]string, 0, len(units))
	toLink := make([]string, 0, len(units))
	for _, v := range units {
		if v.MachineId != fqdn {
			continue
		}
		filename := path.Join(*unitDir, v.Name)
		_, err := os.Stat(filename)
		if err == nil && (!firstRun || v.CurrentVersion == v.TargetVersion) {
			continue
		}
		os.MkdirAll(*unitDir, 0755)

		f, err := api.GetUnit(v.Name, v.TargetVersion)
		if err != nil {
			fmt.Printf("Error getting unit '%s': %s\n", v.Name, err.Error())
			continue
		}
		err = ioutil.WriteFile(filename, f.Serialize(), 0644)
		if err != nil {
			fmt.Printf("Error saving unit '%s': %s\n", v.Name, err.Error())
			continue
		}
		toLink = append(toLink, filename)
		if v.TargetState == "active" {
			toStart = append(toStart, filename)
		}
	}

	_, err = bus.LinkUnitFiles(toLink, true, true)
	if err != nil {
		return err
	}
	if len(toLink) > 0 {
		//daemon-reload
		bus.Reload()
	}
	for _, v := range toStart {
		go StartUnit(path.Base(v))
	}

	err = CleanupUnits(units)
	if err != nil {
		fmt.Println("Could not cleanup units:", err)
	}

	return nil
}

func DestroyUnit(name string) (string, error) {
	ch := make(chan string)
	_, err := bus.StopUnit(name, "replace", ch)
	if err != nil {
		fmt.Printf("Failed to stop unit '%s': %s\n", name, err.Error())
		return "", err
	}
	result := <-ch
	err = os.Remove(path.Join(*unitDir, name))
	return result, err
}

func StartUnit(name string) error {
	ch := make(chan string)
	_, err := bus.StartUnit(name, "replace", ch)
	if err != nil {
		fmt.Printf("Failed to start unit '%s': %s\n", name, err.Error())
		return err
	}
	api.UpdateScheduleCurrent(name, "activating", "")
	result := <-ch
	switch result {
	case "done":
		api.UpdateScheduleCurrent(name, "active", "")
	default:
		api.UpdateScheduleCurrent(name, result, "")
	}
	return nil
}
