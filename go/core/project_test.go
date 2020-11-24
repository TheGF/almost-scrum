package core

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestFindProject(t *testing.T) {
	// _, err := FindProject("../test-data", "")
	// assert.Equal(t, err, ErrNoFound, "Found a project in test folder but none was expected")

}

func TestInitProject(t *testing.T) {
	_, err := InitProject("../test-data/my-scrum")
	assert.Nilf(t, err, "Cannot initialize project: %w", err)

	project, err := FindProject("..", "")
	assert.Nilf(t, err, "Cannot find expected project: %w", err)
	assert.DirExistsf(t, project.Path, "Expected project but none found: %w", err)

	err = ShredProject(project)
	assert.Nilf(t, err, "Cannot shread project: %w", err)

}
