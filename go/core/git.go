package core

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
	"strings"
)

type GitStatus struct {
	AshFiles       []string `json:"ashFiles"`
	StagedFiles    []string `json:"stagedFiles"`
	UntrackedFiles []string `json:"untrackedFiles"`
}

type CommitInfo struct {
	User   string            `json:"user"`
	Header string            `json:"header"`
	Body   map[string]string `json:"body"`
	Files  []string          `json:"files"`
}

type GitClient interface {
 	GetStatus(project *Project) (GitStatus, error)
	Pull(project *Project, user string) (string, error)
	Push(project *Project, user string) error
	Commit(project *Project, commitInfo CommitInfo) (string, error)
}


func prepareGitMessage(commitInfo CommitInfo) string {
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

type GitCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HasGitNative() bool {
	out, _, err := RunCommand("git", "--version")
	return err == nil && strings.HasPrefix(out, "git")
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
