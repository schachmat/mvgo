package iface

import (
	"time"
)

type Departure struct {
	Line        string
	Destination string
	Eta         time.Duration
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
