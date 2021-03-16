// +build windows

package fs

import (
	"os"
)

func GetFileOwner(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return "", nil
}
