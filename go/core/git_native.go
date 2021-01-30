package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type GitNative struct{}

var statusRe = regexp.MustCompile(`\s*([?\w])[?\w]?\s+(.*)`)

func (client GitNative) GetStatus(project *Project) (GitStatus, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	out, err := UseCommand("git", "", "-C", gitFolder, "status", "--porcelain")
	if err != nil {
		return GitStatus{}, err
	}

	gitStatus := GitStatus{
		AshFiles: make([]string, 0),
		Files: make(map[string]GitChange),
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		parts := statusRe.FindStringSubmatch(line)
		if len(parts) < 3 {
			continue
		}

		change := parts[1]
		name := strings.Trim(parts[2], `"`)
		parts = strings.Split(name, string(os.PathSeparator))
		if len(parts) == 0 {
			continue
		}

		if parts[0] == ProjectFolder {
			if parts[1] == ProjectBoardsFolder || (project.Config.IncludeLibInGit && parts[1] != ProjectLibraryFolder) {
				gitStatus.AshFiles = append(gitStatus.AshFiles, name)
			}
		} else {
			gitStatus.Files[name] = GitChange(change)
		}
	}

	elapsed := time.Since(start)
	logrus.Infof("Git Status completed in %s", elapsed)
	return gitStatus, nil
}

//func escapeCredentials(value string) string {
//	tokens := []string{"@", "#", "$", "%", "^", "*", `"`, `'`, "`", `\`}
//
//	for _, token := range tokens {
//		value = strings.ReplaceAll(value, token, `\`+token)
//	}
//	return value
//}

func getGitCredentialAsInput(project *Project, user string) string {
	username, password, err := GetGitCredentials(project, user)
	if err != nil {
		return ""
	}

	username = url.PathEscape(username)
	password = url.PathEscape(password)

	return fmt.Sprintf("%s:%s", username, password)
}

var remoteLocRegex = regexp.MustCompile(`origin\s?(http.*)`)

type gitAction string

const (
	FetchAction = "fetch"
	PushAction  = "push"
)

func getGitRemoteLoc(gitFolder string, action gitAction) (string, error) {
	output, err := UseCommand("git", "", "-C", gitFolder, "remote", "-v")
	if err != nil {
		return "", err
	}

	suffix := fmt.Sprintf(" (%s)", action)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		match := remoteLocRegex.FindStringSubmatch(line)
		if len(match) == 2 {
			if strings.HasSuffix(match[1], suffix) {
				return match[1][0 : len(match[1])-len(suffix)], nil
			}
		}
	}
	return "", ErrNoFound
}

func getRemoteWithCredentials(gitFolder string, action gitAction, credentials string) string {
	remote, err := getGitRemoteLoc(gitFolder, action)
	if err != nil {
		return ""
	}

	parts := strings.Split(remote, "://")
	if len(parts) == 2 {
		return fmt.Sprintf("%s://%s@%s", parts[0], credentials, parts[1])
	} else {
		return ""
	}
}

func (client GitNative) Pull(project *Project, user string) (string, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	credentials := getGitCredentialAsInput(project, user)
	remote := getRemoteWithCredentials(gitFolder, FetchAction, credentials)
	output, err := UseCommand("git", "", "-C", gitFolder, "pull", remote)
	if err != nil {
		return "", err
	}

	elapsed := time.Since(start)
	logrus.Infof("Pull completed in %s", elapsed)
	return output, nil
}

func (client GitNative) Push(project *Project, user string) (string, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	credentials := getGitCredentialAsInput(project, user)
	remote := getRemoteWithCredentials(gitFolder, PushAction, credentials)
	out, err := UseCommand("git", "", "-C", gitFolder, "push", remote)
	if err != nil {
		return out, err
	}

	elapsed := time.Since(start)
	logrus.Infof("Push completed in %s", elapsed)
	return out, nil
}

func (client GitNative) Commit(project *Project, commitInfo CommitInfo) (string, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return "", err
	}

	defer os.Remove(tmpFile.Name())
	message := prepareGitMessage(commitInfo)
	logrus.Debugf("Git message:\n%s", message)

	if _, err = tmpFile.WriteString(message); err != nil {
		return "", err
	}
	_ = tmpFile.Close()

	gitStatus, err := client.GetStatus(project)
	if err != nil {
		return "", err
	}

	addArgs := []string {"-C", gitFolder, "add"}
	rmArgs := []string {"-C", gitFolder, "add", "-u"}
	for _, file := range commitInfo.Files {
		change, found := gitStatus.Files[file]
		if found && change == GitDeleted {
			parts := []string{gitFolder}
			parts = append(parts, strings.Split(file, "/")...)
			p := filepath.Join(parts...)
			if err := ioutil.WriteFile(p, []byte{}, 0666); err == nil {
				rmArgs = append(rmArgs, file)
			}
		} else {
			addArgs = append(addArgs, file)
		}
	}
	if len(addArgs) > 3 {
		out, err := UseCommand("git", "", addArgs...)
		logrus.Debugf("Git add: %s", out)
		if err != nil {
			return out, err
		}
	}
	if len(rmArgs) > 3 {
		out, err := UseCommand("git", "", rmArgs...)
		logrus.Debugf("Git add: %s", out)
		if err != nil {
			return out, err
		}
	}

	out, err := UseCommand("git", "", "-C", gitFolder, "commit",
		"-F", tmpFile.Name())
	logrus.Debugf("Git commit: %s", out)
	if err != nil {
		return out, err
	}

	elapsed := time.Since(start)
	logrus.Infof("Commit completed in %s. Output: %v", elapsed, out)
	return out, nil
}
