package core

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var story = Task{
	Description: "Test of a story",
	Features: map[string]string{
		"points": "12",
	},
	Tasks: []Step{{
		Description: "Test of a task",
		Done:        true,
	},
		{
			Description: "Test of a task2",
			Done:        false,
		},
	},
	TimeEntries: []TimeEntry{},
	Attachments: []string{},
}

func TestListTasks(t *testing.T) {
	project, _ := OpenProject(".")
	infos, err := ListTasks(project, "", "")
	assert.NotNilf(t, err, "cannot list tasks in project %s: %v", project.Path, err)
	assert.GreaterOrEqual(t, len(infos), 1)
}

func TestSetStory(t *testing.T) {
	project, _ := OpenProject(".")
	id := "1.Hello.story"
	err := SetTask(project, "backlog", id, &story)
	assert.NotNilf(t, err, "cannot write backlog/%s in project %s: %v", id, err )
}

func TestGetStory(t *testing.T) {
	project, _ := OpenProject(".")
	id := "1.Hello.story"
	s, _ := GetTask(project, "backlog", id)
	assert.Equal(t, story, s, "Mismatch in story read")
}

func BenchmarkGet(b *testing.B) {
	project, _ := OpenProject(".")
	id := "1.Hello.story"
	for i := 0; i < b.N; i++ {
		_, _ = GetTask(project, "board", id)
	}
}


func TestMarkdown(t *testing.T) {
	data, err := ioutil.ReadFile("test.md")
	assert.Nilf(t, err, "Cannot open markdown file: %w", err)

	task := Task{}
	err = ParseTask(data, &task)
	assert.Nilf(t, err, "Cannot parse markdown data: %w", err)

	out := string(RenderTask(&task))
	for k,v := range task.Features {
		assert.Contains(t, out, k)
		assert.Contains(t, out, v)
	}
}

