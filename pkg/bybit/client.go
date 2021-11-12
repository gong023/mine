package bybit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	Host     = "https://api.bybit.com"
	TestHost = "https://api-testnet.bybit.com"
)

var client = &http.Client{}

type Client struct {
	host      string
	authParam *authParam
}

func NewClient(host, key, secret string) *Client {
	return &Client{
		host: host,
		authParam: &authParam{
			apiKey:    key,
			apiSecret: secret,
		},
	}
}

func (c *Client) GetWalletBalance(ctx context.Context, req *WalletBalanceReq) (*WalletBalanceRes, error) {
	res := &WalletBalanceRes{}
	err := c.doGet(ctx, req, res)
	return res, err
}

func (c *Client) OrderCreate(ctx context.Context, req *OrderCreateReq) (*OrderCreateRes, error) {
	res := &OrderCreateRes{}
	err := c.doPost(ctx, req, res)
	return res, err
}

func (c *Client) OrderCancel(ctx context.Context, req *OrderCancelReq) (*OrderCancelRes, error) {
	res := &OrderCancelRes{}
	err := c.doPost(ctx, req, res)
	return res, err
}

func (c *Client) PositionLeverageSave(ctx context.Context, req *PositionLeverageSaveReq) (*PositionLeverageSaveRes, error) {
	res := &PositionLeverageSaveRes{}
	err := c.doPost(ctx, req, res)
	return res, err
}

func (c *Client) doGet(ctx context.Context, req RequestType, res ResponseType) error {
	param, _, err := c.authParam.composeQuery(req)
	if err != nil {
		return err
	}
	sign, _, err := c.authParam.sign(req)
	if err != nil {
		return err
	}
	param += "&sign=" + sign

	url := fmt.Sprintf("%s%s?%s", c.host, req.Path(), param)
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return c.doRequest(ctx, httpReq, res)
}

func (c *Client) doPost(ctx context.Context, req RequestType, res ResponseType) error {
	body, err := c.authParam.composeBody(req)
	if err != nil {
		return err
	}

	url := c.host + req.Path()
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return c.doRequest(ctx, httpReq, res)
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, res ResponseType) error {
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer func(b io.ReadCloser) { _ = b.Close() }(response.Body)

	if err := json.Unmarshal(body, res); err != nil {
		return err
	}
	if r := res.GetCommon(); r.RetMsg != RetMsgOK {
		return NewResponseError(r, body, req.Method, req.URL)
	}
	return nil
}
