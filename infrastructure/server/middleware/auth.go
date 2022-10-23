package middleware

import (
	"log"
	"time"

	"DaisyClubHouse/domain/entity"
	"DaisyClubHouse/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "daisy-club-id"

type loginParams struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Authorization() *jwt.GinJWTMiddleware {
	auth, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "daisy club",
		Key:         []byte("secret"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entity.UserInfo); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &entity.UserInfo{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var params loginParams
			if err := c.ShouldBind(&params); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			username := params.Username
			password := params.Password

			if (username == "player1" && password == "password") ||
				(username == "player2" && password == "password") {
				loc, _ := time.LoadLocation("Asia/Chongqing")
				return &entity.UserInfo{
					ID:         utils.GenerateRandomID(),
					Username:   "player1",
					Nickname:   "player1",
					CreateTime: time.Now().In(loc),
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*entity.UserInfo); ok && v.Username == "player1" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT ERROR:" + err.Error())
	}

	if errInit := auth.MiddlewareInit(); errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return auth
}
