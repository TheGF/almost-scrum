package core

import (
	"almost-scrum/assets"
	"path/filepath"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestFindFileUpwards(t *testing.T) {
	if runtime.GOOS == "windows" {
		path, _ := FindFileUpwards("c:\\windows\\system32", "win.ini")
		assert.Equal(t, "c:\\windows", path)
	} else {
		path, _ := FindFileUpwards("/usr/bin", "cat")
		assert.Equal(t, "/usr/bin", path)
	}

	path, _ := FindFileUpwards(".", "main.go")
	assert.Equal(t, "go", filepath.Base(path))

	path, _ = FindFileUpwards(".", ".git")
	assert.Equal(t, "almost-scrum", filepath.Base(path))

}

func TestUnzipAsset(t *testing.T) {
	data, _ := assets.Asset("assets/one-week-scrum.zip")
	UnzipFile(data, "/tmp/")

}
