package router

import (
	"net/http"

	"my-base/code"

	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware protects custom Gin routes with GoAdmin's auth session
// and permission table rules.
func AdminAuthMiddleware(c *gin.Context) {
	conn, err := code.GetGoAdminConnection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "database connection not found",
		})
		c.Abort()
		return
	}

	cookie, err := c.Cookie(auth.DefaultCookieKey)
	if err != nil || cookie == "" {
		c.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
		c.Abort()
		return
	}

	user, ok := auth.GetCurUser(cookie, conn)
	if !ok {
		c.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
		c.Abort()
		return
	}

	_ = c.Request.ParseMultipartForm(32 << 20)
	if !checkGoAdminPermission(c, user) {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  "permission denied",
		})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}

func checkGoAdminPermission(c *gin.Context, user models.UserModel) bool {
	if auth.CheckPermissions(user, c.Request.URL.String(), c.Request.Method, c.Request.PostForm) {
		return true
	}
	if fullPath := c.FullPath(); fullPath != "" {
		return auth.CheckPermissions(user, fullPath, c.Request.Method, c.Request.PostForm)
	}
	return false
}
