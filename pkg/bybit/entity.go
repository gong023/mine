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
		ExtCode          string `json:"ext_code"`
		ExtInfo          string `json:"ext_info"`
		RateLimit        int64  `json:"rate_limit"`
		RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
		RateLimitStatus  int64  `json:"rate_limit_status"`
		RetCode          int64  `json:"ret_code"`
		RetMsg           string `json:"ret_msg"`
		TimeNow          string `json:"time_now"`
	}

	WalletBalanceRes struct {
		Response
		Result struct {
			BTC struct {
				AvailableBalance float64 `json:"available_balance"`
				CumRealisedPnl   int64   `json:"cum_realised_pnl"`
				Equity           int64   `json:"equity"`
				GivenCash        int64   `json:"given_cash"`
				OccClosingFee    int64   `json:"occ_closing_fee"`
				OccFundingFee    int64   `json:"occ_funding_fee"`
				OrderMargin      float64 `json:"order_margin"`
				PositionMargin   int64   `json:"position_margin"`
				RealisedPnl      int64   `json:"realised_pnl"`
				ServiceCash      int64   `json:"service_cash"`
				UnrealisedPnl    int64   `json:"unrealised_pnl"`
				UsedMargin       float64 `json:"used_margin"`
				WalletBalance    int64   `json:"wallet_balance"`
			} `json:"BTC"`
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
