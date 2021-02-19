package web

import (
	"almost-scrum/core"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"strings"
)

func libraryRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/library/*path", listLibraryAPI)
	group.PUT("/projects/:project/library/*path", createFolderAPI)
	group.POST("/projects/:project/library", uploadFileAPI)
	group.POST("/projects/:project/library/*path", uploadFileAPI)
	group.DELETE("/projects/:project/library/*path", deleteFileAPI)
	group.POST("/projects/:project/library-stat", getLibraryItemsAPI)
}

func localOpen(c *gin.Context, path string) {
	remoteAddress := c.Request.RemoteAddr
	if !strings.HasPrefix(remoteAddress, "localhost") &&
		!strings.HasPrefix(remoteAddress, "127.0.0.1") {
		logrus.Warn("Request for local access from remote client %s", remoteAddress)
		c.String(http.StatusBadRequest, "Local open only possible from localhost")
	} else {
		mime, _ := mimetype.DetectFile(path)
		logrus.Debugf("Mimetype for %s is %s", path, mime.String())

		for _, valid := range core.MimeForLocalAccess {
			if mime.Is(valid) {
				if err := core.OpenUrl(path); err != nil {
					c.String(http.StatusInternalServerError, "Cannot run locally: %v", err)
				} else {
					c.String(http.StatusNoContent, "Local Open")
				}
				return
			}
		}
		logrus.Warnf("Cannot open %s because mime %s is not supported", path, mime.String())
		c.String(http.StatusNotAcceptable, "Mime %s is not supported", mime.String())
	}
}

func listLibraryAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	list, path, err := core.ListLibrary(project, path)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot read path %s in library", path, err)
		return
	}

	if list == nil {
		logrus.Debugf("Library API - download path %s", path)
		_, local := c.GetQuery("local")

		if local {
			localOpen(c, path)
		} else {
			c.File(path)
		}
	} else {
		logrus.Debugf("Library API - list path %s: %v", path, list)
		c.JSON(http.StatusOK, list)
	}
}

func createFolderAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	if err := core.CreateFolderInLibrary(project, path); err != nil {
		logrus.Warnf("Cannot create folder %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot create folder %s: %v", path, err)
		return
	}
	c.String(http.StatusOK, path)
}

func uploadFileAPI(c *gin.Context) {
	logrus.Debug("Upload file request")
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	move := c.DefaultQuery("move", "")

	if move == "" {
		path, err := core.GetPathInLibrary(project, path)
		if err != nil {
			logrus.Warnf("Cannot resolve path %s: %v", path, err)
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
		logrus.Debugf("File %s uploaded in %s", file.Filename, path)
		listLibraryAPI(c)
	} else {
		if _, err := core.GetPathInLibrary(project, move); err != nil {
			logrus.Warnf("Cannot resolve path %s: %v", path, err)
			_ = c.Error(err)
			c.String(http.StatusNotFound, "Original path %s does not exist", move)
			return
		}
		if err := core.MoveFileInLibrary(project, move, path); err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot move %s to %s", move, path)
			return
		}
		c.String(http.StatusOK, "%s", path)
	}
}

func putFileAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	reader, err := c.Request.GetBody()
	if err != nil {
		logrus.Warnf("Cannot read body in save of file path %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot save file to %s", path)
		return
	}
	if err := core.SetFileInLibrary(project, path, reader); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot save file to %s", path)
		return
	}
	logrus.Debugf("File %s saved in %s", path)
	c.String(http.StatusOK, "")
}

func deleteFileAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	_, recursive := c.GetQuery("recursive")
	if err := core.DeleteFileFromLibrary(project, path, recursive); err != nil {
		logrus.Warnf("Cannot delete path %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot delete file")
		return
	}
	c.String(http.StatusOK, "")
}

func getLibraryItemsAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	var files []string
	if err := c.BindJSON(&files); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	items, err := core.GetLibraryItems(project, files)
	if err != nil {
		logrus.Warnf("Cannot get stat info for %v: %v", files, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot get files stats")
		return
	}
	c.JSON(http.StatusOK, items)
}
