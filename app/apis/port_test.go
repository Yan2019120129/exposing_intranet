package apis

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"my-base/app/service"
)

func TestPortInvalidInputIsBadRequest(t *testing.T) {
	err := errors.Join(service.ErrInvalidPort, errors.New("invalid value"))
	if got := portErrorStatus(err); got != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", got, http.StatusBadRequest)
	}
	if got := portErrorMessage(err); got != "端口参数无效" {
		t.Fatalf("message = %q", got)
	}
}

func TestPortInternalErrorDoesNotLeakDetails(t *testing.T) {
	err := errors.New("Error 1062: duplicate SQL INSERT INTO port")
	if got := portErrorStatus(err); got != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", got, http.StatusInternalServerError)
	}
	message := portErrorMessage(err)
	if message != "服务器内部错误" || strings.Contains(message, "SQL") || strings.Contains(message, "1062") {
		t.Fatalf("unsafe internal error message: %q", message)
	}
}
