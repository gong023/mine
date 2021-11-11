// +build e2e

package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/pkg/bybit"
)

func TestGetWalletBalance(t *testing.T) {
	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	_, err = cli.GetWalletBalance(ctx, &bybit.WalletBalanceReq{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrderCreate(t *testing.T) {
	if r := os.Getenv("TEST_POST"); r == "" {
		t.Skip(t.Name())
	}

	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}

	orders := []*bybit.OrderCreateReq{
		{
			Symbol:      bybit.SymbolBTCUSD,
			Side:        bybit.SideBuy,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         30,
			StopLoss:    60000,
		},
		{
			Symbol:      bybit.SymbolBTCUSD,
			Side:        bybit.SideSell,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         30,
		},
		{
			Symbol:      bybit.SymbolBTCUSD,
			Side:        bybit.SideSell,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         30,
			StopLoss:    70000,
		},
		{
			Symbol:      bybit.SymbolBTCUSD,
			Side:        bybit.SideBuy,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         30,
		},
	}

	ctx := context.Background()
	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	for i, order := range orders {
		res, err := cli.OrderCreate(ctx, order)
		if err != nil {
			t.Fatalf("failed at:%d, %s", i, err)
		}
		forBreakPoint := res
		_ = forBreakPoint
	}
}
