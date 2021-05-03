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
	project, err := core.OpenProject("../../.scrum-to-go")
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

func TestNextVersion(t *testing.T) {
	res, err := getNextVersion("~0.2", false)
	assert.Nil(t, err)
	assert.Equal(t, res, "~0.3")

	res, err = getNextVersion("", false)
	assert.Nil(t, err)
	assert.Equal(t, res, "~0.1")

	res, err = getNextVersion("~0.2", true)
	assert.Nil(t, err)
	assert.Equal(t, res, "~0.2.1")

}

func TestExport(t *testing.T) {
	html, _ := ExportMarkdownToHTML("../../web/src/help/portal.md", "/tmp/export.html",
		"../../web/public")

	print(html)
}
