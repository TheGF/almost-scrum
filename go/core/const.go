package core

import "errors"

// GitFolder is the default folder where git info is stored
const AshVersion = "0.5.0"

// GitFolder is the default folder where git info is stored
const GitFolder = ".git"

// ProjectFolder is the default folder name when the Scrum is used inside a Git repository
const ProjectFolder = ".scrum-to-go"

// ProjectConfigFile is a configuration file for the project
const ProjectConfigFile = "scrum-to-go.yaml"

// ProjectUsersFolder the folder containing users
const ProjectUsersFolder = "users"

// ProjectUsersFolder the folder containing users
const ProjectLibraryFolder = "library"

// ProjectFedFolder the folder containing fed config and temporary files
const ProjectFedFolder = "fed"
const ProjectFedFilesFolder = "fed/files"

// ProjectUsersFolder the folder containing users
const ProjectArchiveFolder = "archive"

const ProjectChatFolder = "chat"


// ProjectUsersFolder the folder containing users
const ProjectModelsFolder = "models"

// ProjectUsersFolder the folder containing users
const ProjectLibraryInlineImagesFolder = "library/.inline-images"

// ProjectBoardsFolder the folder containing boards
const ProjectBoardsFolder = "boards"

const TaskFileExt = ".md"

const IndexFile = "no-git-index.json"

const TypeProperty = "Type"

var (
	// ProjectFolders is the required folders in the project
	ProjectFolders = []string{ProjectBoardsFolder,
		ProjectLibraryFolder, ProjectUsersFolder,
		ProjectLibraryInlineImagesFolder, ProjectModelsFolder,
		ProjectFedFolder, ProjectFedFilesFolder, ProjectChatFolder}

	ProjectTemplatesPath = "assets/templates/"

	// ErrNoFound occurs when an item is not found
	ErrNoFound = errors.New("no such item")

	// ErrNoFound occurs when an item is not found
	ErrInvalidType = errors.New("invalid type")

	// ErrTooMany occurs when too many items are found
	ErrTooMany = errors.New("too many items")

	// ErrExists occurs when an item exists even though it should not exist
	ErrExists = errors.New("already exists")

	ErrMergeConflict = errors.New("merge conflict")
)
