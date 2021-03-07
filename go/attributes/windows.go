// +build windows

package attributes

import (
	"os"
)

func getFileOwner(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return "", nil
}
