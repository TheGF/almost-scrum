package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// TagLink is a link to a story
type TagLink struct {
	Name  string `json:"n"`
	Board string `json:"b"`
	Path  string `json:"p"`
}

// TagLinks is the list of links for a tag
type TagLinks []TagLink

// ListTags lists all the tags in the project p
func ListTags(p Project) ([]string, error) {
	path := filepath.Join(p.Path, ProjectTagsFolder)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Warnf("Error reading tags main folder")
		return nil, err
	}

	tags := make([]string, len(fileInfos))
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		ext := filepath.Ext(fileInfo.Name())
		if ext == ".json" {
			tags = append(tags, strings.TrimSuffix(name, ext))
		}
	}
	return tags, nil
}

// ResolveTag lists all stories linked to the given tag in the project p
func ResolveTag(p Project, tag string) (TagLinks, error) {
	var links TagLinks
	name := fmt.Sprintf("%s.json", tag)
	tagfile := filepath.Join(p.Path, ProjectTagsFolder, name)
	err := ReadJSON(tagfile, &links)
	return links, err
}

func mergeLink(links TagLinks, link TagLink) (TagLinks, bool) {
	for _, l := range links {
		if l.Name == link.Name {
			if l.Board == link.Board && l.Path == link.Path {
				return links, false
			}
			l.Board = link.Board
			l.Path = link.Path
			return links, true
		}
	}
	return append(links, link), true
}

// LinkTag assigns a tag to the given story
func LinkTag(s Store, path string, tag string) error {
	var links TagLinks

	link := TagLink{
		Name:  filepath.Base(path),
		Board: filepath.Base(s.Path),
		Path:  filepath.Dir(path),
	}
	name := fmt.Sprintf("%s.json", tag)
	tagfile := filepath.Join(s.Project.Path, ProjectTagsFolder, name)
	err := ReadJSON(tagfile, &links)
	if err != nil {
		links = TagLinks{link}
		return WriteJSON(tagfile, links)
	} else {
		links, changed := mergeLink(links, link)
		if changed {
			return WriteJSON(tagfile, links)
		}
		return nil
	}
}

// LinkTagsFromStory finds tags stored in a story and add them to the tags index
func LinkTagsFromStory(s Store, path string, story *Story) {
	userTag := fmt.Sprintf("@%s", story.Owner)
	go LinkTag(s, path, userTag)

	for _, tag := range story.Tags {
		go LinkTag(s, path, tag)
	}
}
