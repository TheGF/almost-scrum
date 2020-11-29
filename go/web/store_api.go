package web

import (
	"almost-scrum/core"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func storeRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/stores", listStoresAPI)
	group.GET("/projects/:project/stores/:store", listStoreAPI)
	group.GET("/projects/:project/stores/:store/*path", getStoryAPI)
	group.POST("/projects/:project/stores/:store", postStoryAPI)
	group.POST("/projects/:project/stores/:store/*path", postStoryAPI)
	group.PUT("/projects/:project/stores/:store/*path", putStoryAPI)
}

// GetProjectStore resolves the URL parameters
func getProjectAndStore(c *gin.Context, p *core.Project, s *core.Store) error {
	err := getProject(c, p)
	if err != nil {
		return err
	}
	if s == nil {
		return nil
	}
	store := c.Param("store")
	*s, err = core.GetStore(*p, store)
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Store %s not found in project: %v",
			store, err)
		return err
	}

	return nil
}

func listStoresAPI(c *gin.Context) {
	var p core.Project
	project := c.Param("project")

	err := getProjectAndStore(c, &p, nil)
	if err != nil {
		return
	}

	stores, err := core.ListStores(p)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot list stores: %v", err)
		return
	}
	log.Debugf("listStoresAPI - List stores in project %v: %v", project, stores)
	c.JSON(http.StatusOK, stores)
}

func listStoreAPI(c *gin.Context) {
	var p core.Project
	var s core.Store

	err := getProjectAndStore(c, &p, &s)
	if err != nil {
		return
	}

	stores := core.ListStore(s)
	log.Debugf("listStoreAPI - List stores: %v", stores)
	c.JSON(http.StatusOK, stores)
}

func getStoryAPI(c *gin.Context) {
	var p core.Project
	var s core.Store

	err := getProjectAndStore(c, &p, &s)
	if err != nil {
		return
	}

	path := c.Param("path")
	if !strings.HasSuffix(path, ".story") {
		c.String(http.StatusNotFound, "Story %s does not exist", path)
		return
	}

	story, err := core.GetStory(s, path)
	switch err {
	case core.ErrNoFound:
		c.Error(err)
		c.String(http.StatusNotFound, "Story %s does not exist", path)
	case nil:
		c.JSON(http.StatusOK, story)
	default:
		c.String(http.StatusInternalServerError, "Internal Error %v", err)
	}
}

func postStoryAPI(c *gin.Context) {
	var p core.Project
	var s core.Store

	err := getProjectAndStore(c, &p, &s)
	if err != nil {
		return
	}

	var story core.Story
	path := c.Param("path")
	title := c.DefaultQuery("title", "noname")

	err = c.BindJSON(&story)
	if err != nil {
		log.Warnf("Invalid JSON in request: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := core.GetNextID(p)
	path = fmt.Sprintf("%s/%d. %s.story", path, id, title)
	err = core.SetStory(s, path, &story)
	if err != nil {
		log.Warnf("createStory - Cannot save story to %s: %v", path, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, path)
}

func putStoryAPI(c *gin.Context) {
	var p core.Project
	var s core.Store

	err := getProjectAndStore(c, &p, &s)
	if err != nil {
		return
	}

	var story core.Story
	path := c.Param("path")
	if strings.HasSuffix(path, ".story") {
		err = c.BindJSON(&story)
		if err != nil {
			log.Warnf("Invalid JSON in request: %v", err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err = core.SetStory(s, path, &story)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot update story %s", path)
		}
		c.String(http.StatusOK, "")
	} else {
		err = core.CreateFolder(s, path)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, "Cannot create folder %s", path)
		}
		c.String(http.StatusOK, "")
	}
}
