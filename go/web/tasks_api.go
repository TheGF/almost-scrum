package web

import (
	"almost-scrum/core"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func tasksRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/boards/:board", listTaskAPI)
	group.GET("/projects/:project/boards/:board/:name", getTaskAPI)
	group.POST("/projects/:project/boards/:board", postTaskAPI)
	group.POST("/projects/:project/boards/:board/:name", postTaskAPI)
	group.PUT("/projects/:project/boards/:board/:name", putTaskAPI)
	group.DELETE("/projects/:project/boards/:board/:name", deleteTaskAPI)
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

func listTaskAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	board := c.Param("board")
	if board == "~" {
		board = ""
	}

	filter := c.DefaultQuery("filter", "")
	var keys []string
	if filter != "" {
		keys = strings.Split(filter, ",")
	} else {
		keys = []string{}
	}

	log.Debugf("Search in %s, filter=%s, keys=%v %d", board, filter, keys, len(keys))
	infos, err := core.SearchTask(project, board, true, keys...)
	switch err {
	case os.ErrNotExist:
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

func getTaskAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
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

func postTaskAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	board := c.Param("board")
	title := c.DefaultQuery("title", "")
	move := c.DefaultQuery("move", "")

	if move == "" {
		webUser := getWebUser(c)
		_, name, err := core.CreateTask(project, board, title, webUser)
		if core.IsErr(err, "cannot create task %s", title) {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		_ = core.ReIndex(project)
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
	_ = core.ReIndex(project)
	c.String(http.StatusOK, filepath.Join(board, name))
}

func putTaskAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var task core.Task
	name := c.Param("name")
	board := c.Param("board")
	if err := c.BindJSON(&task); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := core.SetTask(project, board, name, &task); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot update task %s", name)
		return
	}
	_ = core.ReIndex(project)
	c.String(http.StatusOK, "")
}

func deleteTaskAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	board := c.Param("board")
	name := c.Param("name")
	story, err := core.DeleteTask(project, board, name)
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

