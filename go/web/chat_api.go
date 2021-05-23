package web

import (
	"almost-scrum/chat"
	"almost-scrum/core"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func chatRoute(group *gin.RouterGroup) {
	group.POST("/projects/:project/chat", postChatAPI)
	group.GET("/projects/:project/chat", getChatAPI)
	group.GET("/projects/:project/chat/:id/:idx", getChatAttachmentAPI)
	group.POST("/projects/:project/chat/:id/:idx", postChatItemAPI)
	group.DELETE("/projects/:project/chat/:id", deleteChatItemAPI)
}

func getChatAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	s := c.DefaultQuery("start", "0")
	e := c.DefaultQuery("end", "-1")

	start, _ := strconv.Atoi(s)
	end, _ := strconv.Atoi(e)

	messages, err := chat.ListMessages(project, start, end)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot upload message to chat:%v", err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func getChatAttachmentAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	id := c.Param("id")
	idx, _ := strconv.Atoi(c.Param("idx"))
	filepath, contentType := chat.GetMessageAttachmentFilepath(project, id, idx)

	if filepath == "" {
		c.String(http.StatusNotFound, "No content")
		return
	} else {
		c.Writer.Header().Set("Content-Type", contentType)
		c.File(filepath)
	}
}

func postChatAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	message := chat.Message{
		User: getWebUser(c),
		Text: c.PostForm("text"),
	}

	var readers []io.ReadCloser
	for i := 0; ; i++ {
		file, _ := c.FormFile(fmt.Sprintf("file-%d", i))
		if file == nil {
			break
		}
		reader, err := file.Open()
		if err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot open upload stream: %v", err)
			return
		}
		readers = append(readers, reader)
		message.Names = append(message.Names, file.Filename)
	}

	err := chat.AddMessage(project, message, readers)
	for _, reader := range readers {
		reader.Close()
	}

	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot upload message to chat:%v", err)
		return
	}

	c.JSON(http.StatusOK, "")
}

func deleteChatItemAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	id := c.Param("id")
	chat.DeleteChat(project, id)

	c.String(http.StatusOK, "")
}

func postChatItemAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	id := c.Param("id")
	idx, _ := strconv.Atoi(c.Param("id"))
	action := c.Query("action")
	board := c.Query("board")
	title := c.Query("title")
	type_ := c.Query("type")

	var err error
	switch action {
	case "like":
		err = chat.Like(project, id, getWebUser(c))
	case "make_post":
		err = chat.MakeTask(project, board, title, type_, getWebUser(c), id)
	case "make_doc":
		err = chat.MakeDoc(project, getWebUser(c), id, idx)
	}

	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot apply action %s to message %s to chat:%v",
			action, id, err)
		return
	}

	c.JSON(http.StatusOK, "")
}
