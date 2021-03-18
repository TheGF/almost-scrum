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
	group.GET("/projects/:project/fed/config", getConfigAPI)
	group.GET("/projects/:project/fed/status", getStatusAPI)
	group.GET("/projects/:project/fed/diffs", getDiffsAPI)
	group.POST("/projects/:project/fed/import", postImportAPI)
	group.POST("/projects/:project/fed/export", postExportAPI)
	group.POST("/projects/:project/fed/sync", postSyncAPI)
	group.POST("/projects/:project/fed/share", postShareAPI)
}

func getConfigAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	config, err := fed.ReadConfig(project, false)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot read federation config: %v", err)
		return
	}
	logrus.Debugf("Fed config in project: %v", config)
	c.JSON(http.StatusOK, config)
}

func getStatusAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	s := fed.GetStatus(project)
	logrus.Debugf("Fed status in project: %v", s)
	c.JSON(http.StatusOK, s)
}

func getDiffsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	_, sync := c.GetQuery("sync")

	if sync {
		_, err := fed.Sync(project, time.Time{})
		if err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot sync with federation: %v", err)
			return
		}
	}
	diffs, err := fed.GetDiffs(project)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot get fed diffs: %v", err)
		return
	}
	if diffs == nil {
		diffs = []*fed.Diff{}
	}

	logrus.Debugf("Fed diffs in project: %v", diffs)
	c.JSON(http.StatusOK, diffs)
}

func postImportAPI(c *gin.Context) {

	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var diffs []*fed.Diff
	if err := c.BindJSON(&diffs); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	files, err := fed.Import(project, diffs)
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

func postSyncAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	_, err := fed.Sync(project, time.Time{})
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot export fed files: %v", err)
		return
	}

	c.JSON(http.StatusOK, fed.GetStatus(project))

}

type ShareRequest struct {
	Exchanges         []string `json:"exchanges"`
	RemoveCredentials bool     `json:"removeCredentials"`
}
type ShareResponse struct {
	Key   string `json:"key"`
	Token string `json:"token"`
}

func postShareAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var shareRequest ShareRequest
	if err := c.BindJSON(&shareRequest); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	key, token, err := fed.ShareWith(project, shareRequest.Exchanges, shareRequest.RemoveCredentials)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, ShareResponse{
		Key:   key,
		Token: token,
	})
}
