package backends

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	hcs "golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/schachmat/mvgo/iface"
)

type mvgLiveConfig struct {
	bus      bool
	tram     bool
	subway   bool
	suburban bool
	ret      []iface.Departure
}

const (
	mvgLiveDuri = "https://www.mvg-live.de/ims/dfiStaticAuswahl.svc?"
)

func (c *mvgLiveConfig) parse_departure_line(n *html.Node) {
	var line, station, departure string
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode || child.Data != "td" {
			continue
		}
		for _, a := range child.Attr {
			if a.Key != "class" {
				continue
			}
			if a.Val == "lineColumn" {
				line = child.FirstChild.Data
				break
			}
			if a.Val == "stationColumn" {
				station = strings.TrimSpace(child.FirstChild.Data)
				break
			}
			if a.Val == "inMinColumn" {
				departure = child.FirstChild.Data
				break
			}
		}
	}
	eta, err := time.ParseDuration(departure + "m")
	if err != nil {
		log.Println("Unable to parse departure time:", err)
	}
	c.ret = append(c.ret, iface.Departure{
		Line:        line,
		Destination: station,
		EtaNanoSec:  eta,
	})
}

func (c *mvgLiveConfig) parse_departures(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		for _, a := range n.Attr {
			if a.Key == "class" && (a.Val == "rowOdd" || a.Val == "rowEven") {
				c.parse_departure_line(n)
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.parse_departures(child)
	}
}

//func (c *mvgLiveConfig) parse_stations(n *html.Node) {
//	if n.Type == html.ElementNode && n.Data == "a" {
//		for _, a := range n.Attr {
//			if a.Key == "href" && strings.HasPrefix(a.Val, "/ims/dfiStatic") {
//				fmt.Println(a.Val, n.FirstChild.Data)
//			}
//		}
//	}
//	for child := n.FirstChild; child != nil; child = child.NextSibling {
//		c.parse_stations(child)
//	}
//}

func (c *mvgLiveConfig) find_table(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "departureTable departureView" {
				c.parse_departures(n)
				return n
			}
			//			 else if a.Key == "class" && a.Val == "departureTable header" {
			//				c.parse_stations(n)
			//				return n
			//			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		t := c.find_table(child)
		if t != nil {
			return t
		}
	}
	return nil
}

func (c *mvgLiveConfig) Setup() {
	flag.BoolVar(&c.bus, "mvgl-bus", true, "mvg-live backend: show bus departures")
	flag.BoolVar(&c.tram, "mvgl-tram", true, "mvg-live backend: show tram departures")
	flag.BoolVar(&c.subway, "mvgl-subway", true, "mvg-live backend: show subway departures")
	flag.BoolVar(&c.suburban, "mvgl-suburban", true, "mvg-live backend: show suburban departures")
}

func (c *mvgLiveConfig) GetDepartures(station string) []iface.Departure {
	params := make([]string, 5)

	var buf bytes.Buffer
	w := transform.NewWriter(&buf, charmap.Windows1252.NewEncoder())
	fmt.Fprintf(w, station)
	w.Close()
	params = append(params, "haltestelle="+url.QueryEscape(string(buf.Bytes())))

	if c.bus {
		params = append(params, "bus=checked")
	}
	if c.tram {
		params = append(params, "tram=checked")
	}
	if c.subway {
		params = append(params, "ubahn=checked")
	}
	if c.suburban {
		params = append(params, "sbahn=checked")
	}

	res, err := http.Get(mvgLiveDuri + strings.Join(params, "&"))
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	c.ret = []iface.Departure{}
	r, err := hcs.NewReader(res.Body, res.Header["Content-Type"][0])
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	c.find_table(doc)
	return c.ret
}

func init() {
	iface.AllBackends["mvg-live"] = &mvgLiveConfig{}
}
