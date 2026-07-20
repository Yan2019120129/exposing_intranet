package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetSymbolUsesOwnerOnlyPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "client.key")
	t.Setenv("EXPOSING_INTRANET_SYMBOL_PATH", path)

	(&Config{}).SetSymbol("client-symbol")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat client key: %v", err)
	}
	if permissions := info.Mode().Perm(); permissions != 0600 {
		t.Fatalf("client key permissions = %#o, want 0600", permissions)
	}
}
