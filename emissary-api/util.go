package emissaryapi

import (
	"github.com/armon/consul-api"
	"path/filepath"
)

type KeyPutter interface {
	Put(*consulapi.KVPair, *consulapi.WriteOptions) (*consulapi.WriteMeta, error)
}

func containsUnit(units []UnitFile, unitName string) bool {
	for _, v := range units {
		if v.Name == unitName {
			return true
		}
	}
	return false
}

func containsString(strs []string, str string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

func matchAny(patterns []string, name string) bool {
	if patterns == nil {
		return true
	}
	for _, v := range patterns {
		if v == "" {
			return true
		}
		if m, _ := filepath.Match(v, name); m {
			return true
		}
	}
	return false
}
