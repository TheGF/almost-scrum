package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gitClient = GitNative{}

func TestStatus(t *testing.T) {
	project, err := OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitStatus, err := gitClient.GetStatus(project)
	t.Logf("Gitfiles %v", gitStatus)
}


func TestCommit(t *testing.T) {
	project, err := OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	commitInfo := CommitInfo{
		User:     "mp",
		Header:   "This is just a test",
		Body: map[string]string{
			"task 1": "a comment",
			"task 2": "another comment",
			"task 3": "final comment",
		},
		Files:    []string{},
	}

	hash, err := gitClient.Commit(project, commitInfo)
	t.Logf("GitCommit %v", hash)
}

func TestPush(t *testing.T) {
	project, err := OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitClient.Push(project, GetSystemUser())
}

func TestSetGitCredentials(t *testing.T) {
	project, err := OpenProject("../../.ash")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitCredentials := GitSettings{
		Username: "mp",
		Password: "Mariposa83$",//"122abd851b139cb1f73a1f1ded40c6047dd9c3d3",
	}

	SetGitSettings(project, GetSystemUser(), gitCredentials)
}