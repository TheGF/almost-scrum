package web

import (
	"almost-scrum/core"
	"almost-scrum/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func queryRoute(group *gin.RouterGroup) {
	group.POST("/projects/:project/query/tasks", postQueryTasksAPI)
}


func postQueryTasksAPI(c *gin.Context) {
	var project *core.Project
	var q query.Query
	if project = getProject(c); project == nil {
		return
	}

	if err := c.BindJSON(&q); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ts, err := query.QueryTasks(project, q)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ts)
}

