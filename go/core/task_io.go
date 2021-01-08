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
	paragraphMatch  = regexp.MustCompile(`#+\s+(\S+)[^#]+`)
	lineMatch       = regexp.MustCompile(`#+\s+(\w+)|([^\n]+)`)
	propertyMatch   = regexp.MustCompile(`\s*([^:]*):\s+(.*)`)
	partMatch       = regexp.MustCompile(`\[([x ])]\s+(.+)`)
	attachmentMatch = regexp.MustCompile(`\[[^]]*]\([^)]+\)`)
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
	var node *blackfriday.Node
	parser := blackfriday.New()
	node = parser.Parse(input)
	if node == nil || node.FirstChild == nil {
		return
	}

	for node = node.FirstChild; node != nil && node.Type != blackfriday.List; node = node.Next {
	}
	if node == nil {
		return
	}

	for listItem := node.FirstChild; listItem != nil; listItem = listItem.Next {
		text := listItem
		for ; text != nil; text = text.FirstChild {
			if text.Type == blackfriday.Text {
				switch title {
				case "Properties":
					parseProperties(text, task)
				case "Progress":
					parseParts(text, task)
				case "Attachments":
					parseAttachments(text, task)
				}
			}
		}
	}
}

func renderProperties(task *Task, output *bytes.Buffer) {
	if task.Properties == nil {
		return
	}
	output.WriteString("### Properties\n")

	for key, val := range task.Properties {
		output.WriteString("- ")
		output.WriteString(key)
		output.WriteString(": ")
		output.WriteString(val)
		output.WriteString("\n")
	}
}

func renderParts(task *Task, output *bytes.Buffer) {
	if task.Parts == nil {
		return
	}
	output.WriteString("### Progress\n")

	for _, part := range task.Parts {
		output.WriteString("- ")
		if part.Done {
			output.WriteString("[x] ")
		} else {
			output.WriteString("[ ] ")
		}
		output.WriteString(part.Description)
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
	if !strings.HasSuffix(task.Description, "\n") {
		output.WriteString("\n")
	}
	renderProperties(task, &output)
	renderParts(task, &output)

	return output.Bytes()
}

type paragraph struct {
	header string
	title   string
	body    string
}

func splitInParagraph(input []byte) []paragraph {
	paragraphs := []paragraph{}

	lines := lineMatch.FindAllStringSubmatch(string(input), -1)
	for _, line := range lines {
		all := line[0]
		title := line[1]
		part := line[2]

		if title != "" {
			paragraphs = append(paragraphs, paragraph{all+"\n", title, ""})
			continue
		}

		if len(paragraphs) == 0 {
			paragraphs = append(paragraphs, paragraph{"", "", ""})
		}
		paragraphs[len(paragraphs)-1].body += part + "\n"
	}

	return paragraphs
}

func ParseTask(input []byte, task *Task) error {
	var description bytes.Buffer

	task.Properties = map[string]string{}
	task.Parts = []Part{}
	task.Attachments = []string{}

	paragraphs := splitInParagraph(input)
	for _, paragraph := range paragraphs {
		switch paragraph.title {
		case "Properties", "Progress", "Attachments":
			parseList([]byte(paragraph.body), paragraph.title, task)
		default:
			description.WriteString(paragraph.header)
			description.WriteString(paragraph.body)
		}
	}

	task.Description = description.String()
	return nil
}
