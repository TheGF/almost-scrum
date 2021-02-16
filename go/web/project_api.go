package web

import (
	"almost-scrum/core"
	"fmt"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ProjectMapping map[string]*core.Project
type ProjectUsers map[string][]string

var projectMapping = make(ProjectMapping)
var projectUsers = make(ProjectUsers)

type NoAccess struct {
	Message string `json:"message"`
	Users []string `json:"users"`
}

// getProject resolves the URL parameters
func getProject(c *gin.Context) *core.Project {
	name := c.Param("project")

	project, found := projectMapping[name]
	if !found {
		_ = c.Error(core.ErrNoFound)
		c.String(http.StatusNotFound, "Project %s not found in configuration", name)
		return nil
	}

	user := getWebUser(c)
	users := projectUsers[name]
	if _, found := core.FindStringInSlice(users, user); !found {
		noAccess := NoAccess{
			Message: "No access to project",
			Users: users,
		}
		logrus.Warnf("User %s has no access to project %s. Valid users [%s]", user, name,
			strings.Join(users, " "))
		c.JSON(http.StatusForbidden, noAccess)
		return nil
	}

	return project
}

func openProject(name string, path string) error {
	project, err := core.FindProject(path)
	if core.IsErr(err, "cannot open project %s from %s", name, path) {
		return err
	}

	_ = core.ReIndex(project)
	projectMapping[name] = project
	users := core.GetUserList(project)
	projectUsers[name] = users
	logrus.Infof("Open project %s for users [%s]", name, strings.Join(users, " "))
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
		return false
	}
	logrus.Infof("Added project %s (%s) to repo", name, p)
	return true
}

var repository string

func loadRepoProjects(repoPath string) error {
	infos, err := ioutil.ReadDir(repoPath)
	if err != nil {
		return err
	}

	repository = repoPath
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
	group.POST("/projects", createProjectAPI)
	group.GET("/user", getWebUserAPI)
	group.GET("/templates", listProjectTemplatesAPI)

	if err := loadRepoProjects(repoPath); err != nil {
		color.Red("Cannot start server because of invalid repo folder %s: %v", repoPath, err)
		os.Exit(1)
	}
	loadGlobalProjects()
}

func listProjectsAPI(c *gin.Context) {
	var names = make([]string, 0)
	for key := range projectMapping {
		names = append(names, key)
	}

	c.JSON(http.StatusOK, names)
}

type Authorization struct {
	Users []string
	Pam   bool
}


func listProjectTemplatesAPI(c *gin.Context) {
	c.JSON(http.StatusOK, core.ListProjectTemplates())
}

type createOptions struct {
	ProjectName  string   `json:"projectName"`
	ImportFolder string   `json:"importPath"`
	Templates    []string `json:"templates"`
	Inject       bool     `json:"inject"`
	GitUrl       string   `json:"gitUrl"`
}

func createProjectInFolder(c *gin.Context, name string, path string, templates []string, res string)  {
	var err error
	var project *core.Project

	if len(templates) == 0 {
		project, err = core.InitProject(path)
	} else {
		project, err = core.InitProjectFromTemplate(path, templates)
	}
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := getWebUser(c)
	if err = core.SetUserInfo(project, user, &core.UserInfo{}); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Cannot set user '%s' in project '%s'", user, name))
		return
	}
	projectMapping[name] = project
	projectUsers[name] = []string{user}

	logrus.Infof("Project %s created by user '%s'", name, user)
	c.String(http.StatusCreated, res)
}

var cloneUrl = regexp.MustCompile(`http.*\/([^\/]*?)(.git)?$`)
func createProjectFromGit(c *gin.Context, createOptions createOptions) {
	match := cloneUrl.FindStringSubmatch(createOptions.GitUrl)
	if len(match) < 2 {
		c.String(http.StatusBadRequest, "Unsupported URL '%s'", createOptions.GitUrl)
		return
	}
	name := match[1]

	gitClient := core.GetGitClientFromGlobalConfig()
	responseBody, err := gitClient.Clone(createOptions.GitUrl, repository)
	if err != nil {
		c.String(http.StatusInternalServerError, responseBody)
	}

	path := filepath.Join(repository, name, core.ProjectFolder)
	if err := openProject(name, path); err == nil {
		c.String(http.StatusCreated, name)
	} else if createOptions.Inject {
		createProjectInFolder(c, name, path, createOptions.Templates, responseBody)
	} else {
		c.String(http.StatusBadRequest, "Folder %s does not contain a valid project",
			createOptions.ImportFolder)
	}
}

func importFolderFromPath(c *gin.Context, createOptions createOptions) {
	path := createOptions.ImportFolder
	if _, err := os.Stat(path); err != nil {
		c.String(http.StatusBadRequest, "Folder %s does not exists", createOptions.ImportFolder)
		return
	}

	name := filepath.Base(createOptions.ImportFolder)
	if err := openProject(name, createOptions.ImportFolder); err == nil {
		c.String(http.StatusCreated, name)
	} else if createOptions.Inject {
		createProjectInFolder(c, name, path, createOptions.Templates, name)
	} else {
		c.String(http.StatusBadRequest, "Folder %s does not contain a valid project",
			createOptions.ImportFolder)
	}
}

func createProjectFromScratch(c *gin.Context, createOptions createOptions) {
	name := createOptions.ProjectName
	path := filepath.Join(repository, name)
	if _, err := os.Stat(path); err == nil {
		c.String(http.StatusConflict, "Project %s already exists", createOptions.ProjectName)
		return
	}
	createProjectInFolder(c, name, path, createOptions.Templates, name)
}

func createProjectAPI(c *gin.Context) {
	var createOptions createOptions

	if err := c.BindJSON(&createOptions); core.IsErr(err, "Invalid JSON") {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if createOptions.GitUrl != "" {
		createProjectFromGit(c, createOptions)
	} else if createOptions.ImportFolder != "" {
		importFolderFromPath(c, createOptions)
	} else if createOptions.ProjectName != "" {
		createProjectFromScratch(c, createOptions)
	} else {
		c.String(http.StatusBadRequest, "Invalid parameters")
	}
}

//projectRoute add projects related api routes
func projectRoute(group *gin.RouterGroup) {
	group.GET("/projects/:project/info", getProjectInfoAPI)
	group.GET("/projects/:project/boards", listBoardsAPI)
	group.PUT("/projects/:project/boards/:board", createBoardAPI)
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
		LoginUser:     getWebUser(c),
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
