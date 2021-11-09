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
