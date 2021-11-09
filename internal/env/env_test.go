package env

import (
	"context"
	"testing"
)

func TestWithCtx(t *testing.T) {
	cfg := Config{BybitKey: "test"}
	ctx := context.Background()
	ctx = With(ctx, cfg)

	if cfg := MustValue(ctx); cfg.BybitKey != "test" {
		t.Fatalf("want:test, got:%s", cfg.BybitKey)
	}
}
