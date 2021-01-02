package core

import (
	"github.com/go-git/go-git/v5"
	"path/filepath"
)

func GitCommit(project Project) error {
	gitFolder := filepath.Dir(project.Path)

	r, err := git.PlainOpen(gitFolder)
	if err != nil {
		return err
	}

	_, err = r.Worktree()
	if err != nil {
		return err
	}
	return nil
}