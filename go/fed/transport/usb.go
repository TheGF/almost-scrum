package transport

import (
	"almost-scrum/core"
	"errors"
	"fmt"
	usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type USBConfig struct {
	Name     string `json:"name" yaml:"name"`
}

type USBExchange struct {
	config *USBConfig
	remote string
	local  string
}

func GetUSBExchanges(configs ...USBConfig) []Exchange {
	var exchanges []Exchange
	for i := range configs {
		exchanges = append(exchanges, &USBExchange{
			config: &configs[i],
			remote: "",
			local:  "",
		})
	}
	return exchanges
}

func (exchange *USBExchange) ID() string {
	return fmt.Sprintf("ftp-%s", exchange.config.Name)
}

func getMounts() ([]string, error) {
	var mounts []string
	if drives, err := usbdrivedetector.Detect(); err == nil {
		for _, d := range drives {
			folder := filepath.Join(d, "almost-scrum")
			stat, err := os.Stat(folder)
			if err == nil && stat.IsDir() {
				mounts = append(mounts, folder)
			}
		}
		return mounts, nil
	} else {
		return nil, err
	}
}

var ErrNoUSBMounts = errors.New("no suitable USB devices")

func (exchange *USBExchange) Connect(remoteRoot, localRoot string) (UpdatesCh, error) {
	exchange.Disconnect()

	exchange.local = localRoot
	exchange.remote = remoteRoot

	mounts, err := getMounts()
	if err != nil {
		return nil, err
	} else if len(mounts) == 0 {
		return nil, ErrNoUSBMounts
	}

	return nil, nil
}

func (exchange *USBExchange) Disconnect() {
}

func copyFile(source string, file string, dest string) (int64, error) {
	pr := filepath.Join(source, file)

	stat, err := os.Stat(pr)
	if err != nil {
		logrus.Warnf("cannot open %s for USB Media sync", pr)
		return 0, err
	}

	r, err := os.Open(pr)
	if err != nil {
		logrus.Warnf("cannot open %s for USB Media sync", pr)
		return 0, err
	}
	defer r.Close()

	pw := filepath.Join(dest, file)
	_ = os.MkdirAll(filepath.Dir(pw), 0755)
	w, err := os.Create(pw)
	if err != nil {
		logrus.Warnf("cannot create %s for USB Media sync", pw)
		return 0, err
	}
	defer r.Close()

	sz, err := io.Copy(w, r)
	if err != nil {
		logrus.Warnf("cannot copy from %s to %s for USB Media sync", pr, pw)
		return 0, err
	}

	err =  os.Chtimes(pw, stat.ModTime(), stat.ModTime())
	return sz, err
}


func (exchange *USBExchange) List(since time.Time) ([]string, error) {
	mounts, err := getMounts()
	if err != nil {
		return nil, err
	}
	glob := map[string]string{}
	locs := map[string][]string{}

	for _, mount := range mounts {
		root := filepath.Join(mount, exchange.remote)
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && info.ModTime().After(since) {
				if loc, err := filepath.Rel(root, path); err == nil {
					glob[loc] = root
					locs[mount] = append(locs[mount], loc)
				}
			}
			return nil
		})
	}

	go func (files map[string]string, locs map[string][]string){
		for f, source := range files {
			for mount, fs := range locs {
				if _, found := core.FindStringInSlice(fs, f); !found {
					_, _ = copyFile(source, f, mount)
				}
			}
		}
	}(glob, locs)

	var files []string
	for file := range glob {
		files = append(files, file)
	}

	return files, nil
}

func (exchange *USBExchange) Push(loc string) (int64, error) {
	mounts, err := getMounts()
	if err != nil {
		return 0, err
	}
	var sz int64
	for _, mount := range mounts {
		dest := filepath.Join(mount, exchange.remote)
		sz, err = copyFile(exchange.local, loc, dest)
		if err != nil {
			return 0, err
		}
	}

	return sz, nil
}

func (exchange *USBExchange) Pull(loc string) (int64, error) {
	mounts, err := getMounts()
	if err != nil {
		return 0, err
	}

	parts := strings.Split(loc, "/")
	name := parts[len(parts)-1]
	folder := filepath.Join(parts[0:len(parts)-1]...)

	var sz int64
	for _, mount := range mounts {
		source := filepath.Join(mount, exchange.remote, folder)
		dest := filepath.Join(exchange.local, folder)

		sz, err = copyFile(source, name, dest)
		if err != nil {
			return sz, err
		}
	}

	return sz, nil
}

func (exchange *USBExchange) Delete(string, time.Time) error {
	return nil
}

func (exchange *USBExchange) Name() string {
	return exchange.config.Name
}

func (exchange *USBExchange) String() string {
	return fmt.Sprintf("usb %s", exchange.config.Name)
}
