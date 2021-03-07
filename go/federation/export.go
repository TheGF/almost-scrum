package federation

import (
	"almost-scrum/attributes"
	"almost-scrum/core"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/alexmullins/zip"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var exportFolders = []string{
	core.ProjectBoardsFolder,
	core.ProjectArchiveFolder,
	core.ProjectLibraryFolder,
}

const maxExportSizePerFile = 16 * 1024 * 1024

func getFiles(path string, time time.Time) []string {
	var files []string

	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && info.ModTime().After(time) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func addHeader(projectId string, writer *zip.Writer, user string, time time.Time) error {
	config := core.ReadConfig()
	header := Header{
		ProjectID: projectId,
		PeerID:    config.UUID,
		Peer:      "",
		Time:      time,
		User:      user,
	}
	bytes, err := json.Marshal(header)
	if err != nil {
		return nil
	}
	w, err := writer.Create(FedHeaderFile)
	if err != nil {
		logrus.Warnf("cannot create header for federated export: %v", err)
		return err
	}

	if _, err := w.Write(bytes); err != nil {
		return err
	}
	return nil
}

func addFileToPackage(base string, secret string, writer *zip.Writer, file string) (int64, error) {
	stat, err := os.Stat(file)
	if err != nil {
		logrus.Warnf("cannot stat file %s for federated export: %v", file, err)
		return 0, err
	}

	attr, err := attributes.GetExtendedAttr(filepath.Dir(file), stat.Name())
	if err != nil {
		logrus.Warnf("cannot get extended attrs of file %s for federated export: %v", file, err)
		return 0, err
	}

	r, err := os.Open(file)
	if err != nil {
		logrus.Warnf("cannot open file %s for federated export: %v", file, err)
		return 0, err
	}
	defer r.Close()

	name, _ := filepath.Rel(base, file)
	fh := &zip.FileHeader{
		Name:   name,
		Method: zip.Deflate,
		UncompressedSize64: uint64(stat.Size()),
		Comment: hex.EncodeToString(attr.Hash),
	}
	fh.SetModTime(stat.ModTime())
	fh.SetPassword(secret)

	w, err :=  writer.CreateHeader(fh)
	if err != nil {
		logrus.Warnf("cannot encrypt file %s to federated export: %v", file, err)
		return 0, err
	}
	size, err := io.Copy(w, r)
	if err != nil {
		logrus.Warnf("cannot encrypt file %s to federated export: %v", file, err)
		return 0, err
	}

	return size, nil
}

func nextWriter(path string, time time.Time, writer *zip.Writer,
	files []string) (writer_ *zip.Writer, files_ []string, err error) {
	if writer != nil {
		_ = writer.Flush()
		_ = writer.Close()
	}

	zipName := fmt.Sprintf("%x.%x", time.Unix(), len(files))
	path = filepath.Join(path, zipName)
	zipFile, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Cannot create federated export %s: %v", path, err)
		return nil, nil, err
	}
	files_ = append(files, path)
	return zip.NewWriter(zipFile), files_, nil
}

func buildPackage(project *core.Project, dest string, user string, time time.Time, files []string) ([]string, error) {
	dest = filepath.Join(dest, "files")
	if err := os.MkdirAll(dest, 0755); err != nil {
		logrus.Errorf("cannot create output folder %s: %v", dest, err)
		return nil, err
	}

	var zipFiles []string
	total := int64(0)
	writer, zipFiles, err := nextWriter(dest, time, nil, zipFiles)
	if err != nil {
		return nil, err
	}

	addHeader(project.Config.UUID, writer, user, time)
	for _, file := range files {
		size, _ := addFileToPackage(project.Path, project.Config.CipherKey, writer, file)
		total += size
		if size > maxExportSizePerFile {
			writer, zipFiles, err = nextWriter(dest, time, writer, zipFiles)
			if err != nil {
				return nil, err
			}
			total = 0
		}
	}
	_ = writer.Flush()
	writer.Close()

	if len(zipFiles) > 0 {
		last := zipFiles[len(zipFiles)-1]
		name :=  fmt.Sprintf("%s.zip", last)
		os.Rename(last, name)
		zipFiles[len(zipFiles)-1] = name
	}

	return zipFiles, nil
}

func uploadHub(path string, hub Hub) error {
	infos, err := ioutil.ReadDir(filepath.Join(path, hub.ID()))
	if err != nil {
		return err
	}

	for _, info := range infos {
		file := filepath.Join(path, hub.ID(), info.Name())
		if err := hub.Push(file); err == nil {
			if err := os.Remove(file); err != nil {
				logrus.Errorf("cannot remove file %s from output queue in federation export: %v",
					file, err)
			}
		} else {
			logrus.Warnf("cannot export file %s to federated hub %s: %v", file, hub, err)
		}
	}
	return nil
}

func uploadHubs(path string, hubs []Hub) {
	for _, hub := range hubs {
		if err := hub.Connect(); err != nil {
			logrus.Warnf("cannot connect to hub %s: %v", hub, err)
			continue
		}
		_ = uploadHub(path, hub)
		hub.Disconnect()
	}
}

func createLinks(path string, files []string, hubs []Hub) error {
	for _, hub := range hubs {
		p := filepath.Join(path, hub.ID())
		_ = os.Mkdir(p, 0755)
		for _, file := range files {
			name := filepath.Base(file)
			if err := os.Symlink(file, filepath.Join(p, name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func cleanUp(path string, hubs []Hub) {
	p := filepath.Join(path, "files")

	names := map[string]bool{}
	infos, err := ioutil.ReadDir(p)
	if err != nil {
		return
	}
	for _, info := range infos {
		names[info.Name()] = true
	}

	links := map[string]bool{}
	for _, hub := range hubs {
		p := filepath.Join(path, hub.ID())
		infos, err := ioutil.ReadDir(p)
		if err != nil {
			continue
		}
		for _, info := range infos {
			name := info.Name()
			if _, found := names[name]; !found {
				logrus.Errorf("broken link %s in hub %s", name, hub)
				os.Remove(filepath.Join(p, name))
			} else {
				links[name] = true
			}
		}
	}

	for name := range names {
		if _, found := links[name]; !found {
			logrus.Infof("remove output file %s because no hub links to it", name)
			os.Remove(filepath.Join(p, name))
		}
	}
}

func Export(project *core.Project, user string, time time.Time) error {
	if err := CheckTime(); err != nil {
		return err
	}

	var files []string
	hubs, config, err := getHubs(project)
	if err != nil {
		return err
	}

	for _, path := range exportFolders {
		path = filepath.Join(project.Path, path)
		files = append(files, getFiles(path, time)...)
	}

	dest := filepath.Join(project.Path, core.ProjectFedFolder, "out")
	zipFiles, _ := buildPackage(project, dest, user, time, files)
	for _, zipFile := range zipFiles {
		logrus.Infof("Ready zip %s", zipFile)
	}

	if err := createLinks(dest, zipFiles, hubs); err != nil {
		logrus.Errorf("cannot create links to files to be exported: %v", err)
	}

	config.LastExport = time
	WriteConfig(project, config)

	uploadHubs(dest, hubs)
	cleanUp(dest, hubs)

	return nil
}

