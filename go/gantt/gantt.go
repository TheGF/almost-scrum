package gantt

import (
	"almost-scrum/core"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	ganttCache = cache.New(5*time.Minute, 10*time.Minute)
)

type Task struct {
	Board string    `json:"board"`
	Name  string    `json:"name"`
	Task  core.Task `json:"task"`
}

type State struct {
	Tasks   []*Task
	Tracked map[string]time.Time
}

func hasValidType(project *core.Project, task core.Task) bool {
	var hasStart, hasEnd bool
	type_ := task.Properties["Type"]
	for _, model := range project.Models {
		if model.Name == type_ {
			for _, property := range model.Properties {
				switch property.Name {
				case "Start":
					hasStart = true
				case "End":
					hasEnd = true
				}
			}
			break
		}
	}
	return hasStart && hasEnd
}

func DeriveTask(project *core.Project, info core.TaskInfo) (*Task, error) {
	task, err := core.GetTask(project, info.Board, info.Name)
	if err != nil {
		return nil, err
	}
	if !hasValidType(project, task) {
		return nil, nil
	}

	return &Task{
		Board: info.Board,
		Name:     info.Name,
		Task: task,
	}, nil
}

func GetTasks(project *core.Project) ([]*Task, error) {
	var state *State
	s, found := ganttCache.Get(project.Config.UUID)
	if found {
		state = s.(*State)
	} else {
		state = &State{}
		ganttCache.Set(project.Config.UUID, state, cache.DefaultExpiration)
	}

	infos, err := core.ListTasks(project, "", "")
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		modTime, found := state.Tracked[info.Name]
		if found && modTime == info.ModTime {
			continue
		}

		task, err := DeriveTask(project, info)
		if err != nil {
			return nil, err
		}
		if task != nil {
			for idx, t := range state.Tasks {
				if t.Name == task.Name {
					state.Tasks[idx] = task
					goto done
				}
			}
			state.Tasks = append(state.Tasks, task)
		}
		done:
	}

	return state.Tasks, nil
}
