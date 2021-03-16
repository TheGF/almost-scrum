package fed

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/alexmullins/zip"
	"github.com/sirupsen/logrus"
	"io"
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
		if info != nil && !info.IsDir() && info.ModTime().After(time) && info.Name() != fs.AttrsFileName {
			if attr, _ := fs.GetExtendedAttr(path); attr.Public {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}

func addHeader(projectId string, writer *zip.Writer, user string, time time.Time) error {
	config := core.ReadConfig()
	header := Header{
		ProjectID: projectId,
		Host:      config.Host,
		Hostname:  config.Hostname,
		Time:      time,
		User:      user,
	}
	bytes_, err := json.Marshal(header)
	if err != nil {
		return nil
	}
	w, err := writer.Create(HeaderFile)
	if err != nil {
		logrus.Warnf("cannot create header for federated export: %v", err)
		return err
	}

	if _, err := w.Write(bytes_); err != nil {
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

	attr, err := fs.GetExtendedAttr(file)
	if err != nil {
		logrus.Warnf("cannot get extended attrs of file %s for federated export: %v", file, err)
		return 0,err
	}

	hash, err := fs.GetHash(file)
	if err != nil {
		logrus.Warnf("cannot get hash of file %s for federated export: %v", file, err)
		return 0,err
	}
	if bytes.Compare(hash, attr.Origin) == 0 {
		return 0, nil
	}

	r, err := os.Open(file)
	if err != nil {
		logrus.Warnf("cannot open file %s for federated export: %v", file, err)
		return 0,err
	}
	defer r.Close()

	name, _ := filepath.Rel(base, file)
	comment := fmt.Sprintf("%s,%s", attr.Owner, hex.EncodeToString(hash))

	fh := &zip.FileHeader{
		Name:   name,
		Method: zip.Deflate,
		UncompressedSize64: uint64(stat.Size()),
		Comment: comment,
	}
	fh.SetModTime(stat.ModTime())
	fh.SetPassword(secret)

	w, err :=  writer.CreateHeader(fh)
	if err != nil {
		logrus.Warnf("cannot encrypt file %s to federated export: %v", file, err)
		return 0,err
	}
	size, err := io.Copy(w, r)
	if err != nil {
		logrus.Warnf("cannot encrypt file %s to federated export: %v", file, err)
		return 0,err
	}

	logrus.Infof("file %s added to federated export", file)
	return size, nil
}

func nextWriter(path string, time time.Time, writer *zip.Writer,
	files []string) (writer_ *zip.Writer, files_ []string, err error) {
	if writer != nil {
		_ = writer.Flush()
		_ = writer.Close()
	}

	zipName := fmt.Sprintf("exp%x.%x.zip", time.Unix(), len(files))
	path = filepath.Join(path, zipName)
	zipFile, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Cannot create federated export %s: %v", path, err)
		return nil, nil, err
	}
	files_ = append(files, path)
	return zip.NewWriter(zipFile), files_, nil
}

func buildPackage(project *core.Project, dest string, user string, files []string) ([]string, error) {
	if err := os.MkdirAll(dest, 0755); err != nil {
		logrus.Errorf("cannot create output folder %s: %v", dest, err)
		return nil, err
	}

	tm := time.Now()
	var zipFiles []string
	total := int64(0)
	cnt := 0
	writer, zipFiles, err := nextWriter(dest, tm, nil, zipFiles)
	if err != nil {
		return nil, err
	}

	if err = addHeader(project.Config.UUID, writer, user, tm); err != nil {
		logrus.Warnf("cannot add header in zip export: %v", err)
		return nil, err
	}
	for _, file := range files {
		size, _ := addFileToPackage(project.Path, project.Config.CipherKey, writer, file)
		total += size
		cnt += 1
		if size > maxExportSizePerFile {
			writer, zipFiles, err = nextWriter(dest, tm, writer, zipFiles)
			if err != nil {
				return nil, err
			}
			total = 0
		}
	}
	_ = writer.Flush()
	_ = writer.Close()

	return zipFiles, nil
}



func Export(project *core.Project, user string, since time.Time) ([]string, error) {
	var files []string

	config, err := ReadConfig(project)
	if err != nil {
		return nil, err
	}

	if since.IsZero() {
		since = config.LastExport
	}

	for _, path := range exportFolders {
		path = filepath.Join(project.Path, path)
		files = append(files, getFiles(path, since)...)
	}
	if len(files) == 0 {
		logrus.Infof("no changes to export since %s", since)
		return files, nil
	}

	logrus.Infof("file changed since %s: %v", since, files)

	c := core.ReadConfig()
	dest := filepath.Join(project.Path, core.ProjectFedFilesFolder, c.Host)
	zipFiles, err := buildPackage(project, dest, user, files)
	if err != nil {
		return nil, err
	}
	logrus.Infof("ready to commit files %#v", zipFiles)

	config.LastExport = time.Now()
	_ = WriteConfig(project, config)

	return files, nil
}

func ExportLast(project *core.Project, user string) ([]string, error) {
	return Export(project, user, time.Time{})
}
