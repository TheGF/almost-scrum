package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestLinkTag creates a story and link it to a tags
func TestLinkTag(t *testing.T) {
	p, _ := InitProject("../test-data/tag-project")

	s, err := GetStore(p, "backlog")
	assert.Nilf(t, err, "Cannot open backlog: %w", err)

	name := "123.Test.story"
	SetStory(s, name, &Story{})

	err = LinkTag(s, name, "sample")
	assert.Nilf(t, err, "Cannot assign tag: %w", err)

	list, err := ResolveTag(p, "sample")
	println("RESOLVE")
	for _, t := range list {
		fmt.Printf("%v", t)
	}
}
