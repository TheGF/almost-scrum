package core

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// FindFileUpwards looks for the folder where a file with the specified name is present.
func FindFileUpwards(path string, name string) (string, os.FileInfo) {
	path, _ = filepath.Abs(path)
	for parent := "";; path = parent {
		fileInfo, err := os.Stat(filepath.Join(path, name))
		if err == nil {
			logrus.Debugf("FindFileUpward - found %s in %s", name, path)
			return path, fileInfo
		}

		parent = filepath.Dir(path)
		if  parent == path {
			break
		}
		logrus.Debugf("FindFileUpward -  trying parent %s", parent)
	}
	logrus.Debugf("FindFileUpward - no file '%s'", name)
	return "", nil
}

func GetCurrentUser() string {
	user, err := user.Current()
	if err != nil {
		logrus.Fatalf("Cannot get current user: %v", err)
		os.Exit(1)
	}

	return user.Username
}