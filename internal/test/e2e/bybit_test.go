// +build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/pkg/bybit"
)

func TestClient_WalletBalance(t *testing.T) {
	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	res, err := cli.WalletBalance(ctx, &bybit.WalletBalanceReq{})
	if err != nil {
		t.Fatal(err)
	}
	forBreakPoint := res
	_ = forBreakPoint
}

func TestClient_PositionList(t *testing.T) {
	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	res, err := cli.PositionList(ctx, &bybit.PositionListReq{
		Symbol: bybit.SymbolBTCUSD,
	})
	if err != nil {
		t.Fatal(err)
	}
	forBreakPoint := res
	_ = forBreakPoint
}
