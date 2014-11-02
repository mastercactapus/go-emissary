package emissaryapi

import (
	"encoding/json"
	"github.com/armon/consul-api"
	"path"
)

const PrefixPendingUnits = "emissary/pending-units/"
const PrefixScheduledUnits = "emissary/scheduled-units/"

func (c *ApiClient) ScheduleUnit(unit *UnitFile, version string, activate bool) error {
	var e consulapi.UserEvent
	e.Name = "emissary-schedule-unit"
	data, err := json.Marshal(&unit.Eoptions)
	if err != nil {
		return err
	}
	e.Payload = data

	prefix := PrefixPendingUnits + unit.Name + "/"
	var state string
	if activate {
		state = "start"
	} else {
		state = "load"
	}

	_, err = c.kv.Put(&consulapi.KVPair{Key: prefix + "version", Value: []byte(version)}, &c.w)
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
