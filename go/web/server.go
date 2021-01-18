package web

import (
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
	data, err := Asset(p)
	if err != nil {
		log.Warnf("Cannot find resource %s: %v", p, err)
		return
	}
	c.Writer.Write(data)
}

func loadStaticContent(router *gin.Engine) {
	for _, name := range AssetNames() {
		path := fmt.Sprintf("/%s", name[len("build"):])
		data, err := Asset(name)
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
		data, _ := Asset("build/index.html")
		c.Data(http.StatusOK, "text/html", data)
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

	loadStaticContent(r)
	v1 := r.Group("/api/v1")
	projectRoute(v1)
	tasksRoute(v1)
	libraryRoute(v1)
	userRoute(v1)
	indexRoute(v1)
	gitRoute(v1)

	r.Run(fmt.Sprintf(":%s", port))
}

//StartServer starts the embedded server portal.
func StartServer(projectPath string, port string, logLevel string, args []string) {
	if strings.ToUpper(logLevel) != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	loadStaticContent(r)

	authMiddleware := getJWTMiddleware()
	r.POST("/auth/login", authMiddleware.LoginHandler)

	v1 := r.Group("/api/v1")
	v1.GET("/refresh_token", authMiddleware.RefreshHandler)
	v1.Use(authMiddleware.MiddlewareFunc())

	serverRoute(v1)
	projectRoute(v1)
	tasksRoute(v1)
	libraryRoute(v1)
	userRoute(v1)
	indexRoute(v1)
	gitRoute(v1)

	r.Run(fmt.Sprintf(":%s", port))
}
