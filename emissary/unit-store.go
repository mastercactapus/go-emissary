package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/armon/consul-api"
)

type UnitStore struct {
	kv *consulapi.KV
	dc string //Datacenter
}

const PrefixUnitFiles = "emissary/unit-files/"

var ErrInvalidUnitName = errors.New("Invalid unit name")
var ErrUnitNotFound = errors.New("No unit exists in store by that name")

func NewUnitStore(c *consulapi.Client, datacenter string) *UnitStore {
	return &UnitStore{kv: c.KV(), dc: datacenter}
}

func (s *UnitStore) Find(unitName string) (unit *UnitFile, unitVersion string, err error) {

	unit, unitVersion, err = s.GetLatest(unitName)
	if err != nil && strings.Contains(unitName, "@") {
		unitName = UnitPrefixFromName(unitName) + "@" + UnitTypeFromName(unitName)
		unit, unitVersion, err = s.GetLatest(unitName)
	}

	return
}

func (s *UnitStore) GetLatest(unitName string) (unit *UnitFile, unitVersion string, err error) {
	unitVersion, err = s.GetLatestVersion(unitName)
	if err != nil {
		return
	}
	unit, err = s.Get(unitName, unitVersion)
	return
}
func (s *UnitStore) GetLatestVersion(unitName string) (unitVersion string, err error) {
	val, _, err := s.kv.Get(PrefixUnitFiles+unitName+"/latest", &consulapi.QueryOptions{Datacenter: s.dc})
	if err != nil {
		return
	}
	if val == nil {
		return "", ErrUnitNotFound
	}

	return string(val.Value), nil
}
func (s *UnitStore) Get(unitName, tag string) (unit *UnitFile, err error) {
	val, _, err := s.kv.Get(PrefixUnitFiles+unitName+"/"+tag, &consulapi.QueryOptions{Datacenter: s.dc})
	if err != nil {
		return
	}
	if val == nil {
		return nil, ErrUnitNotFound
	}

	unit, err = NewUnitFile(unitName, val.Value)
	return
}
func (s *UnitStore) Set(unitName string, unit *UnitFile) error {
	data := unit.Serialize()
	sum := sha256.Sum224(data)
	hash := hex.EncodeToString(sum[:])
	_, err := s.kv.Put(&consulapi.KVPair{Key: PrefixUnitFiles + unitName + "/" + hash, Value: data}, &consulapi.WriteOptions{Datacenter: s.dc})
	if err != nil {
		return err
	}
	_, err = s.kv.Put(&consulapi.KVPair{Key: PrefixUnitFiles + unitName + "/latest", Value: []byte(hash)}, &consulapi.WriteOptions{Datacenter: s.dc})
	if err != nil {
		return err
	}
	return nil
}
func (s *UnitStore) Exists(unitName string) (bool, error) {
	_, _, err := s.kv.Get(PrefixUnitFiles+unitName+"/latest", &consulapi.QueryOptions{Datacenter: s.dc})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *UnitStore) List(patterns ...string) ([]string, error) {
	list, _, err := consul.KV().Keys("emissary/unit-files/", "/", &consulapi.QueryOptions{Datacenter: s.dc})

	if err != nil {
		fmt.Println("Failed to list units:", err)
		os.Exit(2)
	}
	names := make([]string, 0, len(list))
	for _, v := range list {
		name := path.Base(v[:len(v)-1])
		if matchAny(patterns, name) {
			names = append(names, name)
		}
	}

	return names, nil
}

func (s *UnitStore) Delete(unitName string) error {
	_, err := s.kv.Delete(PrefixUnitFiles+unitName+"/latest", &consulapi.WriteOptions{Datacenter: s.dc})
	if err != nil {
		return err
	}
	return nil
}

func (s *UnitStore) DeleteAll(unitName string) error {
	_, err := s.kv.DeleteTree(PrefixUnitFiles+unitName, &consulapi.WriteOptions{Datacenter: s.dc})
	if err != nil {
		return err
	}
	return nil
}
