package web

import (
	"almost-scrum/assets"
	"almost-scrum/core"
	"fmt"
	"github.com/fatih/color"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func staticHandler(c *gin.Context) {
	p := c.Param("name")
	//	p = fmt.Sprintf("%s%s", "data", p)
	log.Debugf("Static request %s", p)
	data, err := assets.Asset(p)
	if err != nil {
		log.Warnf("Cannot find resource %s: %v", p, err)
		return
	}
	c.Writer.Write(data)
}

func loadStaticContent(router *gin.Engine) {
	for _, name := range assets.AssetNames() {
		if !strings.HasPrefix(name, "build") {
			continue
		}
		path := fmt.Sprintf("/%s", name[len("build"):])
		data, err := assets.Asset(name)
		if err != nil {
			log.Errorf("Cannot read asset %s: %s", name, err)
			continue
		}
		router.GET(path, func(c *gin.Context) {
			m := mime.TypeByExtension(filepath.Ext(path))
			c.Data(http.StatusOK, m, data)
			log.Debugf("Get content of %s with mime %v", path, m)
		})
		log.Debugf("Bound resource %s to %s", name, path)
	}
	router.GET("/", func(c *gin.Context) {
		data, _ := assets.Asset("build/index.html")
		c.Data(http.StatusOK, "text/html", data)
	})

}

func setPortal(r *gin.Engine, portal bool) {
	r.GET("/auth/portal", func (c *gin.Context){
		c.JSON(http.StatusOK, portal)
	})
}

//StartWeb starts the embedded web UI. Only for local usage
func StartWeb(projectPath string, port string, logLevel string, args []string) {
	if strings.ToUpper(logLevel) != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	if err := openProject("~", projectPath); err != nil {
		color.Red("cannot open a project in %s: %v", projectPath, err)
		os.Exit(1)
	}

	setPortal(r, false)
	loadStaticContent(r)
	v1 := r.Group("/api/v1")

	projectRoute(v1)
	tasksRoute(v1)
	libraryRoute(v1)
	userRoute(v1)
	indexRoute(v1)
	gitRoute(v1)

	core.OpenBrowser(fmt.Sprintf("http://127.0.0.1:%s", port))
	r.Run(fmt.Sprintf(":%s", port))
}

//StartServer starts the embedded server portal.
func StartServer(port string, logLevel string, args []string) {
	if len(args) == 0 {
		color.Red("Please provide repo folder")
		os.Exit(1)
	}
	if strings.ToUpper(logLevel) != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	repoPath := args[0]

	r := gin.Default()
	setPortal(r, true)
	loadStaticContent(r)

	authMiddleware := getJWTMiddleware()
	r.POST("/auth/login", authMiddleware.LoginHandler)

	v1 := r.Group("/api/v1")
	v1.GET("/refresh_token", authMiddleware.RefreshHandler)
	v1.Use(authMiddleware.MiddlewareFunc())

	serverRoute(v1, repoPath)
	projectRoute(v1)
	tasksRoute(v1)
	libraryRoute(v1)
	userRoute(v1)
	indexRoute(v1)
	gitRoute(v1)

	core.OpenBrowser(fmt.Sprintf("http://127.0.0.1:%s", port))
	r.Run(fmt.Sprintf(":%s", port))
}
