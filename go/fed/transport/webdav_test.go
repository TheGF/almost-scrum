package transport

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)



var webDAVConfig = WebDAVConfig{
	Name:     "nextcloud",
	URL:      "https://use01.thegood.cloud/remote.php/dav/files/almost_scrum@protonmail.com/",
	Username: "almost_scrum@protonmail.com",
	Password: "ScrumAlmost42@",
	Timeout:  60,
}

func TestConnect(t *testing.T) {

	exchange := GetWebDAVExchanges(webDAVConfig)[0]

	_, err := exchange.Connect("test", "/tmp")
	assert.Nil(t, err)
	locs, err := exchange.List(time.Time{})
	assert.Nil(t, err)

	print(locs)
}

func TestPush(t *testing.T) {

	exchange := GetWebDAVExchanges(webDAVConfig)[0]

	dir, _ := ioutil.TempDir(os.TempDir(), "ash-test")
	_, err := exchange.Connect("test", dir)
	assert.Nil(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "test1"), []byte("test content"), 0755)
	assert.Nil(t, err)

	_, err = exchange.Push("test1")
	assert.Nil(t, err)

	_ = os.RemoveAll(dir)
}

func TestPull(t *testing.T) {
	exchange := GetWebDAVExchanges(webDAVConfig)[0]

	dir, _ := ioutil.TempDir(os.TempDir(), "ash-test")
	_, err := exchange.Connect("test", dir)
	assert.Nil(t, err)

	_, err = exchange.Pull("test1")
	assert.Nil(t, err)

	_, err = os.Stat(filepath.Join(dir, "test1"))
	assert.Nil(t, err)

	_ = os.RemoveAll(dir)
}
