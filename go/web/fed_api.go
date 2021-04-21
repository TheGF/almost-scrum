package web

import (
	"almost-scrum/core"
	"almost-scrum/fed"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

//fedRoute add fed related api routes
func fedRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/fed/config", getConfigAPI)
	group.GET("/projects/:project/fed/status", getStatusAPI)
//	group.GET("/projects/:project/fed/diffs", getDiffsAPI)
	group.POST("/projects/:project/fed/config", setConfigAPI)
	group.POST("/projects/:project/fed/import", postImportAPI)
	group.POST("/projects/:project/fed/export", postExportAPI)
	group.POST("/projects/:project/fed/pull", postPullAPI)
	group.POST("/projects/:project/fed/push", postPushAPI)
	group.POST("/projects/:project/fed/share", postCreateInviteAPI)
	group.POST("/projects/:project/fed/join", postJoinAPI)
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
	c.JSON(http.StatusOK, config)
}

func setConfigAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var config fed.Config
	if err := c.BindJSON(&config); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := fed.WriteConfig(project, &config)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot read federation config: %v", err)
		return
	}
	fed.Disconnect(project)
	c.JSON(http.StatusOK, "")
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

func postImportAPI(c *gin.Context) {

	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var updates []fed.Update
	if err := c.BindJSON(&updates); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	locs, err := fed.Import(project, updates)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot merge fed locs: %v", err)
		return
	}

	c.JSON(http.StatusOK, locs)
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

func postPushAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	stats, err := fed.Push(project)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot export fed files: %v", err)
		return
	}

	c.JSON(http.StatusOK, stats)

}

func postPullAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	stats, err := fed.Pull(project, time.Time{})
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot export fed files: %v", err)
		return
	}

	c.JSON(http.StatusOK, stats)

}

type ShareRequest struct {
	Key               string   `json:"key"`
	Exchanges         []string `json:"exchanges"`
	RemoveCredentials bool     `json:"removeCredentials"`
}

func postCreateInviteAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var request ShareRequest
	if err := c.BindJSON(&request); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	invite, err := fed.CreateInviteForExchanges(project, request.Key, request.Exchanges, request.RemoveCredentials)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, invite)
}

type JoinRequest struct {
	Token string `json:"token"`
	Key   string `json:"key"`
}

func postJoinAPI(c *gin.Context) {
	var request JoinRequest
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	if err := c.BindJSON(&request); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := fed.Join(project, request.Key, request.Token); err == os.ErrInvalid {
		c.String(http.StatusBadRequest, "Invalid token or key")
	} else if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "")
}
