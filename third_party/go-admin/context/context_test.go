package context

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

// TestPostForm 保证原有表单 API 保持兼容。
func TestPostForm(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("name=test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	values := NewContext(req).PostForm()
	if values.Get("name") != "test" {
		t.Fatalf("普通表单字段错误：%q", values.Get("name"))
	}
}

// TestPostFormInterruptedMultipart 验证中断的 multipart 请求会返回解析错误。
func TestPostFormInterruptedMultipart(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("__go_admin_t_", "test-token"); err != nil {
		t.Fatalf("写入测试 token 失败：%v", err)
	}
	part, err := writer.CreateFormFile("storage_path", "large.tar.gz")
	if err != nil {
		t.Fatalf("创建测试文件字段失败：%v", err)
	}
	if _, err = part.Write([]byte("partial file content")); err != nil {
		t.Fatalf("写入测试文件内容失败：%v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	values, err := NewContext(req).PostFormWithError()
	if err == nil {
		t.Fatal("中断的 multipart 请求未返回错误")
	}
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("中断请求返回了非预期错误：%v", err)
	}
	if values != nil {
		t.Fatalf("解析失败时不应返回表单值：%v", values)
	}
}

func TestSlash(t *testing.T) {
	assert.Equal(t, "/abc", slash("/abc"))
	assert.Equal(t, "/", slash(""))
	assert.Equal(t, "/abc", slash("abc/"))
	assert.Equal(t, "/", slash("/"))
	assert.Equal(t, "/abc", slash("/abc/"))
	assert.Equal(t, "/", slash("//"))
}

func TestJoin(t *testing.T) {
	assert.Equal(t, "/abc/abc", join(slash("/abc"), slash("/abc")))
	assert.Equal(t, "/", join(slash("/"), slash("/")))
	assert.Equal(t, "/abc", join(slash("/"), slash("/abc")))
	assert.Equal(t, "/abc", join(slash("abc/"), slash("/")))
	assert.Equal(t, "/abc", join(slash("/abc/"), slash("/")))
}

func TestTree(t *testing.T) {
	tree := tree()
	tree.addPath(stringToArr("/adm"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	tree.addPath(stringToArr("/admi"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	tree.addPath(stringToArr("/admin"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	tree.addPath(stringToArr("/admin/menu/new"), "POST", []Handler{func(ctx *Context) { fmt.Println(1) }})
	tree.addPath(stringToArr("/admin/menu/new"), "GET", []Handler{func(ctx *Context) { fmt.Println(1) }})
	tree.addPath(stringToArr("/admin/info/:__prefix"), "GET", []Handler{
		func(ctx *Context) { fmt.Println("auth") },
		func(ctx *Context) { fmt.Println("init") },
		func(ctx *Context) { fmt.Println("info") },
	})
	tree.addPath(stringToArr("/admin/info/:__prefix/detail"), "GET", []Handler{
		func(ctx *Context) { fmt.Println("auth") },
		func(ctx *Context) { fmt.Println("detail") },
	})

	fmt.Println("/admin/menu/new", "GET")
	h := tree.findPath(stringToArr("/admin/menu/new"), "GET")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	fmt.Println("/admin/menu/new", "POST")
	h = tree.findPath(stringToArr("/admin/menu/new"), "POST")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	fmt.Println("/admin/me/new", "POST")
	h = tree.findPath(stringToArr("/admin/me/new"), "POST")
	assert.Equal(t, h == nil, true)
	printHandler(h)
	fmt.Println("/admin/info/user", "GET")
	h = tree.findPath(stringToArr("/admin/info/user"), "GET")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	fmt.Println("/admin/info/user/detail", "GET")
	h = tree.findPath(stringToArr("/admin/info/user/detail"), "GET")
	assert.Equal(t, h != nil, true)
	printHandler(h)
	fmt.Println("=========== printChildren ===========")
	tree.printChildren()
}

func printHandler(h []Handler) {
	for _, value := range h {
		value(&Context{})
	}
}
