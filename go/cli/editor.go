package cli

import (
	"almost-scrum/core"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
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

func getPoints() []string {
	points := make([]string, 0, 21)
	for i := range points {
		points[i] = strconv.Itoa(i)
	}
	return points
}

func openEmbeddedEditor(s core.Store, path string) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	numSelections := 0
	go func() {
		for _, word := range strings.Split("bla bla bla", " ") {
			if word == "the" {
				word = "[#ff0000]the[white]"
			}
			if word == "to" {
				word = fmt.Sprintf(`["%d"]to[""]`, numSelections)
				numSelections++
			}
			fmt.Fprintf(textView, "%s ", word)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	textView.SetDoneFunc(func(key tcell.Key) {
		currentSelection := textView.GetHighlights()
		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				textView.Highlight()
			} else {
				textView.Highlight("0").ScrollToHighlight()
			}
		} else if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyTab {
				index = (index + 1) % numSelections
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + numSelections) % numSelections
			} else {
				return
			}
			textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
		}
	})
	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	// story, err := core.GetStory(s, path)
	// if err != nil {
	// 	color.Red("Something went wrong with reading story %s: %v", path, err)
	// 	os.Exit(1)
	// }

	// points := getPoints()

	// app := tview.NewApplication()
	// form := tview.NewForm().
	// 	AddInputField("Description", story.Description, 0, nil, func(v string) {
	// 		story.Description = v
	// 	}).
	// 	AddDropDown("Points", points, 0, func(_ string, v int) {
	// 		story.Points = v
	// 	}).
	// 	AddInputField("Last name", "", 20, nil, nil).
	// 	AddCheckbox("Age 18+", false, nil).
	// 	AddPasswordField("Password", "", 10, '*', nil).
	// 	AddButton("Save", nil).
	// 	AddButton("Quit", func() {
	// 		app.Stop()
	// 	})

	// list := tview.NewList().
	// 	AddItem("List item 1", "Some explanatory text", 'a', nil).
	// 	AddItem("List item 2", "Some explanatory text", 'b', nil).
	// 	AddItem("List item 3", "Some explanatory text", 'c', nil).
	// 	AddItem("List item 4", "Some explanatory text", 'd', nil).
	// 	AddItem("Quit", "Press to exit", 'q', func() {
	// 		app.Stop()
	// 	})

	// title := fmt.Sprintf("Story: %s", path)
	// form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	// if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
	// 	panic(err)
	// }

}
