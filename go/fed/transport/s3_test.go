package transport

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var s3Config = S3Config{
	Name:      "nextcloud",
	Endpoint:  "s3.eu-central-1.amazonaws.com",
	AccessKey: "AKIA275FPHM3BLLN2NX5",
	Secret:    "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",
	Bucket:    "almost-scrum",
	UseSSL:    false,
	Location:  "eu-central-1",
}

//var webDAVConfig = WebDAVConfig{
//	Name:     "nextcloud",
//	URL:      "https://shared03.opsone-cloud.ch/remote.php/dav/files/cestino83@protonmail.com",
//	Username: "cestino83@protonmail.com",
//	Password: "83cestino",
//	Timeout:  60,
//}

func TestS3Connect(t *testing.T) {

	exchange := GetS3Exchanges(s3Config)[0]

	_, err := exchange.Connect("test", "/tmp")
	assert.Nil(t, err)
	locs, err := exchange.List(time.Time{})
	assert.Nil(t, err)

	print(locs)
}

func TestS3Push(t *testing.T) {

	exchange := GetS3Exchanges(s3Config)[0]

	dir, _ := ioutil.TempDir(os.TempDir(), "ash-test")
	_, err := exchange.Connect("test", dir)
	assert.Nil(t, err)

	err = ioutil.WriteFile(filepath.Join(dir, "test1"), []byte("test content"), 0755)
	assert.Nil(t, err)

	err = exchange.Push("test1")
	assert.Nil(t, err)

	os.RemoveAll(dir)
}

func TestS3Pull(t *testing.T) {
	exchange := GetS3Exchanges(s3Config)[0]

	dir, _ := ioutil.TempDir(os.TempDir(), "ash-test")
	_, err := exchange.Connect("test", dir)
	assert.Nil(t, err)

	err = exchange.Pull("test1")
	assert.Nil(t, err)

	_, err = os.Stat(filepath.Join(dir, "test1"))
	assert.Nil(t, err)

	os.RemoveAll(dir)
}
