package federation

import (
	"almost-scrum/core"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func getFeeds(dest string, hubs []Hub, time time.Time) (map[string][]Hub, error) {
	previouslyDownloaded := map[string]bool{}
	files, err := ioutil.ReadDir(dest)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		previouslyDownloaded[file.Name()] = true
	}

	feeds := map[string][]Hub{}
	for _, hub := range hubs {
		if names, err := hub.List(time); err == nil {
			for _, name := range names {
				if _, found := previouslyDownloaded[name]; !found {
					feeds[name] = append(feeds[name], hub)
				}
			}
		} else {
			logrus.Warnf("cannot get file list from hub %s: %v", hub, err)
		}
	}

	for _, hs := range feeds {
		rand.Shuffle(len(hs), func(i, j int) {
			hs[i], hs[j] = hs[j], hs[i]
		})
	}

	return feeds, nil
}

func downloadFiles(project *core.Project, since time.Time) error {
	hubs, _, err := connectToHubs(project)
	if err != nil {
		return err
	}
	if len(hubs) == 0 {
		return os.ErrClosed
	}

	dest := filepath.Join(project.Path, core.ProjectFedFolder, "in")
	_ = os.Mkdir(dest, 0755)
	feeds, err := getFeeds(dest, hubs, since)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		logrus.Infof("nothing to import")
		return nil
	}

	for name, hubs := range feeds {
		for _, hub := range hubs {
			if err := hub.Pull(name, dest); err == nil {
				break
			} else {
				logrus.Warnf("cannot pull file %s from hub %s: %v", name, hub, err)
			}
		}
	}
	return nil
}


func Import(project *core.Project, since time.Time) error {
	if err := CheckTime(); err != nil {
		return err
	}

	downloadFiles(project, since)
	MergeFiles(project)
	return nil
}
