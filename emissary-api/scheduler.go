package emissaryapi

import (
	"errors"
	"fmt"
	"github.com/armon/consul-api"
	"path"
)

const PrefixScheduledUnits = "emissary/scheduled-units/"

var ErrNoLock = errors.New("Could not aquire lock")
var ErrAlreadyDeployed = errors.New("Unit is already deployed")

type ScheduledUnit struct {
	Name           string
	Version        string
	Machine        string
	MachineState   string
	MachineVersion string
	Activate       string
}

func (c *ApiClient) AquireUnit(s *ScheduledUnit, sessionId string) error {
	if s.Machine != "" {
		return ErrAlreadyDeployed
	}
	nodeName, err := c.consul.Agent().NodeName()
	if err != nil {
		return err
	}
	gotLock, _, err := c.kv.Acquire(&consulapi.KVPair{Key: PrefixScheduledUnits + s.Name + "/", Session: sessionId}, &c.w)
	if err != nil {
		return err
	}
	if !gotLock {
		return ErrNoLock
	}
	defer c.kv.Release(&consulapi.KVPair{Key: PrefixScheduledUnits + s.Name + "/", Session: sessionId}, &c.w)
	_, err = c.kv.Put(&consulapi.KVPair{Key: PrefixScheduledUnits + s.Name + "/machine", Value: []byte(nodeName), Session: sessionId}, &c.w)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(&consulapi.KVPair{Key: PrefixScheduledUnits + s.Name + "/machine-version", Value: []byte(s.Version), Session: sessionId}, &c.w)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(&consulapi.KVPair{Key: PrefixScheduledUnits + s.Name + "/machine-state", Value: []byte("deploying unit"), Session: sessionId}, &c.w)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ScheduleUnit(unit *UnitFile, version string, activate bool) error {
	var e consulapi.UserEvent
	e.Name = "emissary:pending-unit"
	e.Payload = []byte(unit.Name)

	prefix := PrefixScheduledUnits + unit.Name + "/"
	var state string
	if activate {
		state = "start"
	} else {
		state = "load"
	}

	_, err := c.kv.Put(&consulapi.KVPair{Key: prefix + "version", Value: []byte(version)}, &c.w)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(&consulapi.KVPair{Key: prefix + "activate", Value: []byte(state)}, &c.w)
	if err != nil {
		return err
	}

	_, _, err = c.consul.Event().Fire(&e, &c.w)
	if err != nil {
		return err
	}

	return nil
}

func (c *ApiClient) ScheduledUnit(name string) (*ScheduledUnit, error) {
	s := new(ScheduledUnit)
	s.Name = name
	fmt.Println(PrefixScheduledUnits + name + "/version")
	kv, _, err := c.kv.Get(PrefixScheduledUnits+name+"/version", &c.q)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, ErrUnitNotFound
	}
	s.Version = string(kv.Value)
	kv, _, err = c.kv.Get(PrefixScheduledUnits+name+"/activate", &c.q)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, ErrUnitNotFound
	}
	s.Activate = string(kv.Value)
	kv, _, err = c.kv.Get(PrefixScheduledUnits+name+"/machine", &c.q)
	if err != nil {
		return nil, err
	}
	if kv != nil {
		s.Machine = string(kv.Value)
	}
	kv, _, err = c.kv.Get(PrefixScheduledUnits+name+"/machine-version", &c.q)
	if err != nil {
		return nil, err
	}
	if kv != nil {
		s.MachineVersion = string(kv.Value)
	}
	kv, _, err = c.kv.Get(PrefixScheduledUnits+name+"/machine-state", &c.q)
	if err != nil {
		return nil, err
	}
	if kv != nil {
		s.MachineState = string(kv.Value)
	}
	return s, nil
}

func (c *ApiClient) ScheduledUnits() ([]string, error) {
	names, _, err := c.kv.Keys(PrefixScheduledUnits, "/", &c.q)
	if err != nil {
		return nil, err
	}
	for k, v := range names {
		names[k] = path.Base(v)
	}
	return names, nil
}

func (c *ApiClient) LocalScheduledUnits() (units []UnitFile, err error) {
	name, err := c.consul.Agent().NodeName()
	if err != nil {
		return
	}

	globals, _, err := c.kv.List(PrefixScheduledUnits+"_global/", &c.q)
	if err != nil {
		return
	}
	machine, _, err := c.kv.List(PrefixScheduledUnits+name+"/", &c.q)
	if err != nil {
		return
	}

	units = make([]UnitFile, 0, len(globals)+len(machine))
	for _, v := range globals {
		name := path.Base(v.Key[:len(v.Key)-1])
		if !containsUnit(units, name) {
			unit, err := c.GetUnit(name, string(v.Value))
			if err != nil {
				return nil, err
			}
			units = append(units, *unit)
		}
	}
	for _, v := range machine {
		name := path.Base(v.Key[:len(v.Key)-1])
		if !containsUnit(units, name) {
			unit, err := c.GetUnit(name, string(v.Value))
			if err != nil {
				return nil, err
			}
			units = append(units, *unit)
		}
	}
	return
}
