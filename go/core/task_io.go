// Package core provides basic functionality for Almost Scrum
package core

import (
	"bytes"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
)

var (
	paragraphMatch  = regexp.MustCompile(`#+\s+(\S+)[^#]+`)
	propertyMatch   = regexp.MustCompile(`\s*([^:]*):\s+(.*)`)
	partMatch       = regexp.MustCompile(`\[([x ])]\s+(.+)`)
	attachmentMatch = regexp.MustCompile(`\[[^]]*]\([^\)]+\)`)
)

func parseProperties(node *blackfriday.Node, task *Task) {
	t := string(node.Literal)
	match := propertyMatch.FindStringSubmatch(t)
	if len(match) != 3 {
		return
	}
	key := match[1]
	val := match[2]
	task.Properties[key] = val
	logrus.Debugf("ParseTask - found feature %s: %s", key, val)
}

func parseParts(node *blackfriday.Node, task *Task) {
	t := string(node.Literal)
	match := partMatch.FindStringSubmatch(t)
	if len(match) != 3 {
		return
	}
	done := match[1] == "x"
	description := match[2]
	task.Parts = append(task.Parts, Part{
		Description: description,
		Done:        done,
	})
	logrus.Debugf("ParseTask - found part %s [done = %t]", description, done)
}

func parseAttachments(node *blackfriday.Node, task *Task) {
	if node.Next == nil || node.Next.LinkData.Destination == nil {
		return
	}

	link := string(node.Next.LinkData.Destination)
	task.Attachments = append(task.Attachments, link)
	logrus.Debugf("ParseTask - found attachment %s", link)
}


func parseList(input []byte, title string, task *Task) {
	parser := blackfriday.New()
	node := parser.Parse(input)

	for node = node.FirstChild; node.Type != blackfriday.List; node = node.Next {
		if node == nil {
			return
		}
	}

	task.Properties = map[string]string{}
	for listItem := node.FirstChild; listItem != nil; listItem = listItem.Next {
		text := listItem
		for ; text != nil; text = text.FirstChild {
			if text.Type == blackfriday.Text {
				switch title {
				case "Properties":
					parseProperties(text, task)
				case "Parts":
					parseParts(text, task)
				case "Attachments":
					parseAttachments(text, task)
				}
			}
		}
	}
}

func renderFeatures(task *Task, output *bytes.Buffer) {
	output.WriteString("\n\n### Properties\n")

	for key, val := range task.Properties {
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
		title := string(input[loc[2]:loc[3]])

		switch title {
		case "Properties", "Parts", "Attachments":
			parseList(paragraph, title, task)
		default:
			description.Write(paragraph)
		}
	}

	task.Description = description.String()
	return nil
}