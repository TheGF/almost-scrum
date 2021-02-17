package web

import (
	"almost-scrum/assets"
	"almost-scrum/core"
	"context"
	"fmt"
	"github.com/fatih/color"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

//func staticHandler(c *gin.Context) {
//	p := c.Param("name")
//	logrus.Debugf("Static request %s", p)
//	data, err := assets.Asset(p)
//	if err != nil {
//		logrus.Warnf("Cannot find resource %s: %v", p, err)
//		return
//	}
//	c.Writer.Write(data)
//}

func loadStaticContent(router *gin.Engine) {
	for _, name := range assets.AssetNames() {
		if !strings.HasPrefix(name, "build") {
			continue
		}
		path := fmt.Sprintf("/%s", name[len("build"):])
		data, err := assets.Asset(name)
		if err != nil {
			logrus.Errorf("Cannot read asset %s: %s", name, err)
			continue
		}
		router.GET(path, func(c *gin.Context) {
			m := mime.TypeByExtension(filepath.Ext(path))
			c.Data(http.StatusOK, m, data)
			logrus.Debugf("Get content of %s with mime %v", path, m)
		})
		logrus.Debugf("Bound resource %s to %s", name, path)
	}
	router.GET("/", func(c *gin.Context) {
		data, _ := assets.Asset("build/index.html")
		c.Data(http.StatusOK, "text/html", data)
	})

}

var knownClients = make([]string, 0)

type hello struct {
	Version string `json:"version"`
	Portal bool `json:"portal"`
}

func setHello(r *gin.Engine, portal bool, autoExit bool) {
	r.POST("/auth/hello", func (c *gin.Context){
		id := c.DefaultQuery("id", "")
		if id == "" {
			c.JSON(http.StatusBadRequest, "Provide an id parameter when you say hello")
			return
		}

		if autoExit {
			if _, found := core.FindStringInSlice(knownClients, id); found {
				knownClients = append(knownClients, id)
			}
		}

		c.JSON(http.StatusOK, hello{
			Version: "0.5",
			Portal:  portal,
		})
	})
}

func setBye(r *gin.Engine) {
	r.POST("/auth/bye", func (c *gin.Context) {
		id := c.DefaultQuery("id", "")
		if id == "" {
			c.String(http.StatusBadRequest, "Provide an id parameter when you say bye")
			return
		}

		if idx, found := core.FindStringInSlice(knownClients, id); found {
			knownClients = append(knownClients[0:idx], knownClients[idx+1:]...)
			if err := srv.Shutdown(ctx); err != nil {
				logrus.Fatal("Server forced to shutdown:", err)
			}

		}
		c.String(http.StatusOK, "Have a good day")
	})
}

var srv *http.Server
var ctx context.Context

func runServer(router *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("Shutting down server...")


	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Println("Server exiting")
}

//StartWeb starts the embedded web UI. Only for local usage
func StartWeb(projectPath string, port string, logLevel string, autoExit bool, args []string) {
	if strings.ToUpper(logLevel) != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	if err := openProject("~", projectPath); err != nil {
		color.Red("cannot open a project in %s: %v", projectPath, err)
		os.Exit(1)
	}

	loadStaticContent(r)
	setHello(r, false, autoExit)
	setBye(r)
	v1 := r.Group("/api/v1")

	projectRoute(v1)
	tasksRoute(v1)
	libraryRoute(v1)
	userRoute(v1)
	indexRoute(v1)
	gitRoute(v1)

	core.OpenBrowser(fmt.Sprintf("http://127.0.0.1:%s", port))
	runServer(r, fmt.Sprintf(":%s", port))
}

//StartServer starts the embedded server portal.
func StartServer(port string, logLevel string, autoExit bool, args []string) {
	if len(args) == 0 {
		color.Red("Please provide repo folder")
		os.Exit(1)
	}
	if strings.ToUpper(logLevel) != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	repoPath := args[0]

	r := gin.Default()
	loadStaticContent(r)
	setHello(r, true, autoExit)
	setBye(r)

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
	runServer(r, fmt.Sprintf(":%s", port))
}
