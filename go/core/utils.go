package core

import (
	"archive/zip"
	"bufio"
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// FindFileUpwards looks for the folder where a file with the specified name is present.
func FindFileUpwards(path string, name string) (string, os.FileInfo) {
	path, _ = filepath.Abs(path)
	for parent := ""; ; path = parent {
		fileInfo, err := os.Stat(filepath.Join(path, name))
		if err == nil {
			logrus.Debugf("FindFileUpward - found %s in %s", name, path)
			return path, fileInfo
		}

		parent = filepath.Dir(path)
		if parent == path {
			break
		}
		logrus.Debugf("FindFileUpward -  trying parent %s", parent)
	}
	logrus.Debugf("FindFileUpward - no file '%s'", name)
	return "", nil
}

func IsErr(err error, msg string, args ...interface{}) bool {
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

func GetSystemUser() string {
	u, _ := user.Current()
	return u.Username
}

func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func UnzipFile(data []byte, destination string) error {
	var filenames []string

	reader := bytes.NewReader(data)
	r, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		path := filepath.Join(destination, f.Name)

		filenames = append(filenames, path)
		if f.FileInfo().IsDir() {

			// Creating a new Folder
			_ = os.MkdirAll(path, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		_ = outFile.Close()
		_ = rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	logrus.Debugf("Open browser at %s", url)
	return exec.Command(cmd, args...).Start()
}


func RunProgram(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func UseCommand(command string, input string, arg ...string) (output string, err error) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)
	cmd := exec.Command(command, arg...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err = cmd.Run(); err != nil {
		return buf.String(), err
	}
	return buf.String(), err
}

