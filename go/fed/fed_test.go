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

func getProject(t *testing.T) *core.Project {
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

	data := core.GenerateRandomString(1024 * 1024)
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

	data := core.GenerateRandomString(1024 * 1024)
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

	c, _ := ReadConfig(project, false)
	token, _ := CreateInvite(project, "banana", c)

	logrus.Print(token)

	png, err := qrcode.Encode(token, qrcode.Medium, 256)
	assert.Nil(t, err)
	ioutil.WriteFile("/tmp/qr.png", png, 0755)
}

func TestClaimInvite(t *testing.T) {
	token := `1cc7bde244d4dfc788a10aaac818aa28f2c596b6f85def11c9438410eded0bd3825f8e2b24a5ad17610de8a91ecb9f4da\
3125095c65d9fce44398ca6c432ca3a3ea5a40a20e27814dd32a6debae75741eb171a9797dd08fe8980c1cc45967c7936c93c942\
4399a3d10ec7981d8ce044a22bed88cc2a56736d81017e8993c59308d43d2f84efde8835cdbe0c860e4133e6d2329482bcc9f26f\
899ba4a86d844dfb51ac11ec84c762c53324c12a82fe39e8859536c820eb2f6ee2cdcfbd44b28890932825575930ae21055b1384\
1d8492426c1b4901b5592dd4b2c791521266cd6e977120b6d93c7ec1b0ed6ad9a9f30e062792a19b8814ac6a1db7bd539a8ff843\
b1c25970cf404518bfa0952f19fabb6c538b33771cf73ae67c8751478b402007c1864734164284f3433424e89185039f9f90ccdd\
767d9c56cc4bc85fc6cee5ff7d391c5dd1aa0f058fdd4494879d5b1561c142797d2ca282c9524d489c107ed7e708392490a4be60\
48c0dcd172e20af2b09b6bcb2f0c8b63785ddff0da5fa810c3afd8a235ef0247fa7c38442ae83de62ac0f8a3c5c3e91ba08c8e3e\
a7a65657364c6a9c666e5b1d686f85bda573da5bf2ae92a3124d655646e1231f5d87a476dd6115618152ea40bd2ec9d3ff053ec3\
f2dd180942ecdfce31e9dadfcecc4b560d095acf21d880c7ad518bbce36e0e9fdfa5494b2892d1dbffdbf169c53c89faf65be2db\
f17b049501d7c5e87be716e38599e1cb4e463f3b72d810fc417caa127225ff3a3ec1625013351d547cdde0c7318ef0bf66a7d0bb\
334b54bde8a64175c154dd162d88e4054d7879e132abb327fca9f25d2469d2eb3a7b92309ce9e0b1d3738d0b44d76221d5e9f3cf\
0234a3af303184af3d9e7c0a86d49d0cc26a1dcfdaa51de6d83e1d62c5515f3a2dfb097d22430b360b800fac9a101bf715bf233f\
79ad47477cab18dcf08139b6a83bfd8baebae2da0f32195e12ff8dd67ea83b61c0e0270f2db9afe1b161df2ae50e3905152d02ab\
8398a21d7d8a0034c0c8e68431966ab4c02ca356a9921305a395fb474c8e0ac8bd48d78234352b59bd4102d728fe1ea605bb89ae\
f32f99a03d6f81de7cb56f0cd87f447199a695a5dd02a555ca92eb1357d4a520f781b7b19149390c456e72b9b93c16a8083fa8a2\
3f08805f21c63c27122567b4a8ea24d6925e80b827cf152c46342581071e13ad87c6d0c7d7147ab8def76f023b26918620bfe4cb\
1a283b879fbaf9b5a759615a0911be779851d36eaef83afc251a609ee88cb99`
	key := "banana"

	project := getProject(t)
	Join(project, key, token)
}

func TestMergeConfig(t *testing.T) {
	project := getProject(t)

	mergeFedConfig(project, &Config{
		Secret:        "",
		ReconnectTime: 0,
		PollTime:      0,
		Span:          0,
		LastExport:    time.Time{},
		S3: []transport.S3Config{
			{Name: "Amazon S3", Endpoint: "xx", Bucket: "xx"},
			{Name: "B3", Endpoint: "xx", Bucket: "xx"},
		},
		WebDAV: nil,
		Ftp:    nil,
	})
}
