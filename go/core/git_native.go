package core

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)



type GitNative struct{}

var statusRe = regexp.MustCompile(`\s*([?\w]+)\s+(.*)`)

func (client GitNative) GetStatus(project *Project) (GitStatus, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	out, err := RunCommand("git", "-C", gitFolder, "status", "--porcelain")
	if err != nil {
		return GitStatus{}, err
	}

	gitStatus := GitStatus{
		AshFiles:       make([]string, 0),
		StagedFiles:    make([]string, 0),
		UntrackedFiles: make([]string, 0),
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
			if parts[1] == ProjectBoardsFolder {
				gitStatus.AshFiles = append(gitStatus.AshFiles, name)
			} else if parts[1] == ProjectLibraryFolder && project.Config.IncludeLibInGit {
				gitStatus.AshFiles = append(gitStatus.AshFiles, name)
			}
		} else {
			switch change {
			case "D", "M", "AM":
				gitStatus.StagedFiles = append(gitStatus.StagedFiles, name)
			case "??":
				gitStatus.UntrackedFiles = append(gitStatus.UntrackedFiles, name)
			}
		}
	}

	elapsed := time.Since(start)
	logrus.Infof("Git Status completed in %s", elapsed)
	return gitStatus, nil
}

func (client GitNative) Pull(project *Project, user string) (string, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	out, err := RunCommand("git", "-C", gitFolder, "pull")
	if err != nil {
		return "",err
	}

	elapsed := time.Since(start)
	logrus.Infof("Pull completed in %s", elapsed)
	return out,nil
}

func (client GitNative) Push(project *Project, user string) (string, error) {
	start := time.Now()
	gitFolder := filepath.Dir(project.Path)

	out, err := RunCommand("git", "-C", gitFolder, "push")
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

	var args = []string{"-C", gitFolder, "add"}
	for _, file := range commitInfo.Files {
		args = append(args, file)
	}
	out, err := RunCommand("git", args...)
	logrus.Debugf("Git add: %s", out)
	if err != nil {
		return out, err
	}

	out, err = RunCommand("git", "-C", gitFolder, "commit",
		"-F", tmpFile.Name())
	logrus.Debugf("Git commit: %s", out)
	if err != nil {
		return out, err
	}

	elapsed := time.Since(start)
	logrus.Infof("Commit completed in %s. Output: %v", elapsed, out)
	return out, nil
}
