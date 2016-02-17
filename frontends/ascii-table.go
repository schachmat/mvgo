package frontends

import (
	"fmt"
	"time"

	"github.com/schachmat/mvgo/iface"
)

type mvglAscii struct {
}

func (c *mvglAscii) Setup() {

}

func (c *mvglAscii) RenderDepartures(station string, deps []iface.Departure) {
	fmt.Println("The next departures from", station, "are:")
	for _, dep := range deps {
		fmt.Printf("%3d  %-4s  %s\n", dep.Eta/time.Minute, dep.Line, dep.Destination)
	}
}

func init() {
	iface.AllFrontends["ascii-table"] = &mvglAscii{}
}
