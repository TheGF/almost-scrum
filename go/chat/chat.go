package chat

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"almost-scrum/library"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Message struct {
	Id    string   `json:"id"`
	User  string   `json:"user"`
	Text  string   `json:"text"`
	Names []string `json:"names"`
	Mimes []string `json:"mimes"`
	Likes []string `json:"likes"`
}

func AddMessage(project *core.Project, message Message, attachments []io.ReadCloser) error {
	folder := filepath.Join(project.Path, core.ProjectChatFolder)
	message.Id = fmt.Sprintf("%x", time.Now().UnixNano()/1000)

	for idx, attachment := range attachments {
		filename := filepath.Join(folder, fmt.Sprintf("%s.%x.bin", message.Id, idx))
		writer, err := os.Create(filename)
		if err != nil {
			logrus.Warnf("Cannot open file %s for writing: %v", filename, err)
			return err
		}
		defer writer.Close()

		_, err = io.Copy(writer, attachment)
		if err != nil {
			logrus.Warnf("cannot write file %s in chat: %v", filename, err)
			return err
		}

		mime, _ := mimetype.DetectFile(filename)
		message.Mimes = append(message.Mimes, mime.String())
	}

	return writeMessage(project, message)
}

func writeMessage(project *core.Project, message Message) error {
	filename := filepath.Join(project.Path, core.ProjectChatFolder, fmt.Sprintf("%s.json", message.Id))
	err := fs.WriteJSON(filename, &message)
	if err != nil {
		logrus.Warnf("cannot write file %s in chat: %v", filename, err)
	} else {
		logrus.Debugf("successfully set message %s in chat", filename)
	}
	return err
}

func readMessage(project *core.Project, messageId string) (Message, error) {
	var message Message
	filename := filepath.Join(project.Path, core.ProjectChatFolder, fmt.Sprintf("%s.json", messageId))
	err := fs.ReadJSON(filename, &message)
	if err != nil {
		logrus.Warnf("cannot read message with id %s in chat: %v", messageId, err)
	} else {
		logrus.Debugf("successfully read message with id %s in chat", messageId)
	}
	return message, err
}

func Like(project *core.Project, messageId string, user string) error {
	message, err := readMessage(project, messageId)
	if err != nil {
		return err
	}

	if idx, found := core.FindStringInSlice(message.Likes, user); found {
		message.Likes[idx] = message.Likes[len(message.Likes)-1]
		message.Likes = message.Likes[:len(message.Likes)-1]
	} else {
		message.Likes = append(message.Likes, user)
	}

	return writeMessage(project, message)
}

func MakeTask(project *core.Project, board string, title string, type_ string, owner string, messageId string) error {
	message, err := readMessage(project, messageId)
	if err != nil {
		return err
	}

	task, id, err := core.CreateTask(project, board, title, type_, owner)
	if err != nil {
		return err
	}
	task.Description = message.Text
	return core.SetTask(project, board, id, task)
}

func MakeDoc(project *core.Project, owner string, messageId string, idx int) error {
	message, err := readMessage(project, messageId)
	if err != nil {
		return err
	}

	if idx >= len(message.Names) {
		return os.ErrInvalid
	}
	reader := GetMessageAttachment(project, message.Id, idx)
	if reader == nil {
		return os.ErrNotExist
	}

	_, err = library.SetFileInLibrary(project, message.Names[idx], reader, owner, false)
	reader.Close()
	return err
}

func DeleteChat(project *core.Project, messageId string) {
	folder := filepath.Join(project.Path, core.ProjectChatFolder)
	os.Remove(filepath.Join(folder, fmt.Sprintf("%s.json", messageId)))
	os.Remove(filepath.Join(folder, fmt.Sprintf("%s.bin", messageId)))
}

func reverseIndexes(l, start, end int) (int, int) {
	start = l - start
	if end < 0 {
		end = 0
	} else {
		end = l - end
	}
	return start, end
}

func GetMessageAttachment(project *core.Project, id string, idx int) io.ReadCloser {
	filename := filepath.Join(project.Path, core.ProjectChatFolder, fmt.Sprintf("%s.%x.bin", id, idx))
	file, err := os.Open(filename)
	if err == nil {
		return file
	} else {
		logrus.Errorf("cannot get attachment %s: %v", filename, err)
		return nil
	}
}

func GetMessageAttachmentFilepath(project *core.Project, id string, idx int) (filename string, contentType string) {
	filename = filepath.Join(project.Path, core.ProjectChatFolder, fmt.Sprintf("%s.%x.bin", id, idx))
	m, err := mimetype.DetectFile(filename)
	if err == nil {
		return filename, m.String()
	} else {
		return "", ""
	}
}

func ListMessages(project *core.Project, start int, end int) ([]Message, error) {
	folder := filepath.Join(project.Path, core.ProjectChatFolder)

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			names = append(names, file.Name())
		}
	}

	var messages []Message
	start, end = reverseIndexes(len(names), start, end)
	for i := start - 1; i >= end && i >= 0; i-- {
		var message Message
		err = fs.ReadJSON(filepath.Join(folder, names[i]), &message)
		if err == nil {
			messages = append(messages, message)
		}
	}

	return messages, nil
}
