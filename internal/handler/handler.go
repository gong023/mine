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

	decision := &Decision{webhook: webhook, position: &position.Result}
	conclusion := decision.Make()
	orders := h.NewOrders(conclusion, &balance.Result.USDT)
	for _, order := range orders {
		if _, err := h.client.OrderCreate(ctx, order); err != nil {
			return decision, err
		}
	}
	return decision, nil
}

func (h *Handler) NewOrders(conclusion Conclusion, balance *bybit.WalletBalance) []*bybit.OrderCreateReq {
	return nil
}
