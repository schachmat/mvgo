package backends

import (
	"bytes"
	"encoding/json"
	"fmt"
	hcs "golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/schachmat/mvgo/iface"
)

type efaMvv struct {
}

type efaMvvResponse struct {
	Departures []struct {
		Eta  int `json:"countdown,string"`
		Line struct {
			Destination string `json:"direction"`
			Name        string `json:"symbol"`
		} `json:"servingLine"`
	} `json:"departureList"`
}

const (
	// &language=de // Tests showed this gets ignored by MVV, but set to language nevertheles
	// &name_dm=Freising // This is the actual query string
	efaMvvDuri = "http://efa.mvv-muenchen.de/mobile/XSLT_DM_REQUEST?outputFormat=JSON&stateless=1&coordOutputFormat=WGS84&type_dm=stop&itOptionsActive=1&ptOptionsActive=1&mergeDep=1&useAllStops=1&mode=direct"
)

func (c *efaMvv) Setup() {
}

func (c *efaMvv) GetDepartures(station string) []iface.Departure {
	var buf bytes.Buffer
	w := transform.NewWriter(&buf, charmap.ISO8859_1.NewEncoder())
	fmt.Fprintf(w, station)
	w.Close()

	res, err := http.Get(efaMvvDuri + "&name_dm=" + url.QueryEscape(string(buf.Bytes())))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	r, err := hcs.NewReader(res.Body, res.Header["Content-Type"][0])
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	var resp efaMvvResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Fatal(err)
	}

	var ret []iface.Departure
	for _, dep := range resp.Departures {
		ret = append(ret, iface.Departure{
			Line:        dep.Line.Name,
			Destination: dep.Line.Destination,
			Eta:         iface.JsonDuration(time.Duration(dep.Eta) * time.Minute),
		})
	}

	return ret
}

func init() {
	iface.AllBackends["efa-mvv"] = &efaMvv{}
}
