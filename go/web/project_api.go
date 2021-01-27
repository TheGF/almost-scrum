package web

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProjectMapping map[string]*core.Project

var projectMapping = make(ProjectMapping)

// getProject resolves the URL parameters
func getProject(c *gin.Context) *core.Project {
	name := c.Param("project")

	if project, found := projectMapping[name]; found {
		return project
	} else {
		_ = c.Error(core.ErrNoFound)
		c.String(http.StatusNotFound, "Project %s not found in configuration", name)
		return nil
	}
}

func openProject(name string, path string) error {
	project, err := core.FindProject(path)
	if core.IsErr(err, "cannot open project %s from %s", name, path) {
		return err
	}

	_ = core.ReIndex(project)
	projectMapping[name] = project
	return nil
}

func searchForProject(repoPath string, name string, items ...string) bool {
	elem := []string{repoPath, name}
	p := filepath.Join(append(elem, items...)...)
	if _, err := os.Stat(p); err != nil {
		return false
	}
	if err := openProject(name, p); err != nil {
		logrus.Warnf("Cannot open project %s in path %s: %v", name, p, err)
	} else {
		logrus.Infof("Added project %s (%s) to repo", name, p)
	}
	return true
}

func loadRepoProjects(repoPath string) error {
	infos, err := ioutil.ReadDir(repoPath)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		name := info.Name()

		if searchForProject(repoPath, name, core.ProjectConfigFile) {
			continue
		} else {
			searchForProject(repoPath, name, core.ProjectFolder, core.ProjectConfigFile)
		}
	}
	return nil
}

func loadGlobalProjects() {
	config := core.LoadConfig()
	for name := range config.Projects {
		path := config.Projects[name]
		_, err := core.OpenProject(path)
		if err == nil {
			logrus.Infof("Added project %s (%s) to repo", name, path)
		} else {
			logrus.Warnf("Project %s has invalid path %s", name, path)
		}
	}

}

//serverRoute add routes used in a server setup
func serverRoute(group *gin.RouterGroup, repoPath string) {
	group.GET("/projects", listProjectsAPI)

	if err := loadRepoProjects(repoPath); err != nil {
		color.Red("Cannot start server because of invalid repo folder %s: %v", repoPath, err)
		os.Exit(1)
	}
	loadGlobalProjects()
}

//projectRoute add projects related api routes
func projectRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/info", getProjectInfoAPI)
	group.GET("/projects/:project/boards", listBoardsAPI)
	group.PUT("/projects/:project/boards/:board", createBoardAPI)
}

func listProjectsAPI(c *gin.Context) {
	var names []string
	for key := range projectMapping {
		names = append(names, key)
	}

	c.JSON(http.StatusOK, names)
}

type ProjectInfo struct {
	SystemUser    string             `json:"systemUser"`
	LoginUser     string             `json:"loginUser"`
	CurrentBoard  string             `json:"currentBoard"`
	PropertyModel []core.PropertyDef `json:"propertyModel"`
	GitProject    bool               `json:"gitProject"`
}

func getProjectInfoAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	gitFolder := filepath.Join(filepath.Dir(project.Path), core.GitFolder)
	_, err := os.Stat(gitFolder)
	gitProject := err == nil

	info := ProjectInfo{
		SystemUser:    core.GetSystemUser(),
		LoginUser:     core.GetSystemUser(),
		CurrentBoard:  project.Config.CurrentBoard,
		PropertyModel: project.Config.PropertyModel,
		GitProject:    gitProject,
	}
	c.JSON(http.StatusOK, info)
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

func createBoardAPI(c *gin.Context) {
	var project *core.Project
	if project = getProject(c); project == nil {
		return
	}

	board := c.Param("board")
	if err := core.CreateBoard(project, board); err != nil {
		_ = c.Error(err)
		c.String(http.StatusInternalServerError, "Cannot create board: %v", err)
		return
	}
	logrus.Debugf("createBoardAPI - Board %s created in project: %v", board, project)
	c.JSON(http.StatusCreated, board)
}
