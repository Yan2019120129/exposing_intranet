package dto

import (
	"strings"
	"time"

	"my-base/app/models"
)

// TestPayload 定义创建和更新测试记录的请求体。
type TestPayload struct {
	Name string `json:"name"`
}

// TestListQuery 定义测试记录列表的查询条件。
type TestListQuery struct {
	Name string `form:"name"`
}

// Normalize 规范化查询条件中的文本字段。
func (q *TestListQuery) Normalize() {
	q.Name = strings.TrimSpace(q.Name)
}

// TestItem 定义测试记录的接口响应，保持既有 JSON 字段格式。
type TestItem struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Name      string     `json:"name"`
}

// NewTestItem 将测试模型转换为对外响应对象。
func NewTestItem(item models.Test) TestItem {
	result := TestItem{
		Id:        item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		Name:      item.Name,
	}
	if item.DeletedAt.Valid {
		deletedAt := item.DeletedAt.Time
		result.DeletedAt = &deletedAt
	}
	return result
}

// NewTestItems 将测试模型列表转换为对外响应列表。
func NewTestItems(items []models.Test) []TestItem {
	result := make([]TestItem, 0, len(items))
	for _, item := range items {
		result = append(result, NewTestItem(item))
	}
	return result
}
