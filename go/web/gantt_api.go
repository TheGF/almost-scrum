package web

import (
	"almost-scrum/core"
	"almost-scrum/gantt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ganttRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/gantt", getGanttTasksAPI)
}


func getGanttTasksAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	if tasks, err := gantt.GetTasks(project); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, tasks)
	}
}

