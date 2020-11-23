package core

import (
	"testing"
)

func TestFindProject(t *testing.T) {
	FindProject("../test-data", "")

}

func TestInitProject(t *testing.T) {
	InitProject("../test-data")

}
