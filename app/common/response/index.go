package response

/**
 * 0 表示成功 其他表示失败
 * 0 means success, others means fail
 */

// Resp 同一返回数据结构体
type Resp struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

// Set 设置返回信息
func (r *Resp) Set(status int, v any, m string) {
	r.Msg = m
	r.Data = v
	r.Code = status
}
