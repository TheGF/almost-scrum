package core

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"path/filepath"
	"time"
)

type GitMessagePart struct {
	Task     string
	Progress []string
}

type CommitInfo struct {
	User   string
	Header string
	Body   []GitMessagePart
}

func prepareMessage(commitInfo CommitInfo) string {
	var out bytes.Buffer

	out.WriteString(commitInfo.Header)
	out.WriteString("\n")

	for _, part := range commitInfo.Body {
		out.WriteString(part.Task)
		out.WriteString("\n")
	}
	return out.String()
}

func GitStatus(project *Project) error {
	gitFolder := filepath.Dir(project.Path)

	repo, err := git.PlainOpen(gitFolder)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}

	print(status)
	return nil
}

func GitPull(project *Project) {

}

func GitCommit(project *Project, commitInfo CommitInfo) error {
	gitFolder := filepath.Dir(project.Path)
	boardsFolder := filepath.Join(project.Path, ProjectBoardsFolder)

	userInfo, err := GetUserInfo(project, commitInfo.User)
	if err != nil {
		return err
	}

	r, err := git.PlainOpen(gitFolder)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(boardsFolder)

	message := prepareMessage(commitInfo)
	fmt.Printf("Message %v:", message)

	if message == "" {
		_, err := w.Commit(message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  commitInfo.User,
				Email: userInfo.Email,
				When:  time.Now(),
			},
		})
		if err != nil {
			return nil
		}
//		logrus.Info("Commit %v", commit)
	}
	return nil
}
