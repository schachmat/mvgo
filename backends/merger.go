package backends

import (
	"flag"
	"log"
	"time"

	"github.com/schachmat/mvgo/iface"
)

type mergerConfig struct {
	first  string
	second string
	offset time.Duration
}

func (c *mergerConfig) Setup() {
	flag.StringVar(&c.first, "merger-first", "mvg-live", "merger backend: first `BACKEND`\n    \tThis has higher priority (e.g. live data)")
	flag.StringVar(&c.second, "merger-second", "mvv-efa", "merger backend: second `BACKEND`\n    \tThis has lower priority (e.g. schedule data)")
	flag.DurationVar(&c.offset, "merger-offset", 5*time.Minute, "merger backend: per line `OFFSET` between\n    \tend of primary data end beginning of secondary data")
}

func (c *mergerConfig) GetDepartures(station string) []iface.Departure {
	c1 := make(chan []iface.Departure)
	c2 := make(chan []iface.Departure)

	// get backends
	be1, ok := iface.AllBackends[c.first]
	if !ok {
		log.Fatalf("Merger Backend: Could not find primary backend \"%s\"", c.first)
	}
	be2, ok := iface.AllBackends[c.second]
	if !ok {
		log.Fatalf("Merger Backend: Could not find secondary backend \"%s\"", c.second)
	}

	// fetch data concurrently
	go func() {
		c1 <- be1.GetDepartures(station)
	}()
	go func() {
		c2 <- be2.GetDepartures(station)
	}()
	live := <-c1
	sched := <-c2

	// find max eta per line
	maxEta := make(map[string]time.Duration)
	for _, dep := range live {
		if maxEta[dep.Line] < dep.EtaNanoSec {
			maxEta[dep.Line] = dep.EtaNanoSec
		}
	}

	// remove lines occuring in live data from schedule data to prevent
	// duplicates in output.
	for _, ldep := range live {
		for i, sdep := range sched {
			if ldep.Line == sdep.Line && sdep.EtaNanoSec < maxEta[sdep.Line]+c.offset {
				sched = append(sched[:i], sched[i+1:]...)
			}
		}
	}

	all := append(live, sched...)
	iface.By(iface.Eta).Sort(all)
	return all
}

func init() {
	iface.AllBackends["merger"] = &mergerConfig{}
}
