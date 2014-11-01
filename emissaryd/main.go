package main

import (
	"fmt"
	"os"

	"github.com/armon/consul-api"
)

var consul *consulapi.Client

func RegisterService(unitName string) error {
	service := &consulapi.AgentServiceRegistration{
		Name: unitName,
		Check: &consulapi.AgentServiceCheck{
			TTL: "30s",
		},
	}
	return consul.Agent().ServiceRegister(service)
}

func PassService(unitName string) error {
	return consul.Agent().PassTTL("service:"+unitName, "")
}

func main() {
	conf := consulapi.DefaultConfig()
	c, err := consulapi.NewClient(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	consul = c
	err = RegisterService("emissary")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	PassService("emissary")
}
