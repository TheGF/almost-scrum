// Package core provides basic functionality for Almost Scrum
package core

import "time"

// Task to complete the story and its status
type Task struct {
	Description string `json:"description" yaml:"description"`
	Done        bool   `json:"done" yaml:"done"`
}

// TimeEntry is time used by a user on an activity on a specific day
type TimeEntry struct {
	User  string    `json:"user" yaml:"user"`
	Date  time.Time `json:"date" yaml:"date"`
	Hours int       `json:"hours" yaml:"hours"`
}

// Story contains attributes that define a story
type Story struct {
	Description string      `json:"description" yaml:"description"`
	Points      int         `json:"points" yaml:"points"`
	Tasks       []Task      `json:"tasks" yaml:"tasks"`
	TimeEntries []TimeEntry `json:"timeEntries" yaml:"timeEntries"`
	Attachments []string    `json:"attachments" yaml:"attachments"`
	Users       []string    `json:"users" yaml:"users,flow"`
}
