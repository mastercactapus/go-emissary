package emissaryapi

import (
	"errors"
	"github.com/armon/consul-api"
	"strings"
)

type ScheduledUnit struct {
	Name           string
	TargetState    string
	TargetVersion  string
	CurrentState   string
	CurrentVersion string
	MachineId      string
	Activate       string
	LoadState      string
	ActiveState    string
	SubState       string
}
type ScheduledUnits map[string]*ScheduledUnit

var ErrNoLock = errors.New("Could not establish a lock")

func (c *ApiClient) ScheduledUnits() (ScheduledUnits, error) {
	pairs, _, err := c.kv.List("emissary/scheduled-units/", &c.q)
	if err != nil {
		return nil, err
	}
	units := make(ScheduledUnits, len(pairs))
	for _, v := range pairs {
		if v.Value == nil {
			continue
		}
		s := strings.Split(strings.TrimPrefix(v.Key, "emissary/scheduled-units/"), "/")
		if len(s) != 2 {
			continue
		}
		name := s[0]
		vname := s[1]
		if units[name] == nil {
			units[name] = new(ScheduledUnit)
			units[name].Name = name
		}
		switch vname {
		case "target-state":
			units[name].TargetState = string(v.Value)
		case "target-version":
			units[name].TargetVersion = string(v.Value)
		case "current-state":
			units[name].CurrentState = string(v.Value)
		case "current-version":
			units[name].CurrentVersion = string(v.Value)
		case "current-loadstate":
			if v.Session == "" {
				continue
			}
			units[name].LoadState = string(v.Value)
		case "current-activestate":
			units[name].ActiveState = string(v.Value)
		case "current-substate":
			units[name].SubState = string(v.Value)
		case "machine-id":
			units[name].MachineId = string(v.Value)
		}
	}
	for k, v := range units {
		if v.ActiveState == "" && v.CurrentState == "failed" {
			units[k].ActiveState = "failed"
			units[k].SubState = "failed"
		}

		if v.MachineId == "" {
			units[k].ActiveState = "inactive"
			units[k].SubState = "dead"
		}
	}

	return units, nil
}

func (c *ApiClient) UpdateScheduleTarget(name, targetState, targetVersion string) error {
	topKey := "emissary/scheduled-units/" + name + "/"
	err := c.kvSet(topKey+"target-state", []byte(targetState))
	if err != nil {
		return err
	}
	if targetVersion != "" {
		err = c.kvSet(topKey+"target-version", []byte(targetVersion))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ApiClient) UpdateScheduleCurrent(name, currentState, currentVersion string) error {
	topKey := "emissary/scheduled-units/" + name + "/"
	err := c.kvSetSession(topKey+"current-state", []byte(currentState))
	if err != nil {
		return err
	}
	if currentVersion != "" {
		err = c.kvSetSession(topKey+"current-version", []byte(currentVersion))
		if err != nil {
			return err
		}
	}
	c.consul.Event().Fire(&consulapi.UserEvent{Name: "emissary:current-state:" + name, Payload: []byte(currentState), TagFilter: "emissary-client"}, &c.w)
	return nil
}

func (c *ApiClient) UpdateUnitStates(name, load, active, sub string) error {
	err := c.kvSetSession("emissary/scheduled-units/"+name+"/current-loadstate", []byte(load))
	if err != nil {
		return err
	}
	err = c.kvSetSession("emissary/scheduled-units/"+name+"/current-activestate", []byte(active))
	if err != nil {
		return err
	}
	err = c.kvSetSession("emissary/scheduled-units/"+name+"/current-substate", []byte(sub))
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) LockSchedule(name, machineId string) error {
	if c.sess == "" {
		return ErrNoSession
	}
	key := "emissary/schedule-lock"
	lock, _, err := c.kv.Acquire(&consulapi.KVPair{Key: key, Session: c.sess, Value: []byte(machineId)}, &c.w)
	if err != nil {
		return err
	}
	if !lock {
		return ErrNoLock
	}
	defer c.kv.Release(&consulapi.KVPair{Key: key, Session: c.sess}, &c.w)
	key = "emissary/scheduled-units/" + name + "/machine-id"
	lock, _, err = c.kv.Acquire(&consulapi.KVPair{Key: key, Session: c.sess, Value: []byte(machineId)}, &c.w)
	if err != nil {
		return err
	}
	if !lock {
		return ErrNoLock
	}
	return nil
}
