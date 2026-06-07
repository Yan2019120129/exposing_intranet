package app

import (
	"fmt"
	"my-base/app/models"
	appRouter "my-base/app/router"
	"my-base/configs"
	"my-base/tables"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	"github.com/GoAdminGroup/go-admin/engine"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"
	_ "github.com/GoAdminGroup/themes/sword"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestTestRoutesRealDatabaseConcurrentIsolation(t *testing.T) {
	if os.Getenv("MY_BASE_REAL_DB_TEST") != "1" {
		t.Skip("set MY_BASE_REAL_DB_TEST=1 to run real database concurrency test")
	}

	router, db, closeDB := newRealDatabaseTestRouter(t)
	t.Cleanup(closeDB)

	if err := db.AutoMigrate(&models.Test{}); err != nil {
		t.Fatalf("migrate real test table: %v", err)
	}

	prefix := fmt.Sprintf("codex_concurrent_%d_", time.Now().UnixNano())
	t.Cleanup(func() {
		if err := db.Unscoped().Where("name LIKE ?", prefix+"%").Delete(&models.Test{}).Error; err != nil {
			t.Logf("cleanup real database test data: %v", err)
		}
	})

	const workers = 12
	created := make([]testDTO, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			name := fmt.Sprintf("%screate_%02d", prefix, i)
			body := fmt.Sprintf(`{"name":%q,"id":%d}`, name, 100000+i)
			resp := requestJSON[apiResponse[testDTO]](t, router, http.MethodPost, "/tests", body, http.StatusOK)
			if resp.Code != http.StatusOK {
				t.Errorf("create worker %d business code=%d msg=%q data=%+v", i, resp.Code, resp.Msg, resp.Data)
				return
			}
			if resp.Data.Id == 0 || resp.Data.Id == 100000+i || resp.Data.Name != name {
				t.Errorf("unexpected concurrent create response for worker %d: %+v", i, resp.Data)
				return
			}
			created[i] = resp.Data
		}()
	}
	wg.Wait()

	for i, item := range created {
		if item.Id == 0 {
			t.Fatalf("worker %d did not create a record", i)
		}
	}

	var count int64
	if err := db.Model(&models.Test{}).Where("name LIKE ?", prefix+"%").Count(&count).Error; err != nil {
		t.Fatalf("count concurrent creates: %v", err)
	}
	if count != workers {
		t.Fatalf("expected %d created records for prefix %q, got %d", workers, prefix, count)
	}

	for i := 0; i < workers; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			newName := fmt.Sprintf("%supdated_%02d", prefix, i)
			body := fmt.Sprintf(`{"name":%q,"id":%d}`, newName, created[(i+1)%workers].Id)
			resp := requestJSON[apiResponse[testDTO]](t, router, http.MethodPut, "/tests/"+itoa(created[i].Id), body, http.StatusOK)
			if resp.Code != http.StatusOK {
				t.Errorf("update worker %d business code=%d msg=%q data=%+v", i, resp.Code, resp.Msg, resp.Data)
				return
			}
			if resp.Data.Id != created[i].Id || resp.Data.Name != newName {
				t.Errorf("unexpected concurrent update response for worker %d: %+v", i, resp.Data)
			}
		}()
	}
	wg.Wait()

	for i, item := range created {
		got := requestJSON[apiResponse[testDTO]](t, router, http.MethodGet, "/tests/"+itoa(item.Id), "", http.StatusOK)
		if got.Code != http.StatusOK {
			t.Fatalf("get record %d business code=%d msg=%q data=%+v", i, got.Code, got.Msg, got.Data)
		}
		wantName := fmt.Sprintf("%supdated_%02d", prefix, i)
		if got.Data.Id != item.Id || got.Data.Name != wantName {
			t.Fatalf("record %d was polluted after concurrent updates: got %+v want name %q", i, got.Data, wantName)
		}
	}

	for i := 0; i < workers; i += 2 {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			recorder := performRequest(router, http.MethodDelete, "/tests/"+itoa(created[i].Id), "")
			if recorder.Code != http.StatusOK {
				t.Errorf("delete worker %d status=%d body=%q", i, recorder.Code, recorder.Body.String())
			}
		}()
	}
	wg.Wait()

	if err := db.Model(&models.Test{}).Where("name LIKE ?", prefix+"%").Count(&count).Error; err != nil {
		t.Fatalf("count after concurrent deletes: %v", err)
	}
	if count != workers/2 {
		t.Fatalf("expected %d remaining records after deletes, got %d", workers/2, count)
	}

	for i := 1; i < workers; i += 2 {
		got := requestJSON[apiResponse[testDTO]](t, router, http.MethodGet, "/tests/"+itoa(created[i].Id), "", http.StatusOK)
		if got.Code != http.StatusOK {
			t.Fatalf("get remaining record %d business code=%d msg=%q data=%+v", i, got.Code, got.Msg, got.Data)
		}
		wantName := fmt.Sprintf("%supdated_%02d", prefix, i)
		if got.Data.Id != created[i].Id || got.Data.Name != wantName {
			t.Fatalf("remaining record %d was polluted: got %+v want name %q", i, got.Data, wantName)
		}
	}
}

func newRealDatabaseTestRouter(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../html/*")

	eng := engine.Default()
	if err := eng.AddConfig(configs.GetAdmin()).
		AddGenerators(tables.Generators).
		Use(router); err != nil {
		t.Fatalf("initialize go-admin engine: %v", err)
	}

	db, err := eng.DefaultConnection().GetGorm("default")
	if err != nil {
		eng.DefaultConnection().Close()
		t.Fatalf("get engine gorm: %v", err)
	}

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	appRouter.InitRouter(router)

	return router, db, func() { eng.DefaultConnection().Close() }
}
