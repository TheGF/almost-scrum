package web

import (
	"almost-scrum/core"
	"encoding/hex"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"
var oauth bool = false

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

func getWebUser(c *gin.Context) string {
	if oauth {
		user, _ := c.Get(identityKey)
		return user.(*User).UserName
	} else {
		return core.GetSystemUser()
	}
}

func getJWTMiddleware() *jwt.GinJWTMiddleware {
	config := core.ReadConfig()
	key, _ := hex.DecodeString(config.Secret)

	c := jwt.GinJWTMiddleware{
		Realm:       "Almost Realm",
		Key:         key,
		Timeout:     12 * time.Hour,
		MaxRefresh:  2400 * time.Hour,
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
			var login login
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := login.Username
			password := login.Password

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
		os.Exit(1)
	}

	err = middleware.MiddlewareInit()
	if err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	oauth = true
	return middleware
}
