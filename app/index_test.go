package app

import (
	"bytes"
	"encoding/json"
	"my-base/app/models"
	appRouter "my-base/app/router"
	orm "my-base/module/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func TestFrontendBackendCRUDFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.LoadHTMLGlob("../html/*")
	appRouter.InitRouter(router)

	pageRecorder := performRequest(router, http.MethodGet, "/api/adm/test-page", "")
	if pageRecorder.Code != http.StatusOK {
		t.Fatalf("expected frontend status %d, got %d", http.StatusOK, pageRecorder.Code)
	}
	page := pageRecorder.Body.String()
	for _, expected := range []string{"Test CRUD", "test-form", "fetch('/tests'"} {
		if !strings.Contains(page, expected) {
			t.Fatalf("expected frontend page to contain %q, got %q", expected, page)
		}
	}

	prepareTestDatabase(t)

	created := requestJSON[testResponse](t, router, http.MethodPost, "/api/adm/tests", `{"name":"before"}`, http.StatusOK)
	if created.Data.Id == 0 || created.Data.Name != "before" {
		t.Fatalf("unexpected created test: %+v", created.Data)
	}

	got := requestJSON[testResponse](t, router, http.MethodGet, "/api/adm/tests/"+itoa(created.Data.Id), "", http.StatusOK)
	if got.Data.Name != "before" {
		t.Fatalf("expected detail name %q, got %q", "before", got.Data.Name)
	}

	updated := requestJSON[testResponse](t, router, http.MethodPut, "/api/adm/tests/"+itoa(created.Data.Id), `{"name":"after"}`, http.StatusOK)
	if updated.Data.Id != created.Data.Id || updated.Data.Name != "after" {
		t.Fatalf("unexpected updated test: %+v", updated.Data)
	}

	list := requestJSON[testListResponse](t, router, http.MethodGet, "/api/adm/tests", "", http.StatusOK)
	if len(list.Data) != 1 || list.Data[0].Name != "after" {
		t.Fatalf("unexpected list response: %+v", list.Data)
	}

	deleteRecorder := performRequest(router, http.MethodDelete, "/api/adm/tests/"+itoa(created.Data.Id), "")
	if deleteRecorder.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d", http.StatusOK, deleteRecorder.Code)
	}

	emptyList := requestJSON[testListResponse](t, router, http.MethodGet, "/api/adm/tests", "", http.StatusOK)
	if len(emptyList.Data) != 0 {
		t.Fatalf("expected empty list after delete, got %+v", emptyList.Data)
	}
}

type testResponse struct {
	Data testDTO `json:"data"`
}

type testListResponse struct {
	Data []testDTO `json:"data"`
}

type testDTO struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func prepareTestDatabase(t *testing.T) {
	t.Helper()

	if orm.DB == nil {
		t.Skip("database is not configured")
	}
	sqlDB, err := orm.DB.DB()
	if err != nil {
		t.Skip("database is not available")
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skip("database is not available")
	}
	if err := orm.DB.AutoMigrate(&models.Test{}); err != nil {
		t.Fatalf("auto migrate test table: %v", err)
	}
	if err := orm.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Where("1 = 1").Delete(&models.Test{}).Error; err != nil {
		t.Fatalf("clean test table: %v", err)
	}
}

func requestJSON[T any](t *testing.T, router *gin.Engine, method, path, body string, status int) T {
	t.Helper()

	recorder := performRequest(router, method, path, body)
	if recorder.Code != status {
		t.Fatalf("expected %s %s status %d, got %d with body %q", method, path, status, recorder.Code, recorder.Body.String())
	}

	var result T
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode %s %s response: %v", method, path, err)
	}
	return result
}

func performRequest(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(recorder, request)
	return recorder
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
