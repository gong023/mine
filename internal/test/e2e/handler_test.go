// +build e2e
// +build post

package e2e

import (
	"context"
	"testing"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/internal/handler"
	"github.com/gong023/mine/pkg/bybit"
)

func TestHandler_Start(t *testing.T) {
	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}

	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	h := handler.New(cfg, cli)

	ctx := context.Background()
	webhooks := []string{
		"",
	}
	for _, w := range webhooks {
		decision, err := h.Start(ctx, []byte(w))
		if err != nil {
			t.Fatal(err)
		}
		forBreakPoint := decision
		_ = forBreakPoint
	}
}
