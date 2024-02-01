package middlewares

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddleWareBuilder struct {
}

func (m *LoginMiddleWareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(context)
		userId := sess.Get("userId")
		if userId == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()

		//
		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || (!ok) || now.Sub(lastUpdateTime) > time.Minute {
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				println(err)
			}
		}

	}
}
