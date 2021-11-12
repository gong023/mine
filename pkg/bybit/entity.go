package bybit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	SymbolBTCUSDT = "BTCUSDT"
	SymbolBTCUSD  = "BTCUSD"
	SymbolUSDT    = "USDT"

	RetCodeSuccess             float32 = 0
	RetCodePramsError          float32 = 10001
	RetCodeLeverageNotModified float32 = 34036

	RetMsgOK = "OK"

	SideBuy  = "Buy"
	SideSell = "Sell"
	SideNone = "None"

	OrderTypeMarket = "Market"
	OrderTypeLimit  = "Limit"

	TIFGoodTillCancel    = "GoodTillCancel"
	TIFImmediateOrCancel = "ImmediateOrCancel"
	TIFFillOrKill        = "FillOrKill"
	TIFPostOnly          = "PostOnly"

	TriggerByLastPrice  = "LastPrice"
	TriggerByIndexPrice = "IndexPrice"
	TriggerByMarkPrice  = "MarkPrice"
)

type (
	RequestType interface {
		Path() string
	}

	ResponseType interface {
		GetCommon() Response
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
)

type ResponseError struct {
	Response
	RawResBody []byte
	Method     string
	URL        *url.URL
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("method:%s, url:%s, resBody:%s", r.Method, r.URL, r.RawResBody)
}

func NewResponseError(r Response, rr []byte, method string, url *url.URL) error {
	return &ResponseError{
		Response:   r,
		RawResBody: rr,
		Method:     method,
		URL:        url,
	}
}

type WalletBalanceReq struct {
	Coin string `json:"coin,omitempty"`
}

func (w *WalletBalanceReq) Path() string {
	return "/v2/private/wallet/balance"
}

type (
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

func (w *WalletBalanceRes) GetCommon() Response {
	return w.Response
}

type (
	OrderCreateReq struct {
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
)

func (o *OrderCreateReq) Path() string {
	return "/v2/private/order/create"
}

type (
	OrderCreateRes struct {
		Response
		Result Order `json:"result"`
	}

	Order struct {
		OrderID       string  `json:"order_id,omitempty"`
		UserID        float32 `json:"user_id,omitempty"`
		Symbol        string  `json:"symbol,omitempty"`
		Side          string  `json:"side,omitempty"`
		OrderType     string  `json:"order_type,omitempty"`
		Price         float64 `json:"price,omitempty"`
		Qty           float64 `json:"qty,omitempty"`
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
	}
)

func (o *OrderCreateRes) GetCommon() Response {
	return o.Response
}

type OrderCancelReq struct {
	Symbol      string `json:"symbol"`
	OrderID     string `json:"order_id,omitempty"`
	OrderLinkID string `json:"order_link_id,omitempty"`
}

func (o *OrderCancelReq) Path() string {
	return "/v2/private/order/cancel"
}

type OrderCancelRes struct {
	Response
	Result Order `json:"result"`
}

func (o *OrderCancelRes) GetCommon() Response {
	return o.Response
}

type PositionLeverageSaveReq struct {
	Symbol       string  `json:"symbol"`
	Leverage     float64 `json:"leverage"`
	LeverageOnly bool    `json:"leverage_only,omitempty"`
}

func (p *PositionLeverageSaveReq) Path() string {
	return "/v2/private/position/leverage/save"
}

type PositionLeverageSaveRes struct {
	Response
}

func (p *PositionLeverageSaveRes) GetCommon() Response {
	return p.Response
}

type PositionListReq struct {
	// Though the doc says this is optional, set this
	// otherwise the response structure is changed.
	Symbol string `json:"symbol"`
}

func (p *PositionListReq) Path() string {
	return "/v2/private/position/list"
}

type (
	PositionListRes struct {
		Response
		Result Position `json:"result"`
	}

	Position struct {
		ID                  float32     `json:"id,omitempty"`
		UserID              float32     `json:"user_id,omitempty"`
		RiskID              float32     `json:"risk_id,omitempty"`
		Symbol              string      `json:"symbol,omitempty"`
		Side                string      `json:"side,omitempty"`
		Size                float64     `json:"size,omitempty"`
		PositionValue       json.Number `json:"position_value,omitempty"`
		EntryPrice          json.Number `json:"entry_price,omitempty"`
		Leverage            json.Number `json:"leverage,omitempty"`
		AutoAddMargin       json.Number `json:"auto_add_margin,omitempty"`
		PositionMargin      json.Number `json:"position_margin,omitempty"`
		LiqPrice            json.Number `json:"liq_price,omitempty"`
		BustPrice           json.Number `json:"bust_price,omitempty"`
		OccClosingFee       json.Number `json:"occ_closing_fee,omitempty"`
		OccFundingFee       json.Number `json:"occ_funding_fee,omitempty"`
		TakeProfit          json.Number `json:"take_profit,omitempty"`
		StopLoss            json.Number `json:"stop_loss,omitempty"`
		PositionStatus      string      `json:"position_status,omitempty"`
		DeleverageIndicator json.Number `json:"deleverage_indicator,omitempty"`
		OcCalcData          string      `json:"oc_calc_data,omitempty"`
		OrderMargin         json.Number `json:"order_margin,omitempty"`
		WalletBalance       json.Number `json:"wallet_balance,omitempty"`
		UnrealisedPNL       json.Number `json:"unrealised_pnl,omitempty"`
		RealisedPNL         json.Number `json:"realised_pnl,omitempty"`
		CumRealisedPNL      json.Number `json:"cum_realised_pnl,omitempty"`
		CumCommission       json.Number `json:"cum_commission,omitempty"`
		CrossSeq            json.Number `json:"cross_seq,omitempty"`
		PositionSeq         json.Number `json:"position_seq,omitempty"`
		CreatedAt           time.Time   `json:"created_at,omitempty"`
		UpdatedAt           time.Time   `json:"updated_at,omitempty"`
	}
)

func (p *PositionListRes) GetCommon() Response {
	return p.Response
}
