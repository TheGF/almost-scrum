package core

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type GitChange string

const (
	GitModified  GitChange = "M"
	GitAdd                 = "A"
	GitDeleted             = "D"
	GitRenamed             = "R"
	GitCopied              = "C"
	GitUntracked           = "?"
)

type GitStatus struct {
	AshFiles []string             `json:"ashFiles"`
	Files    map[string]GitChange `json:"files"`
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
	Push(project *Project, user string) (string, error)
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

type GitSettings struct {
	UseGitNative bool   `json:"useGitNative"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func HasGitNative() bool {
	out, err := UseCommand("git", "--version")
	return err == nil && strings.HasPrefix(out, "git")
}

func GetGitCredentials(project *Project, user string) (username string, password string, err error) {
	userInfo, err := GetUserInfo(project, user)
	if err != nil {
		return "", "", err
	}

	credentials, ok := userInfo.Credentials["GitUserPass"]
	if ok {
		credentials, _ = DecryptStringForProject(project, credentials)
		idx := strings.Index(credentials, ":")
		return credentials[0:idx], credentials[1+idx:], nil
	} else {
		return "", "", ErrNoFound
	}

}

func SetGitCredentials(project *Project, user string, gitUsername string, gitPassword string) error {
	userInfo, err := GetUserInfo(project, user)
	if err != nil {
		return err
	}

	credentials := fmt.Sprintf("%s:%s", gitUsername, gitPassword)
	credentials, err = EncryptStringForProject(project, credentials)
	if err != nil {
		return err
	}
	userInfo.Credentials["GitUserPass"] = credentials

	if err := SetUserInfo(project, user, &userInfo); err != nil {
		return err
	}
	return nil
}

func GetGitSettings(project *Project, user string) (GitSettings, error) {
	username, _, _ := GetGitCredentials(project, user)

	return GitSettings{
		UseGitNative: project.Config.UseGitNative,
		Username:     username,
		Password:     "",
	}, nil
}

func SetGitSettings(project *Project, user string, gitSettings GitSettings) error {
	project.Config.UseGitNative = gitSettings.UseGitNative
	if err := WriteProjectConfig(project.Path, &project.Config); err != nil {
		return err
	}

	if gitSettings.Password != "" {
		if err := SetGitCredentials(project, user, gitSettings.Username, gitSettings.Password); err != nil {
			return err
		}
	}
	return nil
}

var mergeConflictRegex = regexp.MustCompile(`(?s)(.*)<<<<<<< HEAD(.*)=======(.*)>>>>>>> (\w+)`)

func FindGitConflict(input string) string {
	match := mergeConflictRegex.FindStringSubmatch(input)
	if len(match) == 0 {
		return ""
	}
	return match[4]
}

type GitSolveConflict int

const (
	GitSolveConflictWithHead GitSolveConflict = iota
	GitSolveConflictWithRemote
)

const (
	mergeStart = "<<<<<<< HEAD"
	mergeSplit = "======="
	mergeEnd   = ">>>>>>>"
)

const (
	mergeCommon int = iota
	mergeHead
	mergeRemote
)

func ResolveGitConflict(input string, gitSolveConflict GitSolveConflict) string {
	var b bytes.Buffer
	var state int

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if line == mergeStart {
			state = mergeHead
			continue
		}
		if line == mergeSplit {
			state = mergeRemote
			continue
		}
		if strings.HasPrefix(line, mergeEnd) {
			state = mergeCommon
			continue
		}

		switch state {
		case mergeCommon:
			b.WriteString(line)
		case mergeHead:
			if gitSolveConflict == GitSolveConflictWithHead {
				b.WriteString(line)
			}
		case mergeRemote:
			if gitSolveConflict == GitSolveConflictWithRemote {
				b.WriteString(line)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}