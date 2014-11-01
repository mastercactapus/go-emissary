package main

import (
	"fmt"
	"path"
)

func UnitNameFromKey(keyName string) string {
	return path.Base(keyName[:len(keyName)-1])
}

func ConfirmYN(prompt string, format ...interface{}) bool {
	yes := []string{"y", "Y", "yes", "Yes", "YES"}
	no := []string{"n", "N", "no", "No", "NO"}

	response := ""
	for {
		fmt.Printf(prompt+" (y/N)", format...)
		fmt.Scanln(&response)
		if response == "" {
			return false
		} else if containsString(yes, response) {
			return true
		} else if containsString(no, response) {
			return false
		}
	}
}

func containsString(strs []string, str string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

func KeyExists(keyName string) bool {
	_, _, err := consul.KV().Get(keyName, nil)
	return err == nil
}
