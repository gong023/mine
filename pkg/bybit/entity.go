package bybit

type (
	Request interface {
		Path() string
	}

	GetRequest interface {
		Request
		IsGet() bool
	}

	PostRequest interface {
		Request
		IsPost() bool
	}

	Response struct {
		ExtCode          string  `json:"ext_code"`
		ExtInfo          string  `json:"ext_info"`
		RateLimit        int64   `json:"rate_limit"`
		RateLimitResetMs int64   `json:"rate_limit_reset_ms"`
		RateLimitStatus  int64   `json:"rate_limit_status"`
		RetCode          float32 `json:"ret_code"`
		RetMsg           string  `json:"ret_msg"`
		TimeNow          string  `json:"time_now"`
	}

	WalletBalanceRes struct {
		Response
		Result struct {
			BTC  WalletBalance `json:"BTC"`
			USDT WalletBalance `json:"USDT"`
			EOS  WalletBalance `json:"EOS"`
			XRP  WalletBalance `json:"XRP"`
			DOT  WalletBalance `json:"DOT"`
		} `json:"result"`
	}

	WalletBalance struct {
		Equity           float64 `json:"equity,omitempty"`
		AvailableBalance float64 `json:"available_balance,omitempty"`
		UsedMargin       float64 `json:"used_margin,omitempty"`
		OrderMargin      float64 `json:"order_margin,omitempty"`
		PositionMargin   float64 `json:"position_margin,omitempty"`
		OccClosingFee    float64 `json:"occ_closing_fee,omitempty"`
		OccFundingFee    float64 `json:"occ_funding_fee,omitempty"`
		WalletBalance    float64 `json:"wallet_balance,omitempty"`
		RealisedPnl      float64 `json:"realised_pnl,omitempty"`
		UnrealisedPnl    float64 `json:"unrealised_pnl,omitempty"`
		CumRealisedPnl   float64 `json:"cum_realised_pnl,omitempty"`
		GivenCash        float64 `json:"given_cash,omitempty"`
		ServiceCash      float64 `json:"service_cash,omitempty"`
	}

	OrderCreateRes struct {
		Response
		Result struct {
			OrderID       string  `json:"order_id,omitempty"`
			UserID        float32 `json:"user_id,omitempty"`
			Symbol        string  `json:"symbol,omitempty"`
			Side          string  `json:"side,omitempty"`
			OrderType     string  `json:"order_type,omitempty"`
			Price         float64 `json:"price,omitempty"`
			Qty           string  `json:"qty,omitempty"`
			TimeInForce   string  `json:"time_in_force,omitempty"`
			OrderStatus   string  `json:"order_status,omitempty"`
			LastExecTime  float64 `json:"last_exec_time,omitempty"`
			LastExecPrice float64 `json:"last_exec_price,omitempty"`
			LeavesQty     float32 `json:"leaves_qty,omitempty"`
			CumExecQty    float32 `json:"cum_exec_qty,omitempty"`
			CumExecValue  float32 `json:"cum_exec_value,omitempty"`
			CumExecFee    float64 `json:"cum_exec_fee,omitempty"`
			RejectReason  string  `json:"reject_reason,omitempty"`
			OrderLinkID   string  `json:"order_link_id,omitempty"`
			CreatedAt     string  `json:"created_at,omitempty"`
			UpdatedAt     string  `json:"updated_at,omitempty"`
		} `json:"result"`
	}
)

type WalletBalanceReq struct {
	Coin string `json:"coin,omitempty"`
}

func (w *WalletBalanceReq) Path() string {
	return "/v2/private/wallet/balance"
}

func (w *WalletBalanceReq) IsGet() bool {
	return true
}

type OrderCreateReq struct {
	Symbol         string  `json:"symbol"`
	Side           string  `json:"side"`
	OrderType      string  `json:"order_type"`
	TimeInForce    string  `json:"time_in_force"`
	Qty            float64 `json:"qty"`
	Price          float64 `json:"price,omitempty"`
	TakeProfit     float64 `json:"take_profit,omitempty"`
	StopLoss       float64 `json:"stop_loss,omitempty"`
	ReduceOnly     bool    `json:"reduce_only,omitempty"`
	TpTriggerBy    string  `json:"tp_trigger_by,omitempty"`
	SlTriggerBy    string  `json:"sl_trigger_by,omitempty"`
	CloseOnTrigger bool    `json:"close_on_trigger,omitempty"`
	OrderLinkID    string  `json:"order_link_id,omitempty"`
}

func (o *OrderCreateReq) Path() string {
	return "/v2/private/order/create"
}

func (o *OrderCreateReq) IsPost() bool {
	return true
}

type OrderCancelReq struct {
	Symbol      string `json:"symbol"`
	OrderID     string `json:"order_id,omitempty"`
	OrderLinkID string `json:"order_link_id,omitempty"`
}

func (o *OrderCancelReq) Path() string {
	return "/v2/private/order/cancel"
}

func (o *OrderCancelReq) IsPost() bool {
	return true
}
