package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Http Http
type Http struct {
	client      *http.Client      //	Client对象
	headers     map[string]string //	头信息
	contentType string            // 	参数格式
	baseURL     string            //	基础URL
}

// NewHttp 创建Http对象
func NewHttp(baseURL string) *Http {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &Http{
		client:      client,
		headers:     make(map[string]string),
		contentType: "application/json",
		baseURL:     baseURL,
	}
}

// GET Get请求
func (_Http *Http) GET(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", _Http.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	// 设置头信息
	for k, v := range _Http.headers {
		req.Header.Set(k, v)
	}

	resp, err := _Http.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// POST Post请求
func (_Http *Http) POST(url string, paramsBytes []byte) ([]byte, error) {
	resp, err := _Http.client.Post(_Http.baseURL+url, _Http.contentType, bytes.NewBuffer(paramsBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// HeaderContentJson 内容类型Json
func (_Http *Http) HeaderContentJson() *Http {
	_Http.contentType = "application/json"
	return _Http
}

// SetHeaderContentType 默认内容类型
func (_Http *Http) SetHeaderContentType(contentType string) *Http {
	_Http.contentType = contentType
	return _Http
}

// SetHeaders 设置头信息
func (_Http *Http) SetHeaders(headers map[string]string) *Http {
	for k, v := range headers {
		_Http.headers[k] = v
	}
	return _Http
}

// GetHttpHost 获取http服务器地址
func GetHttpHost(rawUrl string) string {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	return parsedURL.Scheme + "://" + parsedURL.Host
}

// IsPortEq 判断端口号是是否相等
func IsPortEq(addr, targetAddr string) bool {
	portIndex := strings.Index(addr, ":")
	targetPortIndex := strings.Index(targetAddr, ":")
	if portIndex != -1 && targetPortIndex != -1 && addr[portIndex:] == targetAddr[targetPortIndex:] {
		return true
	}
	return false
}
