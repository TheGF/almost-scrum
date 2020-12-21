package core

import (
	"errors"
)

// GitFolder is the default folder where git info is stored
const GitFolder = ".git"

// ProjectFolder is the default folder name when the Ash is used inside a Git repository
const ProjectFolder = ".ash"

// ProjectConfigFile is a configuration file for the project
const ProjectConfigFile = "ash.yaml"

// ProjectUsersFolder the folder containing users
const ProjectUsersFolder = "users"

// ProjectUsersFolder the folder containing users
const ProjectLibraryFolder = "library"

const ProjectTagsFolder = "tags"

// ProjectFolders is the required folders in the project
var ProjectFolders = []string{"stores/backlog", "stores/sandbox", ProjectLibraryFolder,
	ProjectTagsFolder, ProjectUsersFolder}

// ErrNoFound occurs when an item is not found
var ErrNoFound = errors.New("No such item")

// ErrTooMany occurs when too many items are found
var ErrTooMany = errors.New("Too many items")

// ErrExists occurs when an item exists even though it should not exist
var ErrExists = errors.New("Already exists")
