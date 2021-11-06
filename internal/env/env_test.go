package env

import (
	"context"
	"testing"
)

func TestEnv(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	cfg.BybitKey = "test"
	ctx = With(ctx, cfg)

	if cfg := MustValue(ctx); cfg.BybitKey != "test" {
		t.Fatalf("want:test, got:%s", cfg.BybitKey)
	}
}
