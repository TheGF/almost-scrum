package core

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
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

func IsErr(err error, msg string, args ... interface{}) bool {
	if err != nil {
		if msg == "" {
			logrus.Warnf("Unexpected error: %v", err)
		} else {
			args = append(args, err)
			msg = msg + ":%v"
			logrus.Warnf(msg, args...)
		}
		return true
	} else {
		return false
	}
}
