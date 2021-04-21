package web

import (
	"almost-scrum/core"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func userRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/users", listProjectUsersAPI)
	group.GET("/projects/:project/users/:user", getProjectUserAPI)
	group.PUT("/projects/:project/users/:user", putProjectUserAPI)
	group.DELETE("/projects/:project/users/:user", delProjectUserAPI)
	group.GET("/passwords", listLocalUsersAPI)
	group.POST("/passwords", registerLocalUserAPI)
}

func listProjectUsersAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}
	users := core.GetUserList(project)

	projectLock.Lock()
	defer projectLock.Unlock()
	projectUsers[c.Param("project")] = users
	c.JSON(http.StatusOK, users)
}

func getProjectUserAPI(c *gin.Context) {
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

func putProjectUserAPI(c *gin.Context) {
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

func delProjectUserAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	user := c.Param("user")
	err := core.DelUserInfo(project, user)
	if err != nil {
		log.Warnf("Cannot del user %s info: %v", user, err)
		_ = c.Error(err)
		c.JSON(http.StatusInternalServerError, "Cannot del user info")
		return
	}
	c.JSON(http.StatusOK, "")
}

func getWebUserAPI(c *gin.Context) {
	c.JSON(http.StatusOK, getWebUser(c))
}

func listLocalUsersAPI(c *gin.Context) {
	var users []string

	config := core.ReadConfig()

	for user := range config.Passwords {
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerLocalUserAPI(c *gin.Context) {
	webUser := getWebUser(c)
	var credentials Credentials
	if err := c.BindJSON(&credentials); err != nil || credentials.Username == "" {
		log.Warnf("Invalid credentials: %v", err)
		c.String(http.StatusBadRequest, "Invalid data for username and password")
		return
	}

	config := core.ReadConfig()
	if _, found := config.Passwords[credentials.Username]; found && credentials.Username != webUser &&
		credentials.Password != "" && credentials.Password != "changeme" {
		c.String(http.StatusForbidden, "User '%s' cannot change another user ('%s') password",
			getWebUser(c), credentials.Username)
		return
	}

	if credentials.Password == "" && credentials.Username == webUser {
		c.String(http.StatusForbidden, "Suicide is illegal for user '%s'", webUser)
		return
	}
	if err := core.SetPassword(credentials.Username, credentials.Password); err != nil {
		log.Warnf("Cannot set credentials for user '%s': %v", credentials.Username, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot set credentials for user '%s", credentials.Username)
		return
	}

	c.String(http.StatusOK, "")
}