// +build darwin linux

package library

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func getFileOwner(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	stat := info.Sys().(*syscall.Stat_t)
	user, err := user.LookupId(strconv.Itoa(int(stat.Uid)))
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
