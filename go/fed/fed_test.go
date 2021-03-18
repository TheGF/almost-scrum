package fed

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getProject(t *testing.T) *core.Project{
	path := "/tmp/ash-test"
	_ = os.RemoveAll(path)
	_ = os.MkdirAll(path, 0755)
	_ = core.UnzipProjectTemplates("/tmp/ash-test", []string{"file:/../test-data/test-template.zip"})
	project, err := core.OpenProject(path)
	assert.Nilf(t, err, "Cannot open project: %w", err)

	return project
}

func TestExport(t *testing.T) {
	project := getProject(t)

	since := time.Time{}

	data := core.GenerateRandomString(1024*1024)
	ioutil.WriteFile(filepath.Join(project.Path, core.ProjectLibraryFolder, "Huge.txt"),
		[]byte(data), 0755)

	_, err := Export(project, "marco", since)
	assert.Nilf(t, err, "Cannot export: %v", err)

	Push(project)
	Disconnect(project)
}

func TestImport(t *testing.T) {
	project := getProject(t)

	since := time.Now().AddDate(-1, 0, 0)
	n, err := Pull(project, since)
	assert.Nil(t, err)

	print("pull from ", n)

	diffs, _ := GetDiffs(project)
	_, _ = Import(project, diffs)

	Disconnect(project)
}

func TestSync(t *testing.T) {
	project := getProject(t)

	data := core.GenerateRandomString(1024*1024)
	p := filepath.Join(project.Path, core.ProjectLibraryFolder, "Huge.txt")
	ioutil.WriteFile(p, []byte(data), 0755)

	fs.SetExtendedAttr(p, &fs.ExtendedAttr{
		Owner:   "me",
		Origin:  nil,
		Public:  true,
		Deleted: time.Time{},
	})

	_, err := Export(project, "marco", time.Time{})
	assert.Nilf(t, err, "Cannot export: %v", err)


	failed, err := Sync(project, time.Time{})
	assert.Nil(t, err)
	assert.Equal(t, failed, 0)

	logrus.Printf("status: %#v", GetStatus(project))

}

func TestSharing(t *testing.T) {
	project := getProject(t)

	c,_ := ReadConfig(project, false)
	key, token,_ := Share(project, c)

	logrus.Print(key, token)

	png, err  := qrcode.Encode(key, qrcode.Medium, 256)
	assert.Nil(t, err)
	ioutil.WriteFile("/tmp/qr.png", png, 0755)
}