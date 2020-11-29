package web

import (
	"fmt"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func staticHandler(c *gin.Context) {
	p := c.Param("name")
	p = fmt.Sprintf("%s%s", "data", p)
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
		path := name[len("static"):]
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
}

//StartServer starts the embedded server.
func StartServer(port string, args []string) {
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

	r.Run(fmt.Sprintf(":%s", port))
}
