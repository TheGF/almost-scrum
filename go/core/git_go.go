package core

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func prepareMessage(commitInfo CommitInfo) string {
	var out bytes.Buffer

	out.WriteString(commitInfo.Header)
	out.WriteString("\n\n============\n")

	for task, comment := range commitInfo.Body {
		out.WriteString(task)
		out.WriteString("\n")
		out.WriteString(comment)
		out.WriteString("\n------------\n\n")
	}
	return out.String()
}

type GoGit struct {}

func (client GoGit) GetStatus(project *Project) (GitStatus, error) {
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
			if parts[1] == ProjectBoardsFolder && state.Worktree == git.Unmodified {
				gitStatus.AshFiles = append(gitStatus.AshFiles, name)
			} else if parts[1] == ProjectLibraryFolder && project.Config.IncludeLibInGit &&
				state.Worktree == git.Unmodified {
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

func GitPull(project *Project, username string) (string, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	r, err := git.PlainOpen(gitFolder)
	if err != nil {
		return "", err
	}
	logrus.Debugf("Open git repository %s", gitFolder)

	w, err := r.Worktree()
	if err != nil {
		return "", err
	}
	logrus.Debugf("Worktree successfully open")

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return "", err
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return "", err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	elapsed := time.Since(start)
	logrus.Infof("Pull completed in %s", elapsed)
	return commit.String(), nil
}

func getAuth(project *Project, user string) (transport.AuthMethod, error) {
	userInfo, err := GetUserInfo(project, user)
	if err != nil {
		return nil, err
	}

	credentials, found := userInfo.Credentials["GitUserPass"]
	if !found {
		logrus.Debugf("No Git username and password for user %s", user)
		return nil, nil
	}

	credentials, err = DecryptStringForProject(project, credentials)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(credentials, ":")
	if len(parts) == 2 {
		logrus.Debugf("Fount Git username and password for user %s", user)
		return &http.BasicAuth{
			Username: parts[0],
			Password: parts[1],
		}, nil
	} else {
		return nil, ErrNoFound
	}
}

 func GitPush(project *Project, user string) error {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	r, err := git.PlainOpen(gitFolder)
	if err != nil {
		return err
	}

	auth, err := getAuth(project, user)
	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{
		Auth: auth,
	})
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	logrus.Infof("Push completed in %s", elapsed)
	return nil
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
	logrus.Debugf("Open git repository %s", gitFolder)

	w, err := r.Worktree()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	logrus.Debugf("Worktree successfully open")

	message := prepareMessage(commitInfo)
	logrus.Debugf("Git message:\n%s", message)

	for _, file := range commitInfo.Files {
		if _, err = w.Add(file); err != nil {
			logrus.Warnf("Cannot add file %s to the commit: %v", file, err)
		} else {
			logrus.Debugf("Added file %s to commit", file)
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

type GitCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SetGitCredentials(project *Project, user string, gitCredentials GitCredentials) error {
	userInfo, err := GetUserInfo(project, user)
	if err != nil {
		return err
	}

	if gitCredentials.Password != "" {
		credentials := fmt.Sprintf("%s:%s", gitCredentials.Username, gitCredentials.Password)
		credentials, err := EncryptStringForProject(project, credentials)
		if err != nil {
			return err
		}
		userInfo.Credentials["GitUserPass"] = credentials

		if err := SetUserInfo(project, user, &userInfo); err != nil {
			return err
		}
	}
	return nil
}
