package fed

import (
	"almost-scrum/fed/transport"
	"errors"
	"sync"
	"time"
)

var (
	ErrFedNoTimeServer = errors.New("no time server")
	ErrFedTimeMismatch = errors.New("time mismatch")
	ErrFedCorrupted    = errors.New("corrupted file")
	ErrNoConnections   = errors.New("no available connection to the federation")
)

const HeaderFile = "ash-header.json"

type Header struct {
	ProjectID string    `json:"projectId"`
	Host      string    `json:"host"`
	Hostname  string    `json:"hostname"`
	Time      time.Time `json:"time"`
	User      string    `json:"user"`
}

type Stat struct {
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
	Push     int   `json:"push"`
	Pull     int   `json:"pull"`
}

type Signal struct {
	local      string
	remote     string
	config     *Config
	locs       sync.Map
	exchanges  map[transport.Exchange]bool
	stat       map[transport.Exchange]Stat
	lastExport time.Time
	reconnect  *time.Ticker
	inUse      sync.WaitGroup
}

var (
	states = map[string]*Signal{}
)
