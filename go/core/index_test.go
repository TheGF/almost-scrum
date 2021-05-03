package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestReindex(t *testing.T) {
	project, err := OpenProject("../../.scrum-to-go")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	err = ReIndex(project)
	assert.NotNilf(t, err, "Cannot reindex project: %w", err)
}

func TestReindexTask(t *testing.T) {
	project, err := OpenProject("../test-data")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	err = ReIndex(project)
	assert.NotNilf(t, err, "Cannot reindex project: %w", err)
}


func TestSuggestKeys(t *testing.T) {
	project, err := OpenProject("../../.scrum-to-go")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	ReIndex(project)

	out := SuggestKeys(project, "@", 10)
	println(out)
}