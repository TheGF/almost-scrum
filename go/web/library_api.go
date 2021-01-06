package web

import (
	"almost-scrum/core"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
)

func libraryRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/library/*path", listLibraryAPI)
	group.PUT("/projects/:project/library/*path", createFolderAPI)
	group.POST("/projects/:project/library/*path", uploadFileAPI)
	group.DELETE("/projects/:project/library/*path", deleteFileAPI)
}

func listLibraryAPI(c *gin.Context) {
	var p core.Project

	err := getProject(c, &p)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusNotFound, "Project not found")
		return
	}

	path := c.Param("path")
	list, path, err := core.ListLibrary(p, path)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot read path %s in library", path, err)
		return
	}

	if list == nil {
		log.Debugf("Library API - download path %s", path)
		c.File(path)
	} else {
		log.Debugf("Library API - list path %s: %v", path, list)
		c.JSON(http.StatusOK, list)
	}
}

func createFolderAPI(c *gin.Context) {
	var project core.Project

	err := getProject(c, &project)
	if err != nil {
		return
	}

	path := c.Param("path")
	if err := core.CreateFolderInLibrary(project, path); err != nil {
		log.Warnf("Cannot create folder %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot create folder %s: %v", path, err)
		return
	}
	c.String(http.StatusOK, path)
}

func uploadFileAPI(c *gin.Context) {
	var p core.Project

	err := getProject(c, &p)
	if err != nil {
		return
	}

	path := c.Param("path")
	move := c.DefaultQuery("move", "")

	if move == "" {
		path, err = core.GetPathInLibrary(p, path)
		if err != nil {
			log.Warnf("Cannot resolve path %s: %v", path, err)
			_ = c.Error(err)
			c.String(http.StatusNotFound, "Path not found")
			return
		}

		file, _ := c.FormFile("file")
		if err = c.SaveUploadedFile(file, filepath.Join(path, file.Filename)); err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot upload to %s", path)
			return
		}
		listLibraryAPI(c)
	} else {
		if _, err = core.GetPathInLibrary(p, move); err != nil {
			log.Warnf("Cannot resolve path %s: %v", path, err)
			_ = c.Error(err)
			c.String(http.StatusNotFound, "Original path %s does not exist", move)
			return
		}
		if err = core.MoveFileInLibrary(p, move, path); err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot move %s to %s", move, path)
			return
		}
		c.String(http.StatusOK, "%s", path)
	}
}

func deleteFileAPI(c *gin.Context) {
	var p core.Project

	err := getProject(c, &p)
	if err != nil {
		return
	}

	path := c.Param("path")
	if err = core.DeleteFileFromLibrary(p, path); err != nil {
		log.Warnf("Cannot delete path %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot delete file")
		return
	}
	c.String(http.StatusOK, "")
}
