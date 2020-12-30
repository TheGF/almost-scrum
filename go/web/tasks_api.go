package web

import (
	"almost-scrum/core"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func storeRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/boards", listStoresAPI)
	group.GET("/projects/:project/boards/:board", getStoryAPI)
	group.GET("/projects/:project/boards/:board/:name", getStoryAPI)
	group.POST("/projects/:project/boards/:board", postStoryAPI)
	group.POST("/projects/:project/boards/:board/:name", postStoryAPI)
	group.PUT("/projects/:project/boards/:board/:name", putStoryAPI)
}


func listStoresAPI(c *gin.Context) {
	var p core.Project

	err := getProject(c, &p)
	if err != nil {
		return
	}

	boards, err := core.ListBoards(p)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot list boards: %v", err)
		return
	}
	log.Debugf("listStoresAPI - List boards in project: %v", boards)
	c.JSON(http.StatusOK, boards)
}

func getStoryAPI(c *gin.Context) {
	var project core.Project

	err := getProject(c, &project)
	if err != nil {
		return
	}

	board := c.Param("board")
	name := c.Param("name")
	story, err := core.GetTask(project, board, name)
	switch err {
	case core.ErrNoFound:
		_ = c.Error(err)
		c.String(http.StatusNotFound, "Task %s/%s does not exist", board, name)
	case nil:
		c.JSON(http.StatusOK, story)
	default:
		c.String(http.StatusInternalServerError, "Internal Error %v", err)
	}
}

func postStoryAPI(c *gin.Context) {
	var project core.Project

	if err := getProject(c, &project); err != nil {
		return
	}
	board := c.Param("board")

	var story core.Task
	title := c.DefaultQuery("title", "noname")

	if err := c.BindJSON(&story); err != nil {
		log.Warnf("Invalid JSON in request: %v", err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	name := core.NewTaskName(project, title)
	if err := core.SetTask(project, board, name, &story); core.IsErr(err, "cannot save story to %s", name) {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, name)
}

func putStoryAPI(c *gin.Context) {
	var project core.Project

	err := getProject(c, &project)
	if err != nil {
		return
	}

	var task core.Task
	name := c.Param("name")
	board := c.Param("board")
	if err = c.BindJSON(&task); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err = core.SetTask(project, board, name, &task); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot update task %s", name)
		return
	}
	c.String(http.StatusOK, "")
}
