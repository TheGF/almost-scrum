package chat

import (
	"almost-scrum/core"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

func AddMessage(project *core.Project, user string, reader io.ReadCloser) error {
	folder := filepath.Join(project.Path, core.ProjectChatFolder)
	filename := filepath.Join(folder, fmt.Sprintf("%s.%x.bin", user, time.Now().Unix()))
	writer, err := os.Create(filename)
	if err != nil {
		logrus.Warnf("Cannot open file %s for writing: %v", filename, err)
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		logrus.Warnf("cannot write file %s in library: %v", filename, err)
		return err
	}
	logrus.Debugf("successfully set message %s in chat", filename)
	return err
}
