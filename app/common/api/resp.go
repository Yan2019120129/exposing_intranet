package api

import (
	"github.com/gin-gonic/gin"
	"my-base/app/common/response"
	"net/http"
)

type Api struct {
	val     response.Resp
	context *gin.Context
}

// MakeContext 构建上下文 服务
func (a *Api) MakeContext(c *gin.Context) *Api {
	a.context = c
	return a
}

// Bind 绑定结构体
func (a *Api) Bind(val any) *Api {
	err := a.context.Bind(val)
	if err != nil {
		a.val.Set(http.StatusBadRequest, val, err.Error())
		a.context.JSON(http.StatusBadRequest, a.val)
	}
	return a
}

// SetService 设置服务
func (a *Api) SetService(val any, err error) {
	if err != nil {
		a.fail(err)
		return
	}
	a.success(val)
}

// success 返回结构体
func (a *Api) success(vales any) {
	a.val.Set(http.StatusOK, vales, "")
	a.context.JSON(http.StatusBadRequest, a.val)
}

// success 返回结构体
func (a *Api) fail(err error) {
	a.val.Set(http.StatusInternalServerError, "fail", err.Error())
	a.context.JSON(http.StatusInternalServerError, a.val)
}
