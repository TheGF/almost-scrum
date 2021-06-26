package fs

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
)


var versionsRegex = regexp.MustCompile(`(.*?)(~(\d+\.)+\d+)?(\.\w*)?$`)

func ParsePath(path string) (dir string, prefix string, version string, ext string, err error) {
	dir = filepath.Dir(path)
	name := filepath.Base(path)
	match := versionsRegex.FindStringSubmatch(name)
	if len(match) != 5 {
		err = os.ErrInvalid
		logrus.Errorf("Cannot parse %s: %v", path, err)
		return
	}

	prefix = match[1]
	version = match[2]
	ext = match[4]
	err = nil

	return
}
