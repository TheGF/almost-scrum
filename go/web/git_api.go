package web

import (
	"almost-scrum/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func gitRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/git/status", getGitStatusAPI)
	group.POST("/projects/:project/git/commit", postGitCommitAPI)
}

func getGitStatusAPI(c *gin.Context) {
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

func postGitCommitAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var commitInfo core.CommitInfo
	if err := c.BindJSON(&commitInfo); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hash, err := core.GitCommit(project, commitInfo)
	if err != nil {
		logrus.Warnf("Cannot commit content project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot commit content: %v", err)
		return
	}
	c.JSON(http.StatusOK, hash)
}
