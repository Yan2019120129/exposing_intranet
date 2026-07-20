package utils

import (
	"os/exec"
	"testing"
)

func TestExecReturnsStandardOutput(t *testing.T) {
	output, err := Exec(exec.Command("/bin/sh", "-c", "printf output; printf warning >&2"))
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if output != "output" {
		t.Fatalf("Exec() output = %q, want %q", output, "output")
	}
}
