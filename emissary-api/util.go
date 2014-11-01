package emissaryapi

import (
	"path/filepath"
)

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
