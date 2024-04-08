package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const userKey = "Secret"

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userSession := session.Get(userKey)
		if userSession == nil {
			c.Redirect(
				http.StatusFound,
				"/shines/main/login-page",
			)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func SetSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(userKey))
	return sessions.Sessions("mySession", store)
}

func SaveSession(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Set(userKey, username)
	session.Save()
}

func GetSession(c *gin.Context) string {
	session := sessions.Default(c)
	userSession := session.Get(userKey)
	if userSession == nil {
		return "Anonymous"
	} else {
		return userSession.(string)
	}
}

func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userKey)
	session.Clear()
	session.Save()
}

func CheckSession(c *gin.Context) bool {
	session := sessions.Default(c)
	userSession := session.Get(userKey)
	if userSession == nil {
		return false
	} else {
		return true
	}
}
