// Package core provides basic functionality for Almost Scrum
package core

import (
	"bytes"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	paragraphMatch = regexp.MustCompile("#+\\s+(\\S+)[^#]+")
)

func parseFeatures(input []byte, task *Task) {
	parser := blackfriday.New()
	node := parser.Parse(input)

	for node = node.FirstChild; node.Type != blackfriday.List; node = node.Next {
		if node == nil {
			return
		}
	}

	task.Features = map[string]string{}
	for listItem := node.FirstChild; listItem != nil; listItem = listItem.Next {
		text := listItem
		for ; text != nil; text = text.FirstChild {
			if text.Type == blackfriday.Text {
				t := string(text.Literal)
				parts := strings.Split(t, ":")
				if len(parts) != 2 {
					continue
				}
				key := strings.ToLower(strings.TrimSpace(parts[0]))
				val := strings.TrimSpace(parts[1])
				logrus.Debugf("ParseTask - found feature %s: %s", key, val)
				task.Features[key] = val
			}
		}
	}
}

func renderFeatures(task *Task, output *bytes.Buffer) {
	output.WriteString("\n\n### Features\n")

	for key, val := range task.Features {
		output.WriteString("- ")
		output.WriteString(key)
		output.WriteString(": ")
		output.WriteString(val)
		output.WriteString("\n")
	}
}

func ReadTask(path string, task *Task) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return ParseTask(data, task)
}

func WriteTask(path string, task *Task) error {
	data := RenderTask(task)
	return ioutil.WriteFile(path, data, 0644)
}

func RenderTask(task *Task) []byte {
	var output bytes.Buffer

	output.WriteString(task.Description)
	renderFeatures(task, &output)

	return output.Bytes()
}

func ParseTask(input []byte, task *Task) error {
	var description bytes.Buffer
	locS := paragraphMatch.FindAllSubmatchIndex(input, -1)

	if len(locS) > 0 {
		description.Write(input[0:locS[0][0]])
	}

	for _, loc := range locS {
		paragraph := input[loc[0]:loc[1]]
		title := input[loc[2]:loc[3]]

		switch string(title) {
		case "Features":
			parseFeatures(paragraph, task)
		default:
			description.Write(paragraph)
		}
	}

	task.Description = description.String()
	return nil
}
