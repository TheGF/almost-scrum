package web

import (
	"almost-scrum/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func gitRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/git/status", getGitAPI)
}

func getGitAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	status, err := core.GetGitStatus(project)
	if err != nil {
		logrus.Warnf("Cannot get Git status in project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot get Git status: %v", err)
		return
	}
	c.JSON(http.StatusOK, status)
}

