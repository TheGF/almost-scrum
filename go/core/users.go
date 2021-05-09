package core

import (
	"almost-scrum/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Todo struct {
	Id     int       `json:"id"`
	Action string    `json:"action"`
	Eta    time.Time `json:"eta"`
	Done   bool      `json:"done"`
}

// UserInfo contains information about a user.
type UserInfo struct {
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Office      string            `json:"office"`
	Icon        []byte            `json:"icon"`
	Todo        []Todo            `json:"todo"`
	Credentials map[string]string `json:"credentials"`
}

// GetUserList returns the project users
func GetUserList(project *Project) []string {
	path := filepath.Join(project.Path, ProjectUsersFolder)
	file, _ := os.Open(path)
	names, _ := file.Readdirnames(0)
	users := make([]string, 0, len(names))
	for _, name := range names {
		ext := filepath.Ext(name)
		if ext == ".yaml" {
			users = append(users, name[0:len(name)-len(ext)])
		}
	}
	logrus.Debugf("Users in project %s: %v", project.Path, users)
	return users
}

// GetUserInfo returns information about the specified user
func GetUserInfo(project *Project, user string) (userInfo UserInfo, err error) {
	user = strings.ToLower(user)
	path := filepath.Join(project.Path, ProjectUsersFolder, user+".yaml")
	if err = fs.ReadYaml(path, &userInfo); err != nil {
		logrus.Warnf("Invalid file %s: %v", path, err)
		return
	}
	return
}

//SetUserInfo saves the user info
func DelUserInfo(project *Project, user string) (err error) {
	user = strings.ToLower(user)
	path := filepath.Join(project.Path, ProjectUsersFolder, user+".yaml")
	if err = os.Remove(path); err != nil {
		logrus.Errorf("Cannot remove info for user %s: %v", user, err)
		return
	}
	logrus.Infof("User %s removed", user)
	return
}

//SetUserInfo saves the user info
func SetUserInfo(project *Project, user string, userInfo *UserInfo) (err error) {
	user = strings.ToLower(user)
	path := filepath.Join(project.Path, ProjectUsersFolder, user+".yaml")
	return fs.WriteYaml(path, userInfo)
}
