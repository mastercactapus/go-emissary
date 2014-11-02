package emissaryapi

import (
	"bytes"
	"fmt"
	sdunit "github.com/coreos/go-systemd/unit"
	"io/ioutil"
	"path"
	"strings"
)

var ValidUnitTypes = []string{"service", "socket", "device", "mount",
	"automount", "swap", "target", "path", "timer", "snapshot", "slice",
	"scope"}

var sdTruthy = []string{"1", "yes", "true", "on"}
var sdFalsey = []string{"0", "no", "false", "off"}

// UnitFile represents a systemd unit
type UnitFile struct {
	Options  []*sdunit.UnitOption
	Eoptions Eoptions
	Name     string
}

type Eoptions struct {
	Monitor     bool
	Global      bool
	Requires    []string
	Machines    []string
	Datacenters []string
	Tags        []string
}

// UnitTypeFromName returns the type from a systemd unit name
func UnitTypeFromName(name string) string {
	i := strings.LastIndex(name, ".")
	if i == -1 {
		return ""
	}
	return name[i+1:]
}

// UnitPrefixFromName returns the prefix of a systemd unit name
func UnitPrefixFromName(name string) string {
	i := strings.Index(name, "@")
	if i == -1 {
		i = strings.LastIndex(name, ".")
		if i == -1 {
			return ""
		}
	}
	return name[:i]
}

// UnitInstanceFromName returns the instance of a systemd unit
func UnitInstanceFromName(name string) string {
	s := strings.Index(name, "@")
	if s == -1 {
		return ""
	}
	e := strings.LastIndex(name, ".")
	if e == -1 {
		return ""
	}
	return name[s+1 : e]
}

func (u *UnitFile) Prefix() string {
	return UnitPrefixFromName(u.Name)
}
func (u *UnitFile) Type() string {
	return UnitTypeFromName(u.Name)
}
func (u *UnitFile) Instance() string {
	return UnitInstanceFromName(u.Name)
}
func (u *UnitFile) Serialize() []byte {
	read := sdunit.Serialize(u.Options)
	data, _ := ioutil.ReadAll(read)
	return data
}

func NewEoptions() *Eoptions {
	return &Eoptions{
		Datacenters: make([]string, 0, 5),
		Machines:    make([]string, 0, 5),
		Requires:    make([]string, 0, 10),
		Tags:        make([]string, 0, 30),
		Monitor:     true,
		Global:      false,
	}
}

func EoptionsFromUnitOptions(unitOptions []*sdunit.UnitOption) (e *Eoptions, err error) {
	e = NewEoptions()
	for _, v := range unitOptions {
		if v.Section != "X-Emissary" {
			continue
		}
		switch v.Name {
		case "Requires":
			e.Requires = append(e.Requires, v.Value)
		case "Machine":
			e.Machines = append(e.Machines, v.Value)
		case "Datacenter":
			e.Datacenters = append(e.Datacenters, v.Value)
		case "Tag":
			e.Tags = append(e.Tags, v.Value)
		case "Global":
			e.Global, err = sdBool(v.Value)
			if err != nil {
				return nil, err
			}
		case "Monitor":
			e.Monitor, err = sdBool(v.Value)
			if err != nil {
				return nil, err
			}
		}
	}

	return e, nil
}

func sdBool(val string) (bool, error) {
	str := strings.ToLower(val)
	if containsString(sdTruthy, str) {
		return true, nil
	} else if containsString(sdFalsey, str) {
		return false, nil
	} else {
		return false, fmt.Errorf("Cannot parse as bool '%s'", val)
	}
}

// NewUnitFile creates a new parsed UnitFile with the given name and data
func NewUnitFile(name string, data []byte) (unit *UnitFile, err error) {
	read := bytes.NewReader(data)
	opts, err := sdunit.Deserialize(read)
	if err != nil {
		return
	}
	eopts, err := EoptionsFromUnitOptions(opts)
	if err != nil {
		return
	}
	return &UnitFile{Name: name, Options: opts, Eoptions: *eopts}, nil
}

// NewUnitFromFile creates a new parsed UnitFile, using the basename as the unit name
func NewUnitFromFile(unitFilePath string) (unit *UnitFile, err error) {
	data, err := ioutil.ReadFile(unitFilePath)
	if err != nil {
		return
	}
	return NewUnitFile(path.Base(unitFilePath), data)
}
