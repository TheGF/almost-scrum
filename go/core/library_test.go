package core

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestGetLibraryItems(t *testing.T) {
	project, err := OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	files := []string{"Richiesta ricongiungimento.pdf"}
	items, err := GetLibraryItems(project, files)
	assert.Nil(t, err)
	assert.NotNil(t, items)

}
