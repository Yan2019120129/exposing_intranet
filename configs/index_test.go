package configs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigReportsYAMLErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("port: ["), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := loadConfigFile(path)
	if err == nil || !strings.Contains(err.Error(), "parse") {
		t.Fatalf("loadConfigFile() error = %v, want parse error", err)
	}
}
