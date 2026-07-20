package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger(WithLevel(TraceLevel), WithName("test"))
	h1 := NewHelper(l).WithFields(map[string]interface{}{"key1": "val1"})
	h1.Trace("trace_msg1")
	h1.Warn("warn_msg1")

	h2 := NewHelper(l).WithFields(map[string]interface{}{"key2": "val2"})
	h2.Trace("trace_msg2")
	h2.Warn("warn_msg2")

	h3 := NewHelper(l).WithFields(map[string]interface{}{"key3": "val4"})
	h3.Info("test_msg")
	ctx := context.TODO()
	ctx = context.WithValue(ctx, &loggerKey{}, h3)
	v := ctx.Value(&loggerKey{})
	ll := v.(*Helper)
	ll.Info("test_msg")
}

func TestHelperFieldsDoNotMutateBaseLogger(t *testing.T) {
	var output bytes.Buffer
	base := NewLogger(WithLevel(InfoLevel), WithOutput(&output), WithTimeFormat(""))
	helper := NewHelper(base).WithFields(map[string]interface{}{"requestId": "request-a"})
	helper.Info("request message")
	base.Log(InfoLevel, "base message")

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d log lines, want 2: %q", len(lines), output.String())
	}
	if !strings.Contains(lines[0], "request-a") {
		t.Fatalf("request log missing its fields: %q", lines[0])
	}
	if strings.Contains(lines[1], "request-a") {
		t.Fatalf("base logger inherited request fields: %q", lines[1])
	}
}
