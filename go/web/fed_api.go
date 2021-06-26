package web

import (
	"almost-scrum/core"
	"bytes"
	"encoding/json"
	"github.com/code-to-go/fed/transport"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

//fedRoute add fed related api routes
func fedRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/fed/transport", getTransportListAPI)
	group.GET("/projects/:project/fed/transport/*id", getTransportAPI)
	group.PUT("/projects/:project/fed/transport/*id", putTransportAPI)
	group.GET("/projects/:project/fed/state/:after", getStateAPI)
	group.POST("/projects/:project/fed/resolve", postResolveAPI)
	group.POST("/projects/:project/fed/sync", postSyncAPI)
	group.POST("/projects/:project/fed/invite", postCreateInviteAPI)
	group.POST("/projects/:project/fed/join", postJoinAPI)
}

func getTransportListAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	c.JSON(http.StatusOK, transport.ListExchanges(project.Fed.GetTransport()))
}

func getTransportAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	id := c.Param("id")
	d, err := transport.GetExchange(project.Fed.GetTransport(), id)
	if err != nil {
		c.String(http.StatusNotFound, "no exchange %s: %v", id, err)
		return
	}

	c.String(http.StatusOK, string(d))
}


func putTransportAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	id := c.Param("id")
	id = strings.TrimPrefix(id, "/")
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	transport.SetExchange(project.Fed.GetTransport(), id, buf.Bytes())

	c.JSON(http.StatusOK, "")
}

func getStateAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	tm := time.Time{}
	after := c.Param("after")
	if after != "" {
		if err := json.Unmarshal([]byte(after), &tm); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
//		tm, _ = time.Parse("2006-01-02T15:04:05.999Z07:00", after)
	}

	s := project.Fed.GetState(tm)
	logrus.Debugf("Fed status in project: %v", s)
	c.JSON(http.StatusOK, s)
}

type Resolve struct {
	Extract []string `json:"extract"`
	Ignore []string `json:"ignore"`
}

func postResolveAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var resolve Resolve

	if err := c.BindJSON(&resolve); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := project.Fed.Resolve(resolve.Extract, resolve.Ignore); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func postSyncAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	project.Fed.Sync()
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

	if f, err := project.Fed.Join(request.Key, request.Token); err == os.ErrInvalid {
		c.String(http.StatusBadRequest, "Invalid token or key")
	} else if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	} else {
		project.Fed = f
		c.JSON(http.StatusOK, "")
	}

}
