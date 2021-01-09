package core

import (
	"bytes"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GitFileState int

const (
	GitUntracked GitFileState = iota
	GitReady
	GitMerged
)

type GitFile struct {
	Path  string
	State GitFileState
}

type GitStatus struct {
	AshFiles       []string
	StagedFiles    []string
	UntrackedFiles []string
}

type CommitInfo struct {
	User     string
	Header   string
	Comments map[string]string
	Files    []string
}

func prepareMessage(commitInfo CommitInfo) string {
	var out bytes.Buffer

	out.WriteString(commitInfo.Header)
	out.WriteString("\n\n============\n")

	for task, comment := range commitInfo.Comments {
		out.WriteString(task)
		out.WriteString("\n")
		out.WriteString(comment)
		out.WriteString("\n------------\n\n")
	}
	return out.String()
}

func GetGitStatus(project *Project) (GitStatus, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	repo, err := git.PlainOpen(gitFolder)
	if err != nil {
		return GitStatus{}, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return GitStatus{}, err
	}

	status, err := worktree.Status()
	if err != nil {
		return GitStatus{}, err
	}

	gitStatus := GitStatus{
		AshFiles:       make([]string, 0),
		StagedFiles:    make([]string, 0),
		UntrackedFiles: make([]string, 0),
	}
	for name, state := range status {
		parts := strings.Split(name, string(os.PathSeparator))
		if len(parts) == 0 {
			continue
		}
		if parts[0] == ProjectFolder {
			if parts[1] == ProjectBoardsFolder {
				gitStatus.AshFiles = append(gitStatus.AshFiles, name)
			} else if parts[1] == ProjectLibraryFolder && project.Config.IncludeLibInGit {
				gitStatus.AshFiles = append(gitStatus.AshFiles, name)
			}
		} else {
			switch state.Worktree {
			case git.Modified, git.Added, git.Deleted, git.Renamed:
				gitStatus.StagedFiles = append(gitStatus.StagedFiles, name)
			case git.Untracked:
				gitStatus.UntrackedFiles = append(gitStatus.UntrackedFiles, name)
			}
		}
	}

	elapsed := time.Since(start)
	logrus.Infof("Git Status completed in %s", elapsed)
	return gitStatus, nil
}

func GitPull(project *Project) {

}

func GitCommit(project *Project, commitInfo CommitInfo) (plumbing.Hash, error) {
	start := time.Now()

	gitFolder := filepath.Dir(project.Path)
	userInfo, err := GetUserInfo(project, commitInfo.User)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	r, err := git.PlainOpen(gitFolder)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	w, err := r.Worktree()
	if err != nil {
		return plumbing.ZeroHash, err
	}

	message := prepareMessage(commitInfo)

	for _, file := range commitInfo.Files {
		if _, err = w.Add(file); err != nil {
			logrus.Warnf("Cannot add file %s to the commit: %v", file, err)
		} else {
			logrus.Infof("Added file %s to commit", file)
		}
	}

	hash, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  commitInfo.User,
			Email: userInfo.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		logrus.Warnf("Cannot complete commit: %v", err)
		return plumbing.ZeroHash, err
	}
	elapsed := time.Since(start)
	logrus.Infof("Commit completed in %s. Hash: %v", elapsed, hash)
	return hash, nil
}
