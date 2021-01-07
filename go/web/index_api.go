package web

import (
	"almost-scrum/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func indexRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/index/suggest/:prefix", getSuggestAPI)
}

func getSuggestAPI(c *gin.Context) {
	var project core.Project
	err := getProject(c, &project)
	if err != nil {
		return
	}


	prefix := c.Param("prefix")
	total, _ := strconv.Atoi(c.DefaultQuery("total", "10"))

	suggestions := core.SuggestKeys(project, prefix, total)
	c.JSON(http.StatusOK, suggestions)
}

