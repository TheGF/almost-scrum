package library

import (
	"almost-scrum/assets"
	"almost-scrum/core"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func embedImage(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	mime, err := mimetype.DetectFile(file)
	if err != nil {
		return "", err
	}

	enc := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mime.String(), enc), nil
}

var reImage = regexp.MustCompile(`<img\s+src="~library/([^#]+)#([^"]+)"`)

func replaceImg(html string, libraryFolder string) string {
	var output bytes.Buffer
	pos := 0
	matches := reImage.FindAllStringSubmatchIndex(html, -1)
	for _, match := range matches {

		img := html[match[0]:match[1]]
		loc := html[match[2]:match[3]]
		alt := html[match[4]:match[5]]
		opts := strings.Split(alt, ",")

		imageSrc, err := embedImage(filepath.Join(libraryFolder, loc))
		if err != nil {
			output.WriteString(html[pos:match[1]])
			pos = match[1]
			continue
		}

		classNames := make([]string, 0)
		size := ""
		for _, opt := range opts {
			parts := strings.Split(opt, "=")
			if len(parts) != 2 {
				continue
			}

			switch parts[0] {
			case "align":
				{
					classNames = append(classNames, fmt.Sprintf("%sImage", parts[1]))
				}
			case "size":
				{
					size = fmt.Sprintf("width=\"%s%%\"", parts[1])
				}
			}
		}

		img = fmt.Sprintf(`<img class="%s" %s src="%s"`,
			strings.Join(classNames, " "), size, imageSrc)

		output.WriteString(html[pos:match[0]])
		output.WriteString(img)
		pos = match[1]
	}

	output.WriteString(html[pos:])
	return output.String()
}

func ExportMarkdownToHTML(file string, libraryFolder string) (string, error) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	body := blackfriday.Run(input, blackfriday.WithExtensions(
		blackfriday.CommonExtensions|
			blackfriday.HardLineBreak|
			blackfriday.NoEmptyLineBeforeBlock,
	))
	return replaceImg(string(body), libraryFolder), nil
}

type BookSettings struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Authors  string   `json:"authors"`
	Styles   []string `json:"styles"`
}

//var whitespaceRe = regexp.MustCompile(`\s+`)
func addStyle(style string, output *bytes.Buffer) {
	name := fmt.Sprintf("assets/styles/%s.css", style)
	css, err := assets.Asset(name)
	if err != nil {
		logrus.Warnf("unknown css style %s: %v", style, err)
		return
	}
	//css = whitespaceRe.ReplaceAll(css, []byte{})
	logrus.Debugf("adding css %s name: %s", name, string(css))
	output.Write(css)
}

var sectionTitleRe = regexp.MustCompile(`\d*\.?\s*(.*)\.(pg)|(md)$`)

func CreateBook(project *core.Project, loc string, settings BookSettings) (string, error) {
	var output bytes.Buffer
	folder := filepath.Join(project.Path, core.ProjectLibraryFolder, loc)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		logrus.Errorf("cannot open library folder %s: %v", folder, err)
		return "", err
	}

	title := settings.Title
	subtitle := settings.Subtitle
	output.WriteString(fmt.Sprintf("<html><title>%s</title><head><style type=\"text/css\">", title))

	addStyle("common", &output)
	for _, style := range settings.Styles {
		addStyle(style, &output)
	}

	output.WriteString(`</style></head><body>`)

	output.WriteString(fmt.Sprintf(`<section class="cover"><h1>%s</h1>`, title))
	if subtitle != "" {
		output.WriteString(fmt.Sprintf(`<h2>%s</h2>`, subtitle))
	}
	output.WriteString(`</section>`)

	libraryFolder := filepath.Join(project.Path, core.ProjectLibraryFolder, "")
	for idx, file := range files {
		name := file.Name()
		match := sectionTitleRe.FindStringSubmatch(name)
		if len(match) == 4 {
			sectionTitle := match[1]
			part, err := ExportMarkdownToHTML(filepath.Join(folder, name), libraryFolder)
			if err == nil {
				output.WriteString(
					fmt.Sprintf(`<section id="section-%d"><div class="sectionTitle">%s</div>`,
						idx, sectionTitle))
				output.WriteString(part)
				output.WriteString(`</section>`)
			}
		}
	}

	output.WriteString("</body></html>")
	return output.String(), nil
}
