package chat

import (
	"almost-scrum/core"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestChat(t *testing.T) {

	folder, _ := ioutil.TempDir(os.TempDir(), "stg")

	p, err := core.InitProject(folder, []string{"scrum", "issue-tracker"})
	assert.Nilf(t, err, "Cannot initialize project: %w", err)

	AddMessage(p, Message{
		User: "Me",
		Text: "Hello",
	}, nil)

	AddMessage(p, Message{
		User: "Me",
		Text: "World",
	}, nil)

	msgs, _ := ListMessages(p, 0, 10)
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, "Hello", msgs[1].Text)
}