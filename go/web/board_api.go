package web

import (
	"almost-scrum/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

//projectRoute add projects related api routes
func boardRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/boards", listBoardsAPI)
	group.PUT("/projects/:project/boards/:board", putBoardAPI)
	group.DELETE("/projects/:project/boards/:board", deleteBoardAPI)
}


func listBoardsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	boards, err := core.ListBoards(project)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot list boards: %v", err)
		return
	}
	logrus.Debugf("listBoardsAPI - List boards in project: %v", boards)

	c.JSON(http.StatusOK, boards)
}

func putBoardAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	board := c.Param("board")
	rename := c.DefaultQuery("rename", "")

	if rename != "" {
		if err := core.RenameBoard(project, rename, board); err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot rename board: %v", err)
			return
		}
		c.JSON(http.StatusOK, board)
		return
	}

	if err := core.CreateBoard(project, board); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot create board: %v", err)
		return
	}
	logrus.Debugf("putBoardAPI - Board %s created in project: %v", board, project)
	c.JSON(http.StatusCreated, board)
}

func deleteBoardAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	board := c.Param("board")
	if err := core.DeleteBoard(project, board); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot delete board: %v", err)
		return
	}
	logrus.Debugf("deleteBoardAPI - Board %s deleted in project: %v", board, project)
	c.JSON(http.StatusOK, "")
}
