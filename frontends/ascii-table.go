package frontends

import (
	"flag"
	"fmt"

	"github.com/schachmat/mvgo/iface"
)

type atConfig struct {
	num int
}

func (c *atConfig) Setup() {
	flag.IntVar(&c.num, "ascii-table-num", 10, "ascii-table frontend: `NUMBER` of departures to display\n    \t0 means show all")
}

func (c *atConfig) RenderDepartures(station string, deps []iface.Departure) {
	fmt.Println("The next departures from", station, "are:")
	for i, dep := range deps {
		if c.num != 0 && i >= c.num {
			break
		}
		fmt.Printf("%3d  %-4s  %s\n", uint(dep.EtaNanoSec.Minutes()), dep.Line, dep.Destination)
	}
}

func init() {
	iface.AllFrontends["ascii-table"] = &atConfig{}
}
