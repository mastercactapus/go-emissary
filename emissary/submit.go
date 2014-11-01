package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/armon/consul-api"
	"io/ioutil"
	"os"
	"path"
)

func SubmitUnitsCommand(units []string) {
	for _, v := range units {
		err := SubmitUnit(v)
		if err != nil {
			fmt.Println("Submit failed:", err)
			os.Exit(2)
		}
	}
}

func SubmitUnit(unitPath string) error {
	name := path.Base(unitPath)
	exists := false
	if KeyExists("emissary/unit-files/" + name + "/latest") {
		exists = true
		if !*force && !ConfirmYN("Unit '%s' already exists, update?", name) {
			return fmt.Errorf("Unit '%s' has already been submitted.", name)
		}
	}
	data, err := ioutil.ReadFile(unitPath)
	if err != nil {
		return err
	}
	rawHash := sha256.Sum224(data)
	hash := hex.EncodeToString(rawHash[:])
	var pair consulapi.KVPair
	pair.Key = "emissary/unit-files/" + name + "/" + hash
	pair.Value = data
	_, err = consul.KV().Put(&pair, nil)
	if err != nil {
		return err
	}
	pair.Key = "emissary/unit-files/" + name + "/latest"
	pair.Value = []byte(hash)
	_, err = consul.KV().Put(&pair, nil)
	if err != nil {
		return err
	}
	if *verbose {
		var actionStr string
		if exists {
			actionStr = "Updated"
		} else {
			actionStr = "Submitted"
		}
		fmt.Printf("%s unit '%s'\n", actionStr, name)
	}
	return nil
}
