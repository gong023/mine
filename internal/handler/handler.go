package handler

import (
	"context"
	"errors"
	"math"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/pkg/bybit"
)

type Handler struct {
	config env.Config
	client bybit.ClientType
}

func New(config env.Config, client bybit.ClientType) *Handler {
	return &Handler{config: config, client: client}
}

func (h *Handler) Start(ctx context.Context, webhook []byte) (*Decision, error) {
	position, err := h.client.PositionList(ctx, &bybit.PositionListReq{
		Symbol: h.config.Symbol,
	})
	if err != nil {
		return nil, err
	}

	balance, err := h.client.WalletBalance(ctx, &bybit.WalletBalanceReq{})
	if err != nil {
		return nil, err
	}
	if balance.Result.USDT.AvailableBalance == 0 {
		return nil, errors.New("no available balance for USDT")
	}

	decision := &Decision{
		config:   h.config,
		webhook:  webhook,
		position: position.Result,
	}
	conclusion := decision.Make()
	orders := h.NewOrders(conclusion, position.Result, balance.Result.USDT)
	for _, order := range orders {
		if _, err := h.client.OrderCreate(ctx, order); err != nil {
			return decision, err
		}
	}

	if conclusion != DoNothing {
		_, err := h.client.PositionLeverageSave(ctx, &bybit.PositionLeverageSaveReq{
			Symbol:       h.config.Symbol,
			Leverage:     h.config.Leverage,
			LeverageOnly: true,
		})
		if err != nil {
			r, ok := err.(*bybit.ResponseError)
			if !ok || r.RetCode != bybit.RetCodeLeverageNotModified {
				return decision, err
			}
		}
	}

	return decision, nil
}

// NewOrders issues the orders with isolated margin, and MarketPrice.
func (h *Handler) NewOrders(
	conclusion Conclusion,
	position bybit.Position,
	balance bybit.WalletBalance,
) (orders []*bybit.OrderCreateReq) {
	if conclusion == ReleaseThenLong {
		orders = append(orders, &bybit.OrderCreateReq{
			Symbol:      h.config.Symbol,
			Side:        bybit.SideBuy,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         position.Size,
		})
	}
	if conclusion == ReleaseThenShort {
		orders = append(orders, &bybit.OrderCreateReq{
			Symbol:      h.config.Symbol,
			Side:        bybit.SideSell,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         position.Size,
		})
	}

	qty := round(balance.AvailableBalance / h.config.UseBalance)
	if conclusion == Long || conclusion == ReleaseThenLong {
		orders = append(orders, &bybit.OrderCreateReq{
			Symbol:      h.config.Symbol,
			Side:        bybit.SideBuy,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         qty,
		})
	}
	if conclusion == Short || conclusion == ReleaseThenShort {
		orders = append(orders, &bybit.OrderCreateReq{
			Symbol:      h.config.Symbol,
			Side:        bybit.SideSell,
			OrderType:   bybit.OrderTypeMarket,
			TimeInForce: bybit.TIFGoodTillCancel,
			Qty:         qty,
		})
	}

	return orders
}

func round(src float64) (dest float64) {
	shift := math.Pow(10, 2)
	shifted := src * shift
	t := math.Trunc(shifted)
	if math.Abs(shifted-t) >= 0.5 {
		return (t + math.Copysign(1, shifted)) / shift
	}
	return t / shift
}
