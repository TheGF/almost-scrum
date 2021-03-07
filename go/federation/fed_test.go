package federation

import (
	"almost-scrum/core"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getProject(t *testing.T) *core.Project{
	path := "/tmp/ash-test"
	_ = os.RemoveAll(path)
	_ = os.MkdirAll(path, 0755)
	core.UnzipProjectTemplates("/tmp/ash-test", []string{"file:/../test-data/test-template.zip"})
	project, err := core.OpenProject(path)
	assert.Nilf(t, err, "Cannot open project: %w", err)

	return project
}

func TestExport(t *testing.T) {
	project := getProject(t)
	Export(project, "marco", time.Now().AddDate(-1, 0, 0))
}

func TestImport(t *testing.T) {
	project := getProject(t)
	Import(project, time.Now().AddDate(0, -1, 0))

}