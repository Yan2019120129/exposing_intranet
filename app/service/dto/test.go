package dto

import "strings"

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
