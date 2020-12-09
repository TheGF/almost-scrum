// +build linux

package cli

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func openEditor(path string) {
	var editor string = config.Editor
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Warnf("Cannot open editor '%s': %v", editor, err)
		color.Red("Something went wrong while trying to open editor '%s'", editor)
		os.Exit(1)
	}
}
