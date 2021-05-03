package web

import (
	"almost-scrum/core"
	"github.com/code-to-go/fed"
	"github.com/code-to-go/fed/transport"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//fedRoute add fed related api routes
func fedRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/fed/transport", getTransportAPI)
	group.POST("/projects/:project/fed/transport", setTransportAPI)
	group.GET("/projects/:project/fed/state", getStateAPI)
	group.POST("/projects/:project/fed/merge", postMergeAPI)
	group.POST("/projects/:project/fed/sync", postAutoSyncAPI)
	group.POST("/projects/:project/fed/invite", postCreateInviteAPI)
	group.POST("/projects/:project/fed/join", postJoinAPI)
}

func getTransportAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	config := project.Fed.GetTransport()
	c.JSON(http.StatusOK, config)
}

func setTransportAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var config transport.Config
	if err := c.BindJSON(&config); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	project.Fed.SetTransport(&config)

	c.JSON(http.StatusOK, "")
}

func getStateAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	s := project.Fed.GetState()
	logrus.Debugf("Fed status in project: %v", s)
	c.JSON(http.StatusOK, s)
}

func postMergeAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var actions []fed.Action
	stateId,_ := strconv.Atoi(c.Query("id"))

	if err := c.BindJSON(&actions); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ok := project.Fed.Merge(int64(stateId), actions)
	c.JSON(http.StatusOK, ok)
}

func postAutoSyncAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	p, _ := strconv.Atoi(c.Query("period"))
	period := time.Duration(p) * time.Minute
	project.Fed.AutoSync(period, getWebUser(c))

	c.JSON(http.StatusOK, "")
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

	invite, err := project.Fed.Invite(request.Key, request.Exchanges, request.RemoveCredentials)
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

	if err := fed.Join(request.Key, request.Token, filepath.Join(project.Path, "fed")); err == os.ErrInvalid {
		c.String(http.StatusBadRequest, "Invalid token or key")
	} else if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "")
}
