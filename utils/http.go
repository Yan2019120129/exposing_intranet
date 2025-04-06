package utils

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Http struct {
	header     http.Header
	respHeader http.Header
	params     url.Values
}

func NewHttp() *Http {
	return &Http{header: make(http.Header)}
}

// Set 设置请求头信息
func (h *Http) Set(key, val string) *Http {
	h.header.Set(key, val)
	return h
}

// AddParam 设置get参数
func (h *Http) AddParam(key, val string) *Http {
	if h.params == nil {
		h.params = make(url.Values)
	}
	h.params.Set(key, val)
	return h
}

// GetParam 获取get参数
func (h *Http) GetParam() string {
	return h.params.Encode()
}

// GetRHeader 获取http相应头
func (h *Http) GetRHeader() http.Header {
	return h.respHeader
}

// GetHeader 获取http请求头
func (h *Http) GetHeader() http.Header {
	return h.header
}

// ask 发起请求
func (h *Http) ask(method, url string, body io.Reader) []byte {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	h.header = request.Header
	do, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	h.respHeader = do.Header
	return h.result(do.Body)
}

// Get Get请求
func (h *Http) Get(path string) []byte {
	if h.params != nil {
		path += "?" + h.params.Encode()
	}
	return h.ask("GET", path, nil)
}

// Post Post请求
func (h *Http) Post(url string, ctxType string, s string) []byte {
	h.Set("Content-Type", ctxType)
	var params io.Reader
	if s != "" {
		params = strings.NewReader(s)
	}
	return h.ask("POST", url, params)
}

func (h *Http) result(body io.ReadCloser) []byte {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return data
}
