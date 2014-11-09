package main

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"os"
	"path"
	"time"
)

type ServiceStatus struct {
	Name   string
	Note   string
	Active bool
}

type ServiceState struct {
	Name        string
	Note        string
	Description string
	ActiveState string
	LoadState   string
	SubState    string
}

const ttl = 30 * time.Second
const updateAllInterval = 10 * time.Second

func SystemdMonitorLoop(interval time.Duration) {
	statCh, errCh := bus.SubscribeUnitsCustom(interval, 10, monitorUnitChanged, monitorUnitFilter)
	go func() {
		for {
			select {
			case stat := <-statCh:
				for k, v := range stat {
					if monitorUnitFilter(k) {
						continue
					}
					var err error
					if v == nil {
						err = api.UpdateUnitStates(k, "loaded", "inactive", "dead")
					} else {
						err = api.UpdateUnitStates(k, v.LoadState, v.ActiveState, v.SubState)
					}

					if err != nil {
						fmt.Println("Error updating service:", err)
					}
				}
			case err := <-errCh:
				fmt.Println("DBus error:", err)
			}
		}
	}()
}

func monitorUnitFilter(name string) bool {
	_, err := os.Stat(path.Join(*unitDir, name))
	return err != nil
}

func monitorUnitChanged(a *dbus.UnitStatus, b *dbus.UnitStatus) bool {
	return a.ActiveState != b.ActiveState || a.LoadState != b.LoadState || a.SubState != b.SubState
}
