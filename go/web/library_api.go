package web

import (
	"almost-scrum/core"
	"almost-scrum/library"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"net/http"
	"path/filepath"
	"strings"
)

func libraryRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/library/*path", getLibraryAPI)
	group.GET("/projects/:project/archive/*path", getArchiveAPI)
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
				if err := open.Start(path); err != nil {
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

func getLibraryAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	_, versions := c.GetQuery("versions")
	_, local := c.GetQuery("local")

	if versions {
		items, err := library.GetPreviousVersions(project, path)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, items)
		return
	}

	isDir, err := library.IsDir(project, path)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot read path %s in library", path, err)
		return
	}

	if isDir {
		if items, err := library.List(project, path); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, items)
		}
		return
	}

	absPath, err := library.AbsPath(project, path)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	} else if local {
		localOpen(c, absPath)
	} else {
		c.File(absPath)
	}
}

func getArchiveAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	_, local := c.GetQuery("local")

	absPath, err := library.ArchivePath(project, path)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	} else if local {
		localOpen(c, absPath)
	} else {
		c.File(absPath)
	}
}


func createFolderAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	if err := library.CreateFolder(project, path, getWebUser(c)); err != nil {
		logrus.Warnf("Cannot create folder %s: %v", path, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "cannot create folder %s: %v", path, err)
		return
	}
	c.String(http.StatusOK, path)
}

func uploadFileAPI(c *gin.Context) {
	logrus.Debug("Push file request")
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	move := c.DefaultQuery("move", "")

	if move == "" {
		file, _ := c.FormFile("file")
		reader, err := file.Open()
		if err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot open upload stream: %v", err)
			return
		}
		defer reader.Close()

		dest := filepath.Join(path, file.Filename)
		if _, err = library.SetFileInLibrary(project, dest, reader, getWebUser(c)); err != nil {
			_ = c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot upload to %s", path)
			return
		}
		logrus.Debugf("File %s uploaded in %s", file.Filename, path)
		getLibraryAPI(c)
	} else {
		if err := library.MoveFile(project, move, path); err != nil {
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
	if path, err := library.SetFileInLibrary(project, path, reader, getWebUser(c)); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot save file to %s", path)
		return
	}
	logrus.Debugf("File %s saved in %s", path)
	c.String(http.StatusOK, path)
}

func deleteFileAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	path := c.Param("path")
	_, recursive := c.GetQuery("recursive")
	if err := library.DeleteFile(project, path, recursive); err != nil {
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

	items, err := library.GetItems(project, files)
	if err != nil {
		logrus.Warnf("Cannot get stat info for %v: %v", files, err)
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot get files stats")
		return
	}
	c.JSON(http.StatusOK, items)
}
