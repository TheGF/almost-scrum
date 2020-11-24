package core

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// UserInfo contains information about a user.
type UserInfo struct {
	Email string `json:"email"`
	Icon  []byte `json:"icon"`
}

// GetUserList returns the project users
func GetUserList(project Project) []string {
	path := filepath.Join(project.Path, ProjectUsersFolder)
	file, _ := os.Open(path)
	users, _ := file.Readdirnames(0)
	return users
}

// GetUserInfo returns information about the specified user
func GetUserInfo(project Project, user string) (userInfo UserInfo, err error) {
	path := filepath.Join(project.Path, ProjectUsersFolder, user)
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log.Infof("Invalid file %s: %v", path, err)
		return
	}

	err = yaml.Unmarshal(d, &user)
	if err != nil {
		log.Infof("Invalid file %s: %v", path, err)
		return
	}
	return
}

//SetUserInfo saves the user info
func SetUserInfo(project Project, user string, userInfo *UserInfo) (err error) {
	d, err := yaml.Marshal(userInfo)
	if err != nil {
		log.Errorf("Cannot marshal info for user %s: %v", user, err)
		return
	}
	path := filepath.Join(project.Path, ProjectUsersFolder, user)
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		log.Errorf("Cannot save info for user %s: %v", user, err)
		return
	}
	log.Infof("Story saved to %s", path)
	return
}
