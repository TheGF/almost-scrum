package web

import (
	"almost-scrum/chat"
	"almost-scrum/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

func chatRoute(group *gin.RouterGroup) {
	group.POST("/projects/:project/chat", postChatAPI)
}


func postChatAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	file, _ := c.FormFile("file")
	reader, err := file.Open()
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot open upload stream: %v", err)
		return
	}
	defer reader.Close()

	if err = chat.AddMessage(project, getWebUser(c), reader); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot upload message to chat:%v", err)
		return
	}

	c.JSON(http.StatusOK, "")
}

