package web

import (
	"almost-scrum/core"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ProjectMapping map[string]core.Project
var projectMapping = make(ProjectMapping)


func openProject(name string, path string) error {
	project, err := core.FindProject(path)
	if core.IsErr(err, "cannot open project %s from %s", name, path) {
		return err
	}

	projectMapping[name] = project
	return nil
}

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

// getProject resolves the URL parameters
func getProject(c *gin.Context, p *core.Project) error {
	name := c.Param("project")

	if project, found := projectMapping[name]; found {
		*p = project
		return nil
	} else {
		_ = c.Error(core.ErrNoFound)
		c.String(http.StatusNotFound, "Project %s not found in configuration", name)
		return core.ErrNoFound
	}
}
