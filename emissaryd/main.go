package main

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
	"time"

	"github.com/armon/consul-api"
)

//session-based lock on non _global to take deployment
//

var consul *consulapi.Client
var bus *dbus.Conn
var api *emissaryapi.ApiClient
var sessionId string
var self *emissaryapi.Machine

var (
	unitDir    = kingpin.Flag("unit-directory", "Directory to store scheduled unit files in.").Short('u').Default("/run/emissary/units").String()
	datacenter = kingpin.Flag("datacenter", "Specify the datacenter to operate in.").Short('d').Default("dc1").String()
)

func main() {
	kingpin.Parse()
	err := os.MkdirAll(*unitDir, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
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
	err = api.RegisterSelf([]string{}, 30)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	self, err = api.Self()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	go SystemdMonitorLoop(1 * time.Second)
	SyncUnitLoop(3 * time.Second)
}
