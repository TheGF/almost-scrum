package query

import (
	"almost-scrum/core"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestQueryTasks(t *testing.T) {

	folder, _ := ioutil.TempDir(os.TempDir(), "stg")

	p, err := core.InitProject(folder, []string{"scrum", "issue-tracker"})
	assert.Nilf(t, err, "Cannot initialize project: %w", err)

	_, _, err = core.CreateTask(p, "backlog", "Test1", "feature", core.GetSystemUser())
	assert.Nilf(t, err, "Cannot create task: %w", err)

	_, _, err = core.CreateTask(p, "backlog", "Test4", "issue", core.GetSystemUser())
	assert.Nilf(t, err, "Cannot create task: %w", err)

	t1, _, err := core.CreateTask(p, "backlog", "Test2", "feature", core.GetSystemUser())
	assert.Nilf(t, err, "Cannot create task: %w", err)
	t1.Properties["Status"] = "#Done"
	err = core.SetTask(p, "sandbox", "3. Test3", t1)
	assert.Nilf(t, err, "Cannot save task: %w", err)

	tr, _ := QueryTasks(p, Query{})
	assert.Equal(t, 4, len(tr))
	assert.Equal(t, "", tr[0].Task.Description)

	tr, _ = QueryTasks(p, Query{
		Select:          Select{
			Properties:  true,
		},
		WhereTypes: []WhereType{{
			HasPropertiesAll: []string{"Start", "End"},
		}},
	})
	assert.Equal(t, 3, len(tr))


	tr, _ = QueryTasks(p, Query{
		Select:          Select{
			Properties:  true,
		},
		WhereProperties: []WhereProperty{{
			Name:          "Status",
			ValueIsAnyOf:  []string{"#Done"},
		}},
		WhereBoardIs: []string{"sandbox"},
	})
	assert.Equal(t, 1, len(tr))
	assert.Equal(t, "#Done", tr[0].Task.Properties["Status"])
	
	err = core.ShredProject(p)
	assert.Nilf(t, err, "Cannot shred project: %w", err)

}