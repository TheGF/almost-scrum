package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var gitClient = GoGit{} //GitNative{}

func TestStatus(t *testing.T) {
	project, err := OpenProject("../../.scrum-to-go")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitStatus, err := gitClient.GetStatus(project)
	t.Logf("Gitfiles %v", gitStatus)
}

func TestClone(t *testing.T) {
	out, err := gitClient.Clone("https://github.com/TheGF/AI.git", "/tmp")
	assert.NotNilf(t, err, "Cannot clone: %w", err)

	print(out)
}

func TestCommit(t *testing.T) {
	project, err := OpenProject("../../.scrum-to-go")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	status, err := gitClient.GetStatus(project)

	commitInfo := CommitInfo{
		User:   "mp",
		Header: "This is just a test",
		Body: map[string]string{
			"task 1": "a comment",
			"task 2": "another comment",
			"task 3": "final comment",
		},
		Files: status.AshFiles,
	}

	hash, err := gitClient.Commit(project, commitInfo)
	t.Logf("GitCommit %v", hash)
}

func TestPush(t *testing.T) {
	project, err := OpenProject("../../.scrum-to-go")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitClient.Push(project, GetSystemUser())
}

func TestSetGitCredentials(t *testing.T) {
	project, err := OpenProject("../../.scrum-to-go")
	assert.NotNilf(t, err, "Cannot open project: %w", err)

	gitCredentials := GitSettings{
		UseGitNative: true,
		Username:     "TheGF",
		Password:     "Mariposa83$",
	}

	SetGitSettings(project, GetSystemUser(), gitCredentials)
}

func TestGitConflict(t *testing.T) {
const input = `
Replace with the task description
Just doing some random changes to test conflicts
<<<<<<< HEAD
Ok more changes needed.

### Properties
=======
And more changes
### Properties
- Points: 52
>>>>>>> f975cd6da9c8108c7a08ef6606bc97be5fcbac55
- Owner: @mp
- Points: 5
- Status: #Start
### Progress
### Locs
`

	id := FindGitConflict(input)
	assert.NotEmpty(t, id)

	solved := ResolveGitConflict(input, GitSolveConflictWithHead)
	assert.NotEmpty(t, solved)


}
