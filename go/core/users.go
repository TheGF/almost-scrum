package core

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)



// UserInfo contains information about a user.
type UserInfo struct {
	Email string `json:"email"`
	Icon  []byte `json:"icon"`
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

	return users
}

// GetUserInfo returns information about the specified user
func GetUserInfo(project *Project, user string) (userInfo UserInfo, err error) {
	path := filepath.Join(project.Path, ProjectUsersFolder, user+".yaml")
	if err = ReadYaml(path, &userInfo); err != nil {
		log.Warnf("Invalid file %s: %v", path, err)
		return
	}
	return
}

//SetUserInfo saves the user info
func DelUserInfo(project *Project, user string) (err error) {
	path := filepath.Join(project.Path, ProjectUsersFolder, user+".yaml")
	if err = os.Remove(path); err != nil {
		log.Errorf("Cannot remove info for user %s: %v", user, err)
		return
	}
	log.Infof("User %s removed", user)
	return
}


//SetUserInfo saves the user info
func SetUserInfo(project *Project, user string, userInfo *UserInfo) (err error) {
	path := filepath.Join(project.Path, ProjectUsersFolder, user+".yaml")
	return WriteYaml(path, userInfo)
}
