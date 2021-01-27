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
	group.PUT("/projects/:project/git/settings", putGitSettingsAPI)
	group.GET("/projects/:project/git/settings", getGitSettingsAPI)
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

func putGitSettingsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var gitSettings core.GitSettings
	if err := c.BindJSON(&gitSettings); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := core.SetGitSettings(project, getWebUser(c), gitSettings)
	if err != nil {
		logrus.Warnf("Cannot save git settings in project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot save git settings: %v", err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func getGitSettingsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	gitSettings, err := core.GetGitSettings(project, getWebUser(c))
	if err != nil {
		logrus.Warnf("Cannot get git settings in project %s: %v", project.Path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot get git settings: %v", err)
		return
	}
	c.JSON(http.StatusOK, gitSettings)
}
