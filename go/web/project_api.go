package web

import (
	"almost-scrum/core"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ProjectMapping map[string]core.Project

var projectMapping = make(ProjectMapping)

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

func openProject(name string, path string) error {
	project, err := core.FindProject(path)
	if core.IsErr(err, "cannot open project %s from %s", name, path) {
		return err
	}

	core.ReIndex(&project)
	projectMapping[name] = project
	return nil
}

//serverRoute add routes used in a server setup
func serverRoute(group *gin.RouterGroup) {
	group.GET("/projects", listProjectsAPI)
}

//projectRoute add projects related api routes
func projectRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/info", getProjectInfoAPI)
	group.GET("/projects/:project/boards", listBoardsAPI)
	group.PUT("/projects/:project/boards/:board", createBoardAPI)
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


type ProjectInfo struct {
	SystemUser    string             `json:"system_user"`
	PropertyModel []core.PropertyDef `json:"property_model"`
}

func getProjectInfoAPI(c *gin.Context) {
	var p core.Project
	if err := getProject(c, &p); err != nil {
		return
	}

	info := ProjectInfo{
		SystemUser:    core.GetSystemUser(),
		PropertyModel: p.Config.PropertyModel,
	}
	c.JSON(http.StatusOK, info)
}

func listBoardsAPI(c *gin.Context) {
	var project core.Project

	err := getProject(c, &project)
	if err != nil {
		return
	}

	boards, err := core.ListBoards(project)
	if err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot list boards: %v", err)
		return
	}
	log.Debugf("listBoardsAPI - List boards in project: %v", boards)
	c.JSON(http.StatusOK, boards)
}

func createBoardAPI(c *gin.Context) {
	var project core.Project

	err := getProject(c, &project)
	if err != nil {
		return
	}
	board := c.Param("board")
	if err := core.CreateBoard(project, board); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot create board: %v", err)
		return
	}
	log.Debugf("createBoardAPI - Board %s created in project: %v", board, project)
	c.JSON(http.StatusCreated, board)
}

