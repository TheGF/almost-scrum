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
	group.POST("/projects/:project/git/push", postGitPushAPI)
	group.POST("/projects/:project/git/pull", postGitPullAPI)
	group.PUT("/projects/:project/git/credentials", putGitCredentialsAPI)
}


func getGitStatusAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}
	git := core.GetGitClient(project)

	status, err := git.GetStatus(project)
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
	git := core.GetGitClient(project)

	var commitInfo core.CommitInfo
	if err := c.BindJSON(&commitInfo); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hash, err := git.Commit(project, commitInfo)
	if err != nil {
		logrus.Warnf("Cannot commit content project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot commit content: %v", err)
		return
	}
	c.JSON(http.StatusOK, hash)
}

func postGitPullAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}
	git := core.GetGitClient(project)

	commit, err := git.Pull(project, getWebUser(c))
	if err != nil {
		logrus.Warnf("Cannot pull content project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot commit content: %v", err)
		return
	}
	c.JSON(http.StatusOK, commit)
}

func postGitPushAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}
	git := core.GetGitClient(project)

	err := git.Push(project, getWebUser(c))
	if err != nil {
		logrus.Warnf("Cannot push content project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot commit content: %v", err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func putGitCredentialsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var gitCredentials core.GitCredentials
	if err := c.BindJSON(&gitCredentials); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := core.SetGitCredentials(project, getWebUser(c), gitCredentials)
	if err != nil {
		logrus.Warnf("Cannot save git credentials in project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot save git credentials: %v", err)
		return
	}
	c.JSON(http.StatusOK, "")
}
