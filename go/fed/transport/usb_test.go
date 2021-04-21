package transport

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var usbConfig = USBConfig{
	Name:      "USB Media",
}

func TestUSBConnect(t *testing.T) {

	exchange := GetUSBExchanges(usbConfig)[0]

	_, err := exchange.Connect("test", "/tmp")
	assert.Nil(t, err)
	locs, err := exchange.List(time.Time{})
	assert.Nil(t, err)

	print(locs)
}

func TestUSBPush(t *testing.T) {

	exchange := GetUSBExchanges(usbConfig)[0]

	dir, _ := ioutil.TempDir(os.TempDir(), "ash-test")
	_, err := exchange.Connect("test", dir)
	assert.Nil(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "test1"), []byte("test content"), 0755)
	assert.Nil(t, err)

	_, err = exchange.Push("test1")
	assert.Nil(t, err)

	_ = os.RemoveAll(dir)
}

func TestUSBPull(t *testing.T) {
	exchange := GetUSBExchanges(usbConfig)[0]

	dir, _ := ioutil.TempDir(os.TempDir(), "ash-test")
	_, err := exchange.Connect("test", dir)
	assert.Nil(t, err)

	_, err = exchange.Pull("test1")
	assert.Nil(t, err)

	_, err = os.Stat(filepath.Join(dir, "test1"))
	assert.Nil(t, err)

	_ = os.RemoveAll(dir)
}
