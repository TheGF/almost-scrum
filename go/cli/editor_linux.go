// +build linux

package cli

import (
	"os"
	"os/exec"
)

func openExternalEditor(editor string, path string) {
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	abortIf(err, "")
}
