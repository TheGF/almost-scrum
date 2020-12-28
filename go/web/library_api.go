package web

import (
	"almost-scrum/core"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func libraryRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/library/*path", listLibraryAPI)
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

func uploadFileAPI(c *gin.Context) {
	var p core.Project

	err := getProject(c, &p)
	if err != nil {
		return
	}

	path := c.Param("path")
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
	}
	listLibraryAPI(c)
}

func deleteFileAPI(c *gin.Context) {
	var p core.Project

	err := getProject(c, &p)
	if err != nil {
		return
	}

	path := c.Param("path")
	path, err = core.GetPathInLibrary(p, path)
	if err != nil {
		log.Warnf("Cannot resolve path %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusNotFound, "Path not found")
		return
	}

	err = os.Remove(path)
	if err != nil {
		log.Warnf("Cannot delete path %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot delete file")
		return
	}
	c.String(http.StatusOK, "")
}
