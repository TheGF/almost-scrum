package web

import (
	"almost-scrum/core"
	"encoding/hex"
	"time"

	log "github.com/sirupsen/logrus"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func getJWTMiddleware() *jwt.GinJWTMiddleware {
	config := core.LoadConfig()
	key, _ := hex.DecodeString(config.Secret)

	c := jwt.GinJWTMiddleware{
		Realm:       "Almost Realm",
		Key:         key,
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if core.CheckUser(userID, password) {
				return &User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*User); ok {
				return true
			}

			log.Warnf("Failed authentication for user %v", data)
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token",
		TimeFunc:    time.Now,
	}

	middleware, err := jwt.New(&c)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	err = middleware.MiddlewareInit()
	if err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	return middleware
}