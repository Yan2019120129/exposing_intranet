package app

import (
	"bytes"
	"encoding/json"
	"my-base/app/models"
	appRouter "my-base/app/router"
	"my-base/code"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	adminConfig "github.com/GoAdminGroup/go-admin/modules/config"
	adminDB "github.com/GoAdminGroup/go-admin/modules/db"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestTestRoutesDoNotPolluteData(t *testing.T) {
	router, db := newIsolatedTestRouter(t)

	ping := requestJSON[apiResponse[string]](t, router, http.MethodGet, "/test", "", http.StatusOK)
	if ping.Data != "测试" {
		t.Fatalf("unexpected ping data: %q", ping.Data)
	}

	pageRecorder := performRequest(router, http.MethodGet, "/test-page", "")
	if pageRecorder.Code != http.StatusOK {
		t.Fatalf("expected test page status %d, got %d", http.StatusOK, pageRecorder.Code)
	}
	page := pageRecorder.Body.String()
	for _, expected := range []string{"Test CRUD", "test-form", "fetch('/tests'"} {
		if !strings.Contains(page, expected) {
			t.Fatalf("expected frontend page to contain %q, got %q", expected, page)
		}
	}

	initial := requestJSON[apiResponse[[]testDTO]](t, router, http.MethodGet, "/tests", "", http.StatusOK)
	if len(initial.Data) != 0 {
		t.Fatalf("expected isolated database to start empty, got %+v", initial.Data)
	}

	createdA := requestJSON[apiResponse[testDTO]](t, router, http.MethodPost, "/tests", `{"name":" alpha ","id":999,"createdAt":"2000-01-01T00:00:00Z"}`, http.StatusOK)
	if createdA.Data.Id == 0 || createdA.Data.Id == 999 || createdA.Data.Name != "alpha" {
		t.Fatalf("unexpected first create response: %+v", createdA.Data)
	}

	createdB := requestJSON[apiResponse[testDTO]](t, router, http.MethodPost, "/tests", `{"name":"beta"}`, http.StatusOK)
	if createdB.Data.Id == 0 || createdB.Data.Id == createdA.Data.Id || createdB.Data.Name != "beta" {
		t.Fatalf("unexpected second create response: %+v", createdB.Data)
	}

	list := requestJSON[apiResponse[[]testDTO]](t, router, http.MethodGet, "/tests", "", http.StatusOK)
	if len(list.Data) != 2 || list.Data[0].Name != "alpha" || list.Data[1].Name != "beta" {
		t.Fatalf("unexpected list after creates: %+v", list.Data)
	}

	filtered := requestJSON[apiResponse[[]testDTO]](t, router, http.MethodGet, "/tests?name=alpha", "", http.StatusOK)
	if len(filtered.Data) != 1 || filtered.Data[0].Id != createdA.Data.Id || filtered.Data[0].Name != "alpha" {
		t.Fatalf("unexpected filtered list: %+v", filtered.Data)
	}

	blankFilter := requestJSON[apiResponse[[]testDTO]](t, router, http.MethodGet, "/tests?name=%20%20", "", http.StatusOK)
	if len(blankFilter.Data) != 2 || blankFilter.Data[0].Id != createdA.Data.Id || blankFilter.Data[1].Id != createdB.Data.Id {
		t.Fatalf("unexpected blank-filter list: %+v", blankFilter.Data)
	}

	updatedA := requestJSON[apiResponse[testDTO]](t, router, http.MethodPut, "/tests/"+itoa(createdA.Data.Id), `{"name":"alpha-updated","id":`+itoa(createdB.Data.Id)+`}`, http.StatusOK)
	if updatedA.Data.Id != createdA.Data.Id || updatedA.Data.Name != "alpha-updated" {
		t.Fatalf("unexpected update response: %+v", updatedA.Data)
	}

	gotB := requestJSON[apiResponse[testDTO]](t, router, http.MethodGet, "/tests/"+itoa(createdB.Data.Id), "", http.StatusOK)
	if gotB.Data.Id != createdB.Data.Id || gotB.Data.Name != "beta" {
		t.Fatalf("update polluted second record: %+v", gotB.Data)
	}

	deleteRecorder := performRequest(router, http.MethodDelete, "/tests/"+itoa(createdA.Data.Id), "")
	if deleteRecorder.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d with body %q", http.StatusOK, deleteRecorder.Code, deleteRecorder.Body.String())
	}

	remaining := requestJSON[apiResponse[[]testDTO]](t, router, http.MethodGet, "/tests", "", http.StatusOK)
	if len(remaining.Data) != 1 || remaining.Data[0].Id != createdB.Data.Id || remaining.Data[0].Name != "beta" {
		t.Fatalf("delete polluted remaining records: %+v", remaining.Data)
	}

	var count int64
	if err := db.Model(&models.Test{}).Count(&count).Error; err != nil {
		t.Fatalf("count records: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected database to contain exactly one record, got %d", count)
	}

	var total int64
	if err := db.Unscoped().Model(&models.Test{}).Count(&total).Error; err != nil {
		t.Fatalf("count records including soft-deleted rows: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected two physical records after soft delete, got %d", total)
	}

	deleted := models.Test{}
	if err := db.Unscoped().First(&deleted, createdA.Data.Id).Error; err != nil {
		t.Fatalf("load soft-deleted record: %v", err)
	}
	if !deleted.DeletedAt.Valid {
		t.Fatal("expected deleted record to have a deletion timestamp")
	}
}

func newIsolatedTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../html/*")

	connection := adminDB.GetSqliteDB().InitDB(map[string]adminConfig.Database{
		code.DefaultGoAdminConnectionName: {
			Driver: adminDB.DriverSqlite,
			File:   filepath.Join(t.TempDir(), "test.db"),
		},
	})
	t.Cleanup(func() { _ = connection.Close() })

	db, err := connection.GetGorm(code.DefaultGoAdminConnectionName)
	if err != nil {
		t.Fatalf("get isolated sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Test{}); err != nil {
		t.Fatalf("migrate isolated test table: %v", err)
	}

	router.Use(func(c *gin.Context) {
		c.Set(code.ContextDBKey, connection)
	})
	appRouter.InitRouter(router)
	return router, db
}

type apiResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type testDTO struct {
	Id        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Name      string     `json:"name"`
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
	if _, ok := any(result).(apiResponse[testDTO]); ok {
		responseBody := recorder.Body.String()
		if !strings.Contains(responseBody, `"id":`) || strings.Contains(responseBody, `"ID":`) {
			t.Fatalf("test response must preserve lower-camel JSON fields, got %q", responseBody)
		}
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

func itoa(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
