package web

import (
	"almost-scrum/core"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//projectRoute add projects related api routes
func projectRoute(group *gin.RouterGroup) {
	group.GET("/projects", listProjectsAPI)

}

func listProjectsAPI(c *gin.Context) {
	config := core.LoadConfig()

	keys := make([]string, 0, len(config.Projects))
	for k := range config.Projects {
		path := config.Projects[k]
		_, err := core.OpenProject(path)
		if err == nil {
			keys = append(keys, k)
		} else {
			log.Warnf("Project %s has invalid path %s", k, path)
		}

	}
	c.JSON(http.StatusOK, keys)
}

// GetProjectStore resolves the URL parameters
func getProject(c *gin.Context, p *core.Project) error {
	project := c.Param("project")

	config := core.LoadConfig()
	path := config.Projects[project]
	if path == "" {
		c.Error(core.ErrNoFound)
		c.String(http.StatusNotFound, "Project %s not found in configuration",
			project)
		return core.ErrNoFound
	}

	var err error
	*p, err = core.OpenProject(path)
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Project %s not found at %s: %v",
			project, path, err)
		return err
	}
	return nil
}
