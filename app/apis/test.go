package apis

import (
	"errors"
	"my-base/app/models"
	"my-base/app/service"
	"my-base/app/service/dto"
	"my-base/code/api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Test struct {
	api.Api
}

func (e Test) Ping(ctx *gin.Context) {
	e.MakeContext(ctx)
	e.OK("测试", "查询成功")
}

func (e Test) Page(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "test.tmpl", nil)
}

// List 根据查询条件获取测试记录列表。
func (e Test) List(ctx *gin.Context) {
	s := service.Test{}
	err := e.MakeContext(ctx).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}

	items := make([]models.Test, 0)
	query := dto.TestListQuery{}
	if err = ctx.ShouldBindQuery(&query); err != nil {
		e.Error(http.StatusBadRequest, err, "invalid query parameters")
		return
	}
	if err = s.List(&query, &items); err != nil {
		e.Error(http.StatusInternalServerError, err, err.Error())
		return
	}
	e.OK(dto.NewTestItems(items), "查询成功")
}

func (e Test) Create(ctx *gin.Context) {
	payload := dto.TestPayload{}
	s := service.Test{}
	err := e.MakeContext(ctx).
		MakeOrm().
		Bind(&payload).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	item := models.Test{}
	if err := s.Create(&payload, &item); err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(dto.NewTestItem(item), "创建成功")
}

func (e Test) Get(ctx *gin.Context) {
	s := service.Test{}
	err := e.MakeContext(ctx).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	id, ok := e.parseID(ctx)
	if !ok {
		return
	}

	item := models.Test{}
	if err := s.Get(id, &item); err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(dto.NewTestItem(item), "查询成功")
}

func (e Test) Update(ctx *gin.Context) {
	s := service.Test{}
	err := e.MakeContext(ctx).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	id, ok := e.parseID(ctx)
	if !ok {
		return
	}

	payload := dto.TestPayload{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e.Error(http.StatusBadRequest, err, "invalid json body")
		return
	}

	item := models.Test{}
	if err := s.Update(id, &payload, &item); err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(dto.NewTestItem(item), "修改成功")
}

func (e Test) Delete(ctx *gin.Context) {
	s := service.Test{}
	err := e.MakeContext(ctx).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	id, ok := e.parseID(ctx)
	if !ok {
		return
	}

	if err := s.Delete(id); err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(nil, "删除成功")
}

// parseID 解析路径中的无符号测试记录标识。
func (e Test) parseID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, strconv.IntSize)
	if err != nil || id <= 0 {
		e.Error(http.StatusBadRequest, errors.New("invalid id"), "invalid id")
		return 0, false
	}
	return uint(id), true
}
