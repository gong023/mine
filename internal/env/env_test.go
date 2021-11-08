package env

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err == nil {
		t.Fatal("must be error due to missing required fields")
	}
}

func TestWithCtx(t *testing.T) {
	cfg := Config{BybitKey: "test"}
	ctx := context.Background()
	ctx = With(ctx, cfg)

	if cfg := MustValue(ctx); cfg.BybitKey != "test" {
		t.Fatalf("want:test, got:%s", cfg.BybitKey)
	}
}
