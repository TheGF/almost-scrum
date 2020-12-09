package web

import (
	"fmt"
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
		path := fmt.Sprintf("/%s", name)
		data, err := Asset(name)
		if err != nil {
			log.Errorf("Cannot read asset %s: %s", name, err)
			continue
		}
		router.GET(path, func(c *gin.Context) {
			c.Writer.Write(data)
		})
		log.Debugf("Bound resource %s to %s", name, path)
	}
	router.GET("/", func(c *gin.Context) {
		data, _ := Asset("index.html")
		c.Writer.Write(data)
	})

}

//StartServer starts the embedded server.
func StartServer(port string, logLevel string, args []string) {
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
	projectRoute(v1)
	storeRoute(v1)
	libraryRoute(v1)
	userRoute(v1)

	r.Run(fmt.Sprintf(":%s", port))
}
