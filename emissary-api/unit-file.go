package emissaryapi

import (
	"bytes"
	sdunit "github.com/coreos/go-systemd/unit"
	"io/ioutil"
	"path"
	"strings"
)

var ValidUnitTypes = []string{"service", "socket", "device", "mount",
	"automount", "swap", "target", "path", "timer", "snapshot", "slice",
	"scope"}

// UnitFile represents a systemd unit
type UnitFile struct {
	Options []*sdunit.UnitOption
	Name    string
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

// NewUnitFile creates a new parsed UnitFile with the given name and data
func NewUnitFile(name string, data []byte) (unit *UnitFile, err error) {
	read := bytes.NewReader(data)
	opts, err := sdunit.Deserialize(read)
	if err != nil {
		return
	}
	return &UnitFile{Name: name, Options: opts}, nil
}

// NewUnitFromFile creates a new parsed UnitFile, using the basename as the unit name
func NewUnitFromFile(unitFilePath string) (unit *UnitFile, err error) {
	data, err := ioutil.ReadFile(unitFilePath)
	if err != nil {
		return
	}
	return NewUnitFile(path.Base(unitFilePath), data)
}
