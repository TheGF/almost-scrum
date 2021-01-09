package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestStatus(t *testing.T) {
	project, err := OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitStatus, err := GetGitStatus(project)
	t.Logf("Gitfiles %v", gitStatus)
}
