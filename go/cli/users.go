package cli

import "almost-scrum/core"

func processUsers(projectPath string, args []string) {
	project, err := core.FindProject(projectPath, "")
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		listStores(project)
		return
	}
	switch args[0] {
	case "cd":
		setCurrentStore(project)
	}
}
