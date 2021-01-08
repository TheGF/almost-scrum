package web

import (
	"almost-scrum/core"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func userRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/users", listUsersAPI)
	group.GET("/projects/:project/users/:user", getUserAPI)
	group.PUT("/projects/:project/users/:user", putUserAPI)
}

func listUsersAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	users := core.GetUserList(project)
	c.JSON(http.StatusOK, users)
}

func getUserAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	user := c.Param("user")
	userInfo, err := core.GetUserInfo(project, user)
	if err != nil {
		log.Warnf("Cannot get user %s info: %v", user, err)
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, "Cannot get user info")
		return
	}
	c.JSON(http.StatusOK, userInfo)
}

func putUserAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var userInfo core.UserInfo
	user := c.Param("user")
	_ = c.BindJSON(&userInfo)

	err := core.SetUserInfo(project, user, &userInfo)
	if err != nil {
		log.Warnf("Cannot set user %s info: %v", user, err)
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, "Cannot get user info")
		return
	}
	c.JSON(http.StatusOK, "")

}
