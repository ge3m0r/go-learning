package middlewares

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddleWareBuilder struct {
}

func (m *LoginJWTMiddleWareBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}

		authCode := context.GetHeader("Authorization")
		if authCode == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != context.GetHeader("User-Agent") {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt

		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			context.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		context.Set("user", uc)
	}
}
