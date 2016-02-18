package frontends

import (
	"flag"
	"fmt"
	"time"

	"github.com/schachmat/mvgo/iface"
)

type mvglAscii struct {
	num int
}

func (c *mvglAscii) Setup() {
	flag.IntVar(&c.num, "ascii-table-num", 10, "ascii-table frontend: `NUMBER` of departures to display\n0 means show all")
}

func (c *mvglAscii) RenderDepartures(station string, deps []iface.Departure) {
	fmt.Println("The next departures from", station, "are:")
	for i, dep := range deps {
		if c.num != 0 && i >= c.num {
			break
		}
		fmt.Printf("%3d  %-4s  %s\n", dep.Eta/time.Minute, dep.Line, dep.Destination)
	}
}

func init() {
	iface.AllFrontends["ascii-table"] = &mvglAscii{}
}
