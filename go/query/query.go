package query

import (
	"almost-scrum/core"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

var queryCache = cache.New(5*time.Minute, 10*time.Minute)

func filterByBoard(infos []core.TaskInfo, whereBoardIs []string) []core.TaskInfo {
	var r []core.TaskInfo

	for _, info := range infos {
		if _, found := core.FindStringInSlice(whereBoardIs, info.Board); found {
			r = append(r, info)
		}
	}
	return r
}


func getTaskRef(project *core.Project, info core.TaskInfo) (*TaskRef, error) {
	key := path.Join(project.Config.UUID, info.Board, info.Name)

	t, found := queryCache.Get(key)
	if found && t.(*TaskRef).ModTime == info.ModTime {
		return t.(*TaskRef), nil
	}

	task, err := core.GetTask(project, info.Board, info.Name)
	if err != nil {
		return nil, err
	}

	taskRef := TaskRef{
		Board:   info.Board,
		Name:    info.Name,
		ModTime: info.ModTime,
		Task:    task,
	}

	queryCache.Set(key, &taskRef, cache.DefaultExpiration)
	return &taskRef, nil
}

func hasValidType(taskRef *TaskRef, types []string) bool {
	return core.HasStringInSlice(types, taskRef.Task.Properties["Type"])
}

func hasRequiredProperties(taskRef *TaskRef, whereProperties []WhereProperty) bool {
	for _, whereProperty := range whereProperties {
		for name, value := range taskRef.Task.Properties {
			if whereProperty.Name != name {
				continue
			}

			if len(whereProperty.ValueIsAnyOf) > 0 && !core.HasStringInSlice(whereProperty.ValueIsAnyOf, value) {
				return false
			}

			if len(whereProperty.ValueIsNoneOf) > 0 && core.HasStringInSlice(whereProperty.ValueIsNoneOf, value) {
				return false
			}
		}
	}
	return true
}

func selectContent(refs []*TaskRef, select_ Select) []TaskRef {
	var rs []TaskRef

	for _, ref := range refs {
		r := TaskRef{
			Board:   ref.Board,
			Name:    ref.Name,
			ModTime: ref.ModTime,
			Task:    core.Task{},
		}

		if select_.Description {
			r.Task.Description = ref.Task.Description
		}
		if select_.Properties {
			r.Task.Properties = ref.Task.Properties
		}
		if select_.Parts {
			r.Task.Parts = ref.Task.Parts
		}
		if select_.Files {
			r.Task.Files = ref.Task.Files
		}
		rs = append(rs, r)
	}
	return rs
}

func QueryTasks(project *core.Project, params Query) ([]TaskRef, error) {
	infos, err := core.ListTasks(project, "", "")
	if err != nil {
		return nil, err
	}

	validTypes := QueryTypes(project, params.WhereTypes)
	if params.WhereBoardIs != nil {
		infos = filterByBoard(infos, params.WhereBoardIs)
	}

	var refs []*TaskRef
	for _, info := range infos {
		taskRef, err := getTaskRef(project, info)
		if err != nil {
			logrus.Errorf("cannot get task ref for %s/%s/%s", project.Config.Public.Name, taskRef.Board, taskRef.Name)
			continue
		}

		if len(params.WhereTypes) > 0 && !hasValidType(taskRef, validTypes) {
			continue
		}

		if len(params.WhereProperties) > 0 && !hasRequiredProperties(taskRef, params.WhereProperties) {
			continue
		}
		refs = append(refs, taskRef)
	}

	return selectContent(refs, params.Select), nil
}

func QueryTypes(project *core.Project, params []WhereType) []string {
	var r []string

	for _, where := range params {
		for _, model := range project.Models {
			if len(where.Is) > 0 {
				_, found := core.FindStringInSlice(where.Is, model.Name)
				if !found {
					continue
				}
			}

			if len(where.HasPropertiesAll) > 0 {
				cnt := 0
				for _, propertyDef := range model.Properties {
					if _, found := core.FindStringInSlice(where.HasPropertiesAll, propertyDef.Name); found {
						cnt++
					}
				}
				if cnt != len(where.HasPropertiesAll) {
					continue
				}
			}

			r = append(r, model.Name)
		}
	}
	return r
}
