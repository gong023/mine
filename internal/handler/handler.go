package handler

import (
	"context"

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
	position, err := h.client.PositionList(ctx, &bybit.PositionListReq{Symbol: bybit.SymbolBTCUSD})
	if err != nil {
		return nil, err
	}

	balance, err := h.client.WalletBalance(ctx, &bybit.WalletBalanceReq{})
	if err != nil {
		return nil, err
	}

	decision := &Decision{
		config:   h.config,
		webhook:  webhook,
		position: position.Result,
	}
	conclusion := decision.Make()
	orders := h.NewOrders(conclusion, &position.Result, &balance.Result.USDT)
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

func (h *Handler) NewOrders(
	conclusion Conclusion,
	position *bybit.Position,
	balance *bybit.WalletBalance,
) []*bybit.OrderCreateReq {
	return nil
}
