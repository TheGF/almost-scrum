package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var story = Story{
	Description: "Test of a story",
	Points:      12,
	Users:       []string{"mp", "mol"},
	Tasks: []Task{{
		Description: "Test of a task",
		Done:        true,
	},
		{
			Description: "Test of a task2",
			Done:        false,
		},
	},
	TimeEntries: []TimeEntry{},
	Attachments: []string{},
}

func TestStoreList(t *testing.T) {
	store := Store{Path: ".."}
	list := List(store)
	assert.GreaterOrEqual(t, len(list), 1)
}

func TestStoreSet(t *testing.T) {

	store := Store{Path: "../test-data"}
	assert.Nilf(t, Set(store, "1.Hello.story", &story), "Cannot write to store")
}

func TestStoreGet(t *testing.T) {
	store := Store{Path: "../test-data"}
	s, _ := Get(store, "1.Hello.story")
	assert.Equal(t, story, s, "Mismatch in story read")
}

func BenchmarkGet(b *testing.B) {
	store := Store{Path: "../test-data"}
	for i := 0; i < b.N; i++ {
		Get(store, "1.Hello.story")
	}

}
