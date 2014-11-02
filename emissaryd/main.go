package main

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"os"

	"github.com/armon/consul-api"
)

var consul *consulapi.Client
var bus *dbus.Conn
var api *emissaryapi.ApiClient

func AddService(name, note string, active bool) error {
	status[name] = &ServiceStatus{name, note, active}
	return api.RegisterService(name, ttl)
}
func RmService(name string) error {
	if status[name] != nil {
		delete(status, name)
	}
	return api.DeregisterService(name)
}

func main() {
	conf := consulapi.DefaultConfig()
	c, err := consulapi.NewClient(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	consul = c
	api = emissaryapi.NewClient(c, "")
	bus, err = dbus.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	err = AddService("emissary", "Running", true)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	statusUpdates := make(chan *ServiceStatus, 256)
	go MonitorServiceLoop(statusUpdates)
	go UpdateImmediateLoop(statusUpdates)
	UpdateAllLoop()
}
