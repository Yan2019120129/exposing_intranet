package router

import (
	"net/http"

	"my-base/code"

	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 校验自定义路由的后台会话和接口权限。
func AdminAuthMiddleware(c *gin.Context) {
	user, ok := getAdminUser(c)
	if !ok {
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

// SystemFileAccessMiddleware 校验系统文件列表访问权限并注入后台用户。
func SystemFileAccessMiddleware(c *gin.Context) {
	user, ok := getAdminUser(c)
	if !ok {
		return
	}
	if !auth.CheckPermissions(user, "/admin/info/system_files", http.MethodGet, nil) {
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

// getAdminUser 根据会话获取当前后台用户。
func getAdminUser(c *gin.Context) (models.UserModel, bool) {
	conn, err := code.GetGoAdminConnection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "database connection not found",
		})
		c.Abort()
		return models.UserModel{}, false
	}

	cookie, err := c.Cookie(auth.DefaultCookieKey)
	if err != nil || cookie == "" {
		c.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
		c.Abort()
		return models.UserModel{}, false
	}

	user, ok := auth.GetCurUser(cookie, conn)
	if !ok {
		c.Redirect(http.StatusFound, config.Url(config.GetLoginUrl()))
		c.Abort()
		return models.UserModel{}, false
	}
	return user, true
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
