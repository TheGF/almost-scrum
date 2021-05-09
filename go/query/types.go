package query

import (
	"almost-scrum/core"
	"time"
)

type TaskRef struct {
	Board   string    `json:"board"`
	Name    string    `json:"name"`
	ModTime time.Time `json:"modTime"`
	Task    core.Task `json:"task"`
}

type Select struct {
	Description bool `json:"description"`
	Properties  bool `json:"properties"`
	Parts       bool `json:"parts"`
	Files       bool `json:"files"`
}

type WhereType struct {
	Is               []string `json:"is"`
	HasPropertiesAll []string `json:"hasProperties"`
}

type WhereProperty struct {
	Name          string   `json:"property"`
	ValueIsAnyOf  []string `json:"is"`
	ValueIsNoneOf []string `json:"isNot"`
}

type Query struct {
	Select          Select          `json:"select"`
	WhereTypes      []WhereType     `json:"whereTypes"`
	WhereProperties []WhereProperty `json:"whereProperties"`
	WhereBoardIs    []string        `json:"whereBoardIs"`
}

type State map[string]*TaskRef

//type State struct {
//	Tasks   map[string]*TaskRef
//	Tracked map[string]time.Time
//}
