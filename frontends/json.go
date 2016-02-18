package frontends

import (
	"os"
	"flag"
	"log"
	"encoding/json"

	"github.com/schachmat/mvgo/iface"
)

type feJson struct {
	noIndent bool
}

func (c *feJson) Setup() {
	flag.BoolVar(&c.noIndent, "json-no-indent", false, "json frontend: do not indent the output")
}

func (c *feJson) RenderDepartures(station string, deps []iface.Departure) {
	var b []byte
	var err error
	if c.noIndent {
		b, err = json.Marshal(deps)
	} else {
		b, err = json.MarshalIndent(deps, "", "\t")
	}
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}

func init() {
	iface.AllFrontends["json"] = &feJson{}
}
