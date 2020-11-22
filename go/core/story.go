// Package core provides basic functionality for Almost Scrum
package core

// Story contains attributes that define a story
type Story struct {
	Description string   `json:"description"`
	Points      int      `json:"points"`
	Attachments []string `json:"attachments"`
	Users       []string `json:"users"`
}
