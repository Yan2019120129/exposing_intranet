package utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Http struct {
	header     http.Header
	respHeader http.Header
	client     http.Client
	params     url.Values
}

func NewHttp() *Http {
	return &Http{header: make(http.Header),
		client: *http.DefaultClient,
	}
}

// Set 设置请求头信息
func (h *Http) Set(key, val string) *Http {
	h.header.Set(key, val)
	return h
}

// SetTransport 设置连接信息
func (h *Http) SetTransport(t *http.Transport) {
	h.client.Transport = t
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
func (h *Http) ask(method, url string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	request.Header = h.header
	do, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	h.respHeader = do.Header
	v, e := h.result(do.Body)
	h = NewHttp()
	return v, e
}

// Get Get请求
func (h *Http) Get(path string) ([]byte, error) {
	if h.params != nil {
		path += "?" + h.params.Encode()
	}
	return h.ask("GET", path, nil)
}

// Post Post请求
func (h *Http) Post(url string, s string) ([]byte, error) {
	return h.PostFormat(url, "application/json", s)
}

// Post Form-data
func (h *Http) PostFormData(url string) ([]byte, error) {
	return h.PostFormat(url, "application/json", "")
}

// PostFormat Post请求,带格式
func (h *Http) PostFormat(url string, ctxType string, s string) ([]byte, error) {
	h.Set("Content-Type", ctxType)
	var params io.Reader
	if s != "" {
		params = strings.NewReader(s)
	}
	if h.params != nil {
		url += "?" + h.params.Encode()
	}
	return h.ask("POST", url, params)
}

// PostFile 发起 multipart/form-data 文件上传请求
func (h *Http) PostFile(targetURL, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("upload", file.Name())
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	h.Set("Content-Type", writer.FormDataContentType())

	return h.ask("POST", targetURL, &requestBody)
}

// Put 发起 PUT 请求，默认 application/json
func (h *Http) Put(url string, s string) ([]byte, error) {
	return h.PutFormat(url, "application/json", s)
}

// PutFormat 发起 PUT 请求，带指定 Content-Type
func (h *Http) PutFormat(url string, ctxType string, s string) ([]byte, error) {
	h.Set("Content-Type", ctxType)
	var params io.Reader
	if s != "" {
		params = strings.NewReader(s)
	}
	return h.ask("PUT", url, params)
}

// Delete 发起 DELETE 请求
func (h *Http) Delete(url string) ([]byte, error) {
	if h.params != nil {
		url += "?" + h.params.Encode()
	}
	return h.ask("DELETE", url, nil)
}

// DeleteWithBody 发起带请求体的 DELETE 请求（部分接口需要）
func (h *Http) DeleteWithBody(url string, s string) ([]byte, error) {
	var params io.Reader
	if s != "" {
		params = strings.NewReader(s)
	}
	h.Set("Content-Type", "application/json")
	return h.ask("DELETE", url, params)
}

func (h *Http) result(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
