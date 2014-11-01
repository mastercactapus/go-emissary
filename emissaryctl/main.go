package main

import (
	"fmt"
	"github.com/armon/consul-api"
	"github.com/mastercactapus/go-emissary/emissary-api"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
)

var (
	verbose              = kingpin.Flag("verbose", "Output extra information.").Default("false").Short('v').Bool()
	force                = kingpin.Flag("force", "Don't confirm for unsafe actions.").Default("false").Short('f').Bool()
	noblock              = kingpin.Flag("no-block", "Don't wait for actions to be completed remotely.").Default("false").Bool()
	dc                   = kingpin.Flag("datacenter", "Limit actions to a particular datacenter").Default("").Short('d').String()
	submit               = kingpin.Command("submit", "Submit/update one or more unit files.")
	submitUnits          = submit.Arg("unit-file", "One or more unit files to submit.").Required().Strings()
	listUnitFiles        = kingpin.Command("list-unit-files", "List all submitted unit files.")
	listUnitFilesPattern = listUnitFiles.Arg("PATTERN...", "Pattern(s) to match against unit names.").Strings()
	load                 = kingpin.Command("load", "Schedules units in the cluster.")
	loadUnits            = load.Arg("units", "One or more units to schedule in the cluster.").Strings()
	unload               = kingpin.Command("unload", "Unschedule units in the cluster.")
	unloadUnits          = unload.Arg("units", "One or more units to unschedule in the cluster.").Strings()
	listMachines         = kingpin.Command("list-machines", "List all known nodes.")
	listMachinesPattern  = listMachines.Arg("PATTERN...", "Pattern(s) to match against machine names.").Strings()
)

var consul *consulapi.Client
var store *emissaryapi.UnitStore

func main() {
	parsed := kingpin.Parse()
	conf := consulapi.DefaultConfig()
	c, err := consulapi.NewClient(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	consul = c
	store = emissaryapi.NewUnitStore(c, *dc)

	switch parsed {
	case "submit":
		submitUnitsCommand(*submitUnits)
	case "list-unit-files":
		listUnitFilesCommand(*listUnitFilesPattern...)
	case "load":
		loadUnitsCommand(*loadUnits)
	case "unload":
		unloadUnitsCommand(*unloadUnits)
	case "list-machines":
		listMachinesCommand(*listMachinesPattern...)
	}

}
