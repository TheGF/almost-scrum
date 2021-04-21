package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"errors"
	"github.com/patrickmn/go-cache"
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
	Version  string    `json:"version"`
	Host     string    `json:"host"`
	Hostname string    `json:"hostname"`
	Time     time.Time `json:"time"`
	User     string    `json:"user"`
}

type Throughput struct {
	Upload   int64   `json:"upload"`
	Download int64   `json:"download"`
}

type Transfer struct {
	Exchange string        `json:"exchange"`
	Locs     []string      `json:"locs"`
	Size     int64         `json:"size"`
	Error    error         `json:"error"`
	Issues   []error       `json:"issues"`
	Elapsed  time.Duration `json:"elapsed"`
}

type Connection struct {
	local      string
	remote     string
	config     *Config
	locs       sync.Map
	exchanges  map[transport.Exchange]bool
	throughput map[string]*Throughput
	exports    map[string]time.Time
	reconnect  *time.Ticker
	mutex      sync.Mutex
	checkTime  error
}

var (
	connections = cache.New(5*time.Minute, 10*time.Minute)
)

type syncItem struct {
	folders        []string
	includePrivate bool
	neverDelete    bool
	prefix         string
}

var syncItems = []syncItem{
	{
		folders:        []string{core.ProjectBoardsFolder},
		includePrivate: true,
		prefix:         "bo",
	},
	{
		folders: []string{
			core.ProjectArchiveFolder,
			core.ProjectLibraryFolder,
		},
		includePrivate: false,
		prefix:         "li",
	},
	{
		folders:        []string{core.ProjectModelsFolder},
		includePrivate: true,
		prefix:         "sy",
	},
}
