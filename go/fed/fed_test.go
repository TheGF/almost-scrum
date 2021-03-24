package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
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

func TestInvite(t *testing.T) {
	project := getProject(t)

	c,_ := ReadConfig(project, false)
	invite,_ := CreateInvite(project, c)

	logrus.Print(invite.Key, invite.Token)

	png, err  := qrcode.Encode(invite.Key, qrcode.Medium, 256)
	assert.Nil(t, err)
	ioutil.WriteFile("/tmp/qr.png", png, 0755)
}


func TestClaimInvite(t *testing.T) {
	token := `1d1980595b6d4d19b41a5224b748665a0bf7fc7db511dcc2577bd1cf01ba66aa39c295ad6cb223bfd61be9c6c0de8d92d167
              9acda08a4bdc18a3a0ae4189d0cfe283ba89255d7d53a6ffe5848ea59e4d8a38917989202757ce3dc52d22d0f38b93020443
              4b2ecce698dcd6fcc88cb868fd4977e013139fe29b8621de71dcdcc6e6b102188ced8e8ea3fa073498941175183723562651
              eecbe3278be18c00dc2aae59cef20c488d447c25d2868112a83fe1d10197c02876c307bc7666d705c117ca7818ca1bdc8d05
              d93904e5a28146567dab1d7e4569ccac94594633b725400b19c409a4a1842629101470840600950d38d3a7f6953fb510eb33
              06397eeedf3b2a6945b4a143410f9b81962978a4ad0a8e7259500b9e25855b3a0c3ccc60f10b`
	key := "NfHK5a84jjJkwzDk"

	invite := Invite{
		Key: key,
		Token: token,
	}

	ClaimInvite(invite, "/tmp")
}

func TestMergeConfig(t *testing.T) {
	project := getProject(t)

	mergeFedConfig(project, &Config{
		Secret:        "",
		ReconnectTime: 0,
		PollTime:      0,
		Span:          0,
		LastExport:    time.Time{},
		S3:            [] transport.S3Config{
			{Name: "Amazon S3", Endpoint: "xx", Bucket: "xx"},
			{Name: "B3", Endpoint: "xx", Bucket: "xx"},
		},
		WebDAV:        nil,
		Ftp:           nil,
	})
}