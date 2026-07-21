package apis

import (
	"errors"
	"net/http"

	"my-base/app/service"
	"my-base/code/api"
	"my-base/code/contract"

	"github.com/gin-gonic/gin"
)

type Port struct {
	api.Api
}

type PortRequest = contract.PortRequest
type PortResponse = contract.PortResponse
type PortMappingInfo = contract.PortMappingInfo

func (e Port) Manage(ctx *gin.Context) {
	e.MakeContext(ctx)

	var req PortRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, PortResponse{Success: false, Message: "参数错误"})
		return
	}
	if err := e.MakeOrm().Errors; err != nil {
		ctx.JSON(http.StatusInternalServerError, PortResponse{Success: false, Message: "服务器内部错误"})
		return
	}

	result, err := service.NewPortService(e.Orm).Manage(service.PortCommand{
		Symbol:     req.Symbol,
		Action:     req.Action,
		ServerPort: req.ServerPort,
		LocalAddr:  req.LocalAddr,
		Comment:    req.Comment,
	})
	if err != nil {
		status := portErrorStatus(err)
		if status == http.StatusInternalServerError && e.Logger != nil {
			e.Logger.Error("manage port mapping:", err)
		}
		ctx.JSON(status, PortResponse{Success: false, Message: portErrorMessage(err)})
		return
	}

	data := make([]PortMappingInfo, 0, len(result.Mappings))
	for _, mapping := range result.Mappings {
		data = append(data, PortMappingInfo{
			ServerPort: mapping.ServerPort,
			LocalAddr:  mapping.LocalAddr,
			Comment:    mapping.Comment,
			Status:     mapping.Status,
		})
	}

	switch req.Action {
	case "add":
		ctx.JSON(http.StatusOK, PortResponse{Success: true, Message: "端口映射添加成功", Data: data})
	case "del":
		ctx.JSON(http.StatusOK, PortResponse{Success: true, Message: "端口映射删除成功"})
	default:
		ctx.JSON(http.StatusOK, PortResponse{Success: true, Message: "获取成功", Data: data})
	}
}

func portErrorStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrClientNotRegistered):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrPortNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrInvalidPort), errors.Is(err, service.ErrInvalidAction):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrPortConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func portErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrClientNotRegistered):
		return "客户端未注册"
	case errors.Is(err, service.ErrPortNotFound):
		return "端口映射不存在"
	case errors.Is(err, service.ErrInvalidPort):
		return "端口参数无效"
	case errors.Is(err, service.ErrInvalidAction):
		return "无效的操作类型，支持: add, del, list"
	case errors.Is(err, service.ErrPortConflict):
		return "服务端口或本地端口已被占用"
	default:
		return "服务器内部错误"
	}
}
