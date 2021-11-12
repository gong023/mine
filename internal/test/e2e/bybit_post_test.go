// +build e2e
// +build post

package e2e

import (
	"context"
	"testing"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/pkg/bybit"
)

func TestClient_OrderCreate(t *testing.T) {
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
		//{
		//	Symbol:      bybit.SymbolBTCUSD,
		//	Side:        bybit.SideBuy,
		//	OrderType:   bybit.OrderTypeMarket,
		//	TimeInForce: bybit.TIFGoodTillCancel,
		//	Qty:         30,
		//},
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

func TestClient_PositionLeverageSave(t *testing.T) {
	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	cli := bybit.NewClient(cfg.BybitHost, cfg.BybitKey, cfg.BybitSec)
	res, err := cli.PositionLeverageSave(ctx, &bybit.PositionLeverageSaveReq{
		Symbol:       bybit.SymbolBTCUSD,
		Leverage:     20,
		LeverageOnly: true,
	})
	if err != nil {
		re, ok := err.(*bybit.ResponseError)
		if !ok || re.Response.RetCode != bybit.RetCodeLeverageNotModified {
			t.Fatal(err)
		}
	}
	forBreakPoint := res
	_ = forBreakPoint
}
