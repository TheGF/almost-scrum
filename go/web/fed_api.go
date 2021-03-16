package web

import (
	"almost-scrum/core"
	"almost-scrum/fed"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//fedRoute add fed related api routes
func fedRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/fed/hubs", getLogsAPI)
	group.GET("/projects/:project/fed/logs", getLogsAPI)
	group.POST("/projects/:project/fed/merge", postMergeAPI)
	group.POST("/projects/:project/fed/export", postExportAPI)
}

func getLogsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	err := fed.Pull(project, time.Now().AddDate(0, 0, -7))
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot get import: %v", err)
		return
	}
	logs, err := fed.Match(project)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot get fed logs: %v", err)
		return
	}
	if logs == nil {
		logs = []*fed.Diff{}
	}

	logrus.Debugf("Fed logs in project: %v", logs)
	c.JSON(http.StatusOK, logs)
}


func postMergeAPI(c *gin.Context) {

	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var logs []*fed.Diff
	if err := c.BindJSON(&logs); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	files, err := fed.Diff(project, logs)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot merge fed files: %v", err)
		return
	}

	c.JSON(http.StatusOK, files)
}


func postExportAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	since, found := c.GetQuery("since")

	var files []string
	var err error
	if found {
		tm, err := time.Parse(time.RFC3339, since)
		if err != nil {
			c.String(http.StatusBadRequest, "wrong format for parameter since: %s", since)
			return
		}
		files, err = fed.Export(project, getWebUser(c), tm)
	} else {
		files, err = fed.ExportLast(project, getWebUser(c))
	}
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot export fed files: %v", err)
		return
	}

	c.JSON(http.StatusOK, files)

}

