package iface

import (
	"sort"
	"time"
)

// Sorting stuff totaly stolen from pkg sort documentation. Use like this:
//     By(Eta).Sort(departures)
func (by By) Sort(deps []Departure) {
	ps := &depSorter{
		deps: deps,
		by:   by,
	}
	sort.Sort(ps)
}

type By func(d1, d2 *Departure) bool

type depSorter struct {
	deps []Departure
	by   func(d1, d2 *Departure) bool
}

func (s *depSorter) Len() int {
	return len(s.deps)
}

func (s *depSorter) Swap(i, j int) {
	s.deps[i], s.deps[j] = s.deps[j], s.deps[i]
}

func (s *depSorter) Less(i, j int) bool {
	return s.by(&s.deps[i], &s.deps[j])
}

var (
	Line = func(d1, d2 *Departure) bool {
		return d1.Line < d2.Line
	}
	Destination = func(d1, d2 *Departure) bool {
		return d1.Destination < d2.Destination
	}
	Eta = func(d1, d2 *Departure) bool {
		return d1.EtaNanoSec < d2.EtaNanoSec
	}
	EtaDesc = func(d1, d2 *Departure) bool {
		return Eta(d2, d1)
	}
)

type Departure struct {
	Line        string
	Destination string
	EtaNanoSec  time.Duration
}

type Backend interface {
	Setup()
	GetDepartures(station string) []Departure
}

type Frontend interface {
	Setup()
	RenderDepartures(station string, deps []Departure)
}

var (
	AllBackends  = make(map[string]Backend)
	AllFrontends = make(map[string]Frontend)
)
