package web

import (
	"almost-scrum/core"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func tasksRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/boards/:board", listStoryAPI)
	group.GET("/projects/:project/boards/:board/:name", getStoryAPI)
	group.POST("/projects/:project/boards/:board", postStoryAPI)
	group.POST("/projects/:project/boards/:board/:name", postStoryAPI)
	group.PUT("/projects/:project/boards/:board/:name", putStoryAPI)
}

func getRange(c *gin.Context, max int) (start int, end int) {
	startParam := c.DefaultQuery("start", "")
	endParam := c.DefaultQuery("end", "")

	start = 0
	end = max

	if startParam != "" {
		if n, err := strconv.Atoi(startParam); err == nil {
			if n < max {
				start = n
			} else {
				start = max
			}
		}
	}
	if endParam != "" {
		if n, err := strconv.Atoi(endParam); err == nil {
			if n < max {
				end = n
			}
		}
	}

	return start, end
}

func listStoryAPI(c *gin.Context) {
	var project core.Project

	err := getProject(c, &project)
	if err != nil {
		return
	}

	board := c.Param("board")
	if board == "~" {
		board = ""
	}

	filter := c.DefaultQuery("filter", "")

	infos, err := core.ListTasks(project, board, filter)
	switch err {
	case core.ErrNoFound:
		_ = c.Error(err)
		c.String(http.StatusNotFound, "Board %s does not exist", board)
	case nil:
		start, end := getRange(c, len(infos))
		infos = infos[start:end]
		c.JSON(http.StatusOK, &infos)
	default:
		c.String(http.StatusInternalServerError, "Internal Error %v", err)
	}
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
	title := c.DefaultQuery("title", "")
	move := c.DefaultQuery("move", "")

	if move == "" {
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
		return
	}

	parts := strings.Split(move, "/")
	if len(parts) != 2 {
		c.String(http.StatusBadRequest, "parameter move is invalid")
		return
	}
	oldBoard := parts[0]
	oldName := parts[1]
	id, _ := core.ExtractTaskId(oldName)
	name := oldName
	if title != "" {
		name = fmt.Sprintf("%d.%s", id, title)
	}

	if err := core.MoveTask(project, oldBoard, oldName, board, name);
		core.IsErr(err, "cannot move story %s/%s to %s/%s",
			oldBoard, oldName, board, name ) {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
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
