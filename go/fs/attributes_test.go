// +build skip


package fs

import (
	"almost-scrum/core"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestGetSetExtendedAttr(t *testing.T) {
	path := "/tmp/text.txt"
	f, err := os.Create(path)
	assert.Nil(t, err)
	f.Close()

	attr, err := GetExtendedAttr(filepath.Dir(path), filepath.Base(path))
	cacheWg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, attr.Owner, core.GetSystemUser())

}
