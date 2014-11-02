package emissaryapi

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

func (c *ApiClient) FindUnit(unitName string) (unit *UnitFile, unitVersion string, err error) {

	unit, unitVersion, err = c.GetLatestUnit(unitName)
	if err != nil && strings.Contains(unitName, "@") {
		unitName = UnitPrefixFromName(unitName) + "@" + UnitTypeFromName(unitName)
		unit, unitVersion, err = c.GetLatestUnit(unitName)
	}

	return
}

func (c *ApiClient) GetLatestUnit(unitName string) (unit *UnitFile, unitVersion string, err error) {
	unitVersion, err = c.GetLatestUnitVersion(unitName)
	if err != nil {
		return
	}
	unit, err = c.GetUnit(unitName, unitVersion)
	return
}
func (c *ApiClient) GetLatestUnitVersion(unitName string) (unitVersion string, err error) {
	val, _, err := c.kv.Get(PrefixUnitFiles+unitName+"/latest", &consulapi.QueryOptions{Datacenter: c.dc})
	if err != nil {
		return
	}
	if val == nil {
		return "", ErrUnitNotFound
	}

	return string(val.Value), nil
}
func (c *ApiClient) GetUnit(unitName, tag string) (unit *UnitFile, err error) {
	val, _, err := c.kv.Get(PrefixUnitFiles+unitName+"/"+tag, &consulapi.QueryOptions{Datacenter: c.dc})
	if err != nil {
		return
	}
	if val == nil {
		return nil, ErrUnitNotFound
	}

	unit, err = NewUnitFile(unitName, val.Value)
	return
}
func (c *ApiClient) SetLatestUnit(unit *UnitFile) error {
	data := unit.Serialize()
	sum := sha256.Sum224(data)
	hash := hex.EncodeToString(sum[:])
	_, err := c.kv.Put(&consulapi.KVPair{Key: PrefixUnitFiles + unit.Name + "/" + hash, Value: data}, &c.w)
	if err != nil {
		return err
	}
	_, err = c.kv.Put(&consulapi.KVPair{Key: PrefixUnitFiles + unit.Name + "/latest", Value: []byte(hash)}, &c.w)
	if err != nil {
		return err
	}
	return nil
}
func (c *ApiClient) UnitExists(unitName string) (bool, error) {
	_, _, err := c.kv.Get(PrefixUnitFiles+unitName+"/latest", &c.q)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ApiClient) ListUnits(patterns ...string) ([]string, error) {
	list, _, err := c.kv.Keys("emissary/unit-files/", "/", &c.q)

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

func (c *ApiClient) DeleteUnit(unitName string) error {
	_, err := c.kv.Delete(PrefixUnitFiles+unitName+"/latest", &c.w)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DestroyUnit(unitName string) error {
	_, err := c.kv.DeleteTree(PrefixUnitFiles+unitName, &c.w)
	if err != nil {
		return err
	}
	return nil
}
