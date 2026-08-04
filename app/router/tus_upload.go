package router

import (
	"net/http"
	"strings"

	fileService "my-base/module/file"

	adminModels "github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/gin-gonic/gin"
)

// RegisterTusUploadRouter 注册需要系统文件权限的 tus 协议路由。
func RegisterTusUploadRouter(router *gin.Engine, endpoint string, uploadHandler http.Handler) {
	if router == nil || uploadHandler == nil {
		return
	}
	endpoint = "/" + strings.Trim(endpoint, "/")
	group := router.Group(endpoint, SystemFileAccessMiddleware, tusUploadUserMiddleware)
	group.Any("/*path", gin.WrapH(uploadHandler))
}

// tusUploadUserMiddleware 将 GoAdmin 用户转换为 tus 服务可信上下文。
func tusUploadUserMiddleware(ctx *gin.Context) {
	value, exists := ctx.Get("user")
	user, ok := value.(adminModels.UserModel)
	if !exists || !ok || user.Id <= 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "unauthorized",
		})
		return
	}
	requestContext := fileService.WithTusUser(ctx.Request.Context(), fileService.TusUser{
		ID:         user.Id,
		SuperAdmin: user.IsSuperAdmin(),
	})
	ctx.Request = ctx.Request.WithContext(requestContext)
	ctx.Next()
}
