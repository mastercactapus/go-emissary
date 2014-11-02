package main

import (
	"time"
)

type ServiceStatus struct {
	Name   string
	Note   string
	Active bool
}

const ttl = 30 * time.Second
const updateAllInterval = 10 * time.Second

var status = make(map[string]*ServiceStatus, 100)

func UpdateAllLoop() {
	for {
		for _, v := range status {
			if v == nil {
				continue
			}
			if v.Active {
				api.ServicePass(v.Name, v.Note)
			} else {
				api.ServiceFail(v.Name, v.Note)
			}
		}
		time.Sleep(updateAllInterval)
	}
}

func UpdateImmediateLoop(statusCh chan *ServiceStatus) {
	for {
		s := <-statusCh
		if s.Active {
			api.ServicePass(s.Name, s.Note)
		} else {
			api.ServiceFail(s.Name, s.Note)
		}
	}
}

func MonitorServiceLoop(statusCh chan *ServiceStatus) {
	statCh, _ := bus.SubscribeUnits(time.Millisecond * 250)
	for {
		for k, v := range <-statCh {
			if status[k] == nil {
				continue
			}
			changed := false
			if v == nil {
				changed = updateStatusCheck(k, "Missing", false)
			} else {
				changed = updateStatusCheck(k, v.ActiveState+"/"+v.SubState, v.ActiveState == "active")
			}
			if changed {
				statusCh <- status[k]
			}
		}
	}
}

func updateStatusCheck(name, note string, active bool) bool {
	changed := status[name].Note != note
	changed = changed || status[name].Active != active
	status[name] = &ServiceStatus{name, note, active}
	return changed
}
