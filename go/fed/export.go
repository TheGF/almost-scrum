package fed

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexmullins/zip"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)


const maxExportSizePerFile = 16 * 1024 * 1024

func shouldBeExported(file string, includePrivate bool) (bool, error) {
	attr, err := fs.GetExtendedAttr(file)
	if err != nil {
		return false, err
	}
	if !includePrivate && !attr.Public {
		return false, err
	}

	hash, err := fs.GetHash(file)
	if err != nil {
		logrus.Warnf("cannot get hash of file %s for federated export: %v", file, err)
		return false, err
	}

	logrus.Debugf("Compare hash %v = %v", hash, attr.ImportHash)

	if bytes.Compare(hash, attr.ExportHash) == 0 {
		logrus.Debugf("current hash and last export hash for file %s are the same. Skip export", file)
		return false, nil
	}

	return true, nil
}

func getFiles(path string, includePrivate bool, time time.Time) []string {
	var files []string

	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && info.ModTime().After(time) && info.Name() != fs.AttrsFileName {
			if ok, err := shouldBeExported(path, includePrivate); err != nil {
				logrus.Warnf("issues with file %s: %v", path, err)
			} else if ok {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}

const v1 = "v1"

func addHeader(origin string, writer *zip.Writer, user string, time time.Time) error {
	header := Header{
		Version:  v1,
		Host:     origin,
		Time:     time,
		User:     user,
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
		logrus.Warnf("cannot Throughput file %s for federated export: %v", file, err)
		return 0, err
	}

	attr, err := fs.GetExtendedAttr(file)
	if err != nil {
		logrus.Warnf("cannot get extended attrs of file %s for federated export: %v", file, err)
		return 0, err
	}

	hash, err := fs.GetHash(file)
	if err != nil {
		logrus.Warnf("cannot get hash of file %s for federated export: %v", file, err)
		return 0, err
	}

	r, err := os.Open(file)
	if err != nil {
		logrus.Warnf("cannot open file %s for federated export: %v", file, err)
		return 0, err
	}
	defer r.Close()

	name, _ := filepath.Rel(base, file)
	comment := fmt.Sprintf("%s,%v,%v", attr.Owner, attr.ImportHash, hash )

	fh := &zip.FileHeader{
		Name:               name,
		Method:             zip.Deflate,
		UncompressedSize64: uint64(stat.Size()),
		Comment:            comment,
	}
	fh.SetModTime(stat.ModTime())
	fh.SetPassword(secret)

	w, err := writer.CreateHeader(fh)
	if err != nil {
		logrus.Warnf("cannot encrypt file %s to federated export: %v", file, err)
		return 0, err
	}
	size, err := io.Copy(w, r)
	if err != nil {
		logrus.Warnf("cannot encrypt file %s to federated export: %v", file, err)
		return 0, err
	}

	attr.ExportHash = hash
	fs.SetExtendedAttr(file, attr)

	logrus.Infof("export file %s: owner %s, importHash %v, exportHash %v", file, attr.Owner, attr.ImportHash,
		hash)
	return size, nil
}

func nextWriter(path string, prefix string, time time.Time, writer *zip.Writer,
	files []string) (writer_ *zip.Writer, files_ []string, err error) {
	if writer != nil {
		_ = writer.Flush()
		_ = writer.Close()
	}

	zipName := fmt.Sprintf("%s%x.%x.zip", prefix, time.Unix(), len(files))
	path = filepath.Join(path, zipName)
	zipFile, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Cannot create federated export %s: %v", path, err)
		return nil, nil, err
	}
	files_ = append(files, path)
	return zip.NewWriter(zipFile), files_, nil
}

func buildPackage(project *core.Project, origin string, prefix string, user string, files []string) ([]string, error) {
	dest := filepath.Join(project.Path, core.ProjectFedFilesFolder, origin)

	if err := os.MkdirAll(dest, 0755); err != nil {
		logrus.Errorf("cannot create output dest %s: %v", origin, err)
		return nil, err
	}

	tm := time.Now()
	var zipFiles []string
	total := int64(0)
	cnt := 0
	writer, zipFiles, err := nextWriter(dest, prefix, tm, nil, zipFiles)
	if err != nil {
		return nil, err
	}

	if err = addHeader(origin, writer, user, tm); err != nil {
		logrus.Warnf("cannot add header in zip export: %v", err)
		return nil, err
	}
	for _, file := range files {
		size, _ := addFileToPackage(project.Path, project.Config.CipherKey, writer, file)
		total += size
		cnt += 1
		if size > maxExportSizePerFile {
			writer, zipFiles, err = nextWriter(origin, prefix, tm, writer, zipFiles)
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

func export(project *core.Project, item syncItem, user string,
	host string, since time.Time) ([]string, error) {
	var files []string
	for _, path := range item.folders {
		path = filepath.Join(project.Path, path)
		files = append(files, getFiles(path, item.includePrivate, since)...)
	}
	if len(files) == 0 {
		return nil, nil
	}
	logrus.Infof("file changed since %s: %v", since, files)

	origin := fmt.Sprintf("%s.%s", project.Config.Public.Name, host)
	zipFiles, err := buildPackage(project, origin, item.prefix, user, files)
	if err != nil {
		return nil, err
	}
	logrus.Infof("created export files %v", zipFiles)

	var locs []string
	for _, file := range files {
		loc, _ := filepath.Rel(project.Path, file)
		locs = append(locs, loc)
	}

	return locs, nil
}

func Export(project *core.Project, user string, since time.Time) ([]string, error) {
	var files []string

	signal, err := Connect(project)
	if err != nil {
		return nil, err
	}

	config, err := ReadConfig(project, false)
	if err != nil {
		return nil, err
	}
	if since.IsZero() {
		since = config.LastExport
	}

	c := core.ReadConfig()
	for _, item := range syncItems {
		files_, err := export(project, item, user, c.Host, since)
		if err != nil {
			return nil, err
		}
		files = append(files, files_...)
	}
	if len(files) == 0 {
		logrus.Infof("no changes to export since %s", since)
		return nil, nil
	}

	logrus.Infof("file changed since %s: %v", since, files)


	for _, file := range files {
		signal.exports[file] = time.Time{}
	}

	config.LastExport = time.Now()
	_ = WriteConfig(project, config)

	return files, nil
}

func ExportLast(project *core.Project, user string) ([]string, error) {
	return Export(project, user, time.Time{})
}
