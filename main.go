package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/schachmat/ingo"
	_ "github.com/schachmat/mvgo/backends"
	_ "github.com/schachmat/mvgo/frontends"
	"github.com/schachmat/mvgo/iface"
)

func main() {
	configpath := os.Getenv("MVGORC")
	if configpath == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("%v\nYou can set the environment variable MVGORC to point to your config file as a workaround.", err)
		}
		configpath = path.Join(usr.HomeDir, ".mvgorc")
	}

	// initialize backends and frontends (flags and default config)
	for _, be := range iface.AllBackends {
		be.Setup()
	}
	for _, fe := range iface.AllFrontends {
		fe.Setup()
	}

	// initialize global flags and default config
	station := flag.String("station", "Marienplatz", "Which `STATION` should be querried")
	selectedBackend := flag.String("backend", "mvg-live", "`BACKEND` to be used")
	selectedFrontend := flag.String("frontend", "ascii-table", "`FRONTEND` to be used")

	// read/write config and parse flags
	ingo.Parse(configpath)

	// non-flag argument overwrites station flag
	if len(flag.Args()) > 0 {
		*station = flag.Args()[0]
	}

	// get selected backend and fetch the departure data from it
	be, ok := iface.AllBackends[*selectedBackend]
	if !ok {
		log.Fatalf("Could not find selected backend \"%s\"", *selectedBackend)
	}
	r := be.GetDepartures(*station)

	// get selected frontend and render the result
	fe, ok := iface.AllFrontends[*selectedFrontend]
	if !ok {
		log.Fatalf("Could not find selected frontend \"%s\"", *selectedFrontend)
	}
	fe.RenderDepartures(*station, r)
}
