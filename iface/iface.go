package iface

import (
	"time"
)

type JsonDuration time.Duration

func (d JsonDuration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(d).String() + "\""), nil
}

type Departure struct {
	Line        string
	Destination string
	Eta         JsonDuration `json:"string"`
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
