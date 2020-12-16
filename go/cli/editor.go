package cli

import (
	"almost-scrum/core"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rivo/tview"
)

func openEditor(s core.Store, path string) {
	var editor string = config.Editor
	if editor == "" {
		openEmbeddedEditor(s, path)
	} else {
		p := core.GetStoryAbsPath(s, path)
		openExternalEditor(editor, p)
		if story, err := core.GetStory(s, path); err == nil {
			core.LinkTagsFromStory(s, path, &story)
		}
	}
}

func openEmbeddedEditor(s core.Store, path string) {
	story, err := core.GetStory(s, path)
	if err != nil {
		color.Red("Something went wrong with reading story %s: %v", path, err)
		os.Exit(1)
	}

	app := tview.NewApplication()
	form := tview.NewForm().
		AddInputField("Description", story.Description, 0, nil, nil).
		AddDropDown("Points", []string{"1", "2", "3", "5", "9"}, 0, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddCheckbox("Age 18+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", nil).
		AddButton("Quit", func() {
			app.Stop()
		})
	title := fmt.Sprintf("Story: %s", path)
	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
