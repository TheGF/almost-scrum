package library

import (
	"almost-scrum/core"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestGetLibraryItems(t *testing.T) {
	project, err := core.OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	files := []string{"Richiesta ricongiungimento.pdf"}
	items, err := GetItems(project, files)
	assert.Nil(t, err)
	assert.NotNil(t, items)

}
func TestGetPreviousVersions(t *testing.T) {
	project, err := core.OpenProject("/tmp/Hello")

	items, err := GetPreviousVersions(project, "/architecture/network/SpecK-v0.2.docx")
	assert.Nil(t, err)
	assert.NotNil(t, items)

}
