// Package core provides basic functionality for Almost Scrum
package core

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	idMatch = regexp.MustCompile(`^([\pN]+)\.(.*)`)
)

// Part to complete the story and its status
type Part struct {
	Description string `json:"description" yaml:"description"`
	Done        bool   `json:"done" yaml:"done"`
}

// TimeEntry is time used by a user on an activity on a specific day
type TimeEntry struct {
	User  string    `json:"user" yaml:"user"`
	Date  time.Time `json:"date" yaml:"date"`
	Hours int       `json:"hours" yaml:"hours"`
}

// TaskInfo is the result of List operation
type TaskInfo struct {
	ID      uint16    `json:"id"`
	Board   string    `json:"board"`
	Name    string    `json:"name"`
	ModTime time.Time `json:"modTime"`
}

// Task contains attributes that define a story
type Task struct {
	Description string            `json:"description"`
	Properties  map[string]string `json:"properties"`
	Parts       []Part            `json:"parts"`
	Files       []string          `json:"files"`
	ConflictId  string            `json:"conflictId"`
}

// ListBoardTasks list the tasks in the board
func ListTasks(project *Project, board string, filter string) ([]TaskInfo, error) {
	var infos = make([]TaskInfo, 0)

	if board == "" {
		boards, err := ListBoards(project)
		if IsErr(err, "Cannot list boards in %s", project.Path) {
			return infos, err
		}
		for _, board := range boards {
			err := listTasksForBoard(project, board, filter, &infos)
			if IsErr(err, "cannot list tasks in %s/%s", project.Path, board) {
				return infos, err
			}
		}
	} else {
		err := listTasksForBoard(project, board, filter, &infos)
		if IsErr(err, "cannot list tasks in %s/%s", project.Path, board) {
			return infos, err
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime.After(infos[j].ModTime)
	})

	return infos, nil
}

func listTasksForBoard(project *Project, board string, filter string, infos *[]TaskInfo) error {
	p := filepath.Join(project.Path, ProjectBoardsFolder, board)
	fileInfos, err := ioutil.ReadDir(p)
	if IsErr(err, "cannot read board %s", board) {
		return err
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		ext := path.Ext(name)
		if ext != TaskFileExt {
			continue
		}
		name = strings.TrimSuffix(name, ext)
		if filter != "" && !strings.Contains(name, filter) {
			continue
		}
		id, _ := ExtractTaskId(name)
		if id == 0 {
			continue
		}

		*infos = append(*infos, TaskInfo{
			ID:      id,
			Board:   board,
			Name:    name,
			ModTime: fileInfo.ModTime(),
		})
	}
	return nil
}

// GetTask a story in the Board
func GetTask(project *Project, board string, name string) (task Task, err error) {
	p := filepath.Join(project.Path, ProjectBoardsFolder, board, name+TaskFileExt)
	if err = ReadTask(p, &task); IsErr(err, "Cannot read task %s/%s", board, name) {
		return task, err
	}
	return task, nil
}

// GetTaskPath returns the absolute path of a story
func GetTaskPath(project *Project, board string, name string) string {
	p := filepath.Join(project.Path, ProjectBoardsFolder, board, name+TaskFileExt)
	p, _ = filepath.Abs(p)
	return p
}

func ExtractTaskId(name string) (uint16, string) {
	match := idMatch.FindStringSubmatch(name)
	if len(match) < 3 {
		return 0, ""
	}
	id, _ := strconv.Atoi(match[1])
	return uint16(id), match[2]
}

func CreateTask(project *Project, board string, title string, type_ string, owner string) (*Task, string, error) {
	task := Task{
		Description: "",
		Properties:  map[string]string{},
		Parts:       []Part{},
		Files:       []string{},
	}

	for _, model := range project.Models {
		if model.Name == type_ {
			if err := ParseTask(model.Template, &task); err != nil {
				return nil, "", err
			}
			for _, property := range model.Properties {
				if _, found := task.Properties[property.Name]; !found {
					task.Properties[property.Name] = property.Default
					logrus.Debugf("Set property %s = %s", property.Name, property.Default)
				}
			}
			task.Properties["Type"] = type_
			task.Properties["Owner"] = "@" + owner

			name := NewTaskName(project, title)
			if err := SetTask(project, board, name, &task); err != nil {
				return nil, "", err
			}

			project.TasksCount += 1
			return &task, name, nil
		}
	}
	return nil, "", ErrInvalidType
}

//SetTask a story in the Board
func SetTask(project *Project, board string, id string, task *Task) error {
	p := filepath.Join(project.Path, ProjectBoardsFolder, board, id+TaskFileExt)
	if err := WriteTask(p, task); IsErr(err, "cannot save task %s/%s", board, id) {
		return err
	}

	return nil
}

// GetTask a story in the Board
func DeleteTask(project *Project, board string, name string) (task Task, err error) {
	p := filepath.Join(project.Path, ProjectBoardsFolder, board, name+TaskFileExt)
	if err = os.Remove(p); IsErr(err, "Cannot delete task %s/%s", board, name) {
		return task, err
	}

	return task, nil
}

// TouchTask set the modified time to current time. It applies to stories and folders
func TouchTask(project *Project, board string, name string) error {
	currentTime := time.Now().Local()
	p := filepath.Join(project.Path, ProjectBoardsFolder, board, name+TaskFileExt)
	if err := os.Chtimes(p, currentTime, currentTime); IsErr(err, "cannot touch %s/%s", board, name) {
		return err
	}

	return nil
}

func MoveTask(project *Project, oldBoard string, oldName string, board string, name string) error {
	source := filepath.Join(project.Path, ProjectBoardsFolder, oldBoard, oldName+TaskFileExt)
	target := filepath.Join(project.Path, ProjectBoardsFolder, board, name+TaskFileExt)

	if err := os.Rename(source, target); err != nil {
		return err
	}

	currentTime := time.Now().Local()
	_ = os.Chtimes(target, currentTime, currentTime)

	return nil
}
