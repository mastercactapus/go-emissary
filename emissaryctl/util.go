package main

import "fmt"

func isPath(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == '.' || s[0] == '/'
}

func confirmYN(prompt string, format ...interface{}) bool {
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
