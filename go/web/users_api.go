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
	var p core.Project
	err := getProject(c, &p)
	if err != nil {
		return
	}

	users := core.GetUserList(p)
	c.JSON(http.StatusOK, users)
}

func getUserAPI(c *gin.Context) {
	var p core.Project
	err := getProject(c, &p)
	if err != nil {
		return
	}

	user := c.Param("user")
	userInfo, err := core.GetUserInfo(p, user)
	if err != nil {
		log.Warnf("Cannot get user %s info: %v", user, err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Cannot get user info")
		return
	}
	c.JSON(http.StatusOK, userInfo)
}

func putUserAPI(c *gin.Context) {
	var p core.Project
	err := getProject(c, &p)
	if err != nil {
		return
	}

	var userInfo core.UserInfo
	user := c.Param("user")
	c.BindJSON(&userInfo)

	err = core.SetUserInfo(p, user, &userInfo)
	if err != nil {
		log.Warnf("Cannot set user %s info: %v", user, err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Cannot get user info")
		return
	}
	c.JSON(http.StatusOK, "")

}
