package apis

import (
	"errors"
	"net/http"

	"my-base/app/service"
	"my-base/code/api"
	"my-base/code/contract"

	"github.com/gin-gonic/gin"
)

type Client struct {
	api.Api
}

type RegisterClientRequest = contract.AuthRequest
type RegisterClientResponse = contract.AuthResponse

func (e Client) Register(ctx *gin.Context) {
	e.MakeContext(ctx)

	var req RegisterClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, RegisterClientResponse{Success: false, Message: "参数错误"})
		return
	}
	if err := e.MakeOrm().Errors; err != nil {
		ctx.JSON(http.StatusInternalServerError, RegisterClientResponse{Success: false, Message: "服务器内部错误"})
		return
	}

	symbol, err := service.NewClientService(e.Orm).Register(service.RegisterClientInput{
		Username: req.Username,
		Password: req.Password,
		Hostname: req.Hostname,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, RegisterClientResponse{Success: false, Message: "用户名或密码错误"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, RegisterClientResponse{Success: false, Message: "服务器内部错误"})
		return
	}

	ctx.JSON(http.StatusOK, RegisterClientResponse{Success: true, Symbol: symbol, Message: "认证成功"})
}
