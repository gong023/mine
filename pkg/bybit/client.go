package bybit

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/iancoleman/orderedmap"
)

const (
	Host     = "https://api.bybit.com"
	TestHost = "https://api-testnet.bybit.com"
)

var client = &http.Client{}

type Client struct {
	host      string
	apiKey    string
	apiSecret string
}

func NewClient(host, key, secret string) *Client {
	return &Client{
		host:      host,
		apiKey:    key,
		apiSecret: secret,
	}
}

func (c *Client) GetWalletBalance(ctx context.Context, req *WalletBalanceReq) (res *WalletBalanceRes, err error) {
	getReq, err := c.composeGet(req)
	if err != nil {
		return res, err
	}
	rawResBody, err := c.doRequest(ctx, getReq, &res)
	if err != nil {
		return res, err
	}
	if res.Response.RetCode != RetCodeSuccess {
		return res, fmt.Errorf("failed GET request:%s, res:%s", req.Path(), rawResBody)
	}
	return res, nil
}

func (c *Client) OrderCreate(ctx context.Context, req *OrderCreateReq) (res *OrderCreateRes, err error) {
	postReq, err := c.composePost(req)
	if err != nil {
		return res, err
	}
	rawResBody, err := c.doRequest(ctx, postReq, &res)
	if err != nil {
		return res, err
	}
	if res.Response.RetCode != RetCodeSuccess {
		return res, fmt.Errorf("failed POST request:%s, resBody:%s", req.Path(), rawResBody)
	}
	return res, nil
}

func (c *Client) PositionLeverageSave(ctx context.Context, req *PositionLeverageSaveReq) (res *PositionLeverageSaveRes, err error) {
	postReq, err := c.composePost(req)
	if err != nil {
		return res, err
	}
	rawResBody, err := c.doRequest(ctx, postReq, &res)
	if err != nil {
		return res, err
	}
	if res.Response.RetCode != RetCodeSuccess {
		return res, fmt.Errorf("failed POST request:%s, resBody:%s", req.Path(), rawResBody)
	}
	return res, nil
}

func (c *Client) composeGet(req GetRequest) (*http.Request, error) {
	param, _, err := c.composeQuery(req)
	if err != nil {
		return nil, err
	}
	sign, _, err := c.sign(req)
	if err != nil {
		return nil, err
	}
	param += "&sign=" + sign

	url := fmt.Sprintf("%s%s?%s", c.host, req.Path(), param)
	return http.NewRequest(http.MethodGet, url, nil)
}

func (c *Client) composePost(req PostRequest) (*http.Request, error) {
	body, err := c.composeBody(req)
	if err != nil {
		return nil, err
	}

	url := c.host + req.Path()
	return http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
}

func (c *Client) sign(r Request) (string, int64, error) {
	q, ts, err := c.composeQuery(r)
	if err != nil {
		return "", 0, err
	}
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	if _, err := io.WriteString(h, q); err != nil {
		return "", 0, err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), ts, nil
}

func (c *Client) composeQuery(r Request) (string, int64, error) {
	sb, err := json.Marshal(r)
	if err != nil {
		return "", 0, err
	}
	var srcMap map[string]interface{}
	if err := json.Unmarshal(sb, &srcMap); err != nil {
		return "", 0, err
	}

	srcMap["api_key"] = c.apiKey
	timestamp := time.Now().UnixMilli()
	srcMap["timestamp"] = timestamp

	var keys []string
	for k := range srcMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// url.Values shouldn't work because it's map, doesn't care the order.
	dest := ""
	for _, k := range keys {
		val := ""
		switch reflect.TypeOf(srcMap[k]).Kind() {
		case reflect.Float64:
			val = strconv.FormatFloat(srcMap[k].(float64), 'f', -1, 64)
		case reflect.Int64:
			val = strconv.FormatInt(srcMap[k].(int64), 10)
		case reflect.Int:
			val = strconv.Itoa(srcMap[k].(int))
		case reflect.String:
			val = srcMap[k].(string)
		}
		dest += fmt.Sprintf("%s=%s&", k, val)
	}

	return dest[0 : len(dest)-1], timestamp, nil
}

func (c *Client) composeBody(req PostRequest) ([]byte, error) {
	sb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var srcMap map[string]interface{}
	if err := json.Unmarshal(sb, &srcMap); err != nil {
		return nil, err
	}

	sign, timestamp, err := c.sign(req)
	if err != nil {
		return nil, err
	}

	o := orderedmap.New()
	for k, v := range srcMap {
		o.Set(k, v)
	}
	o.Set("api_key", c.apiKey)
	o.Set("timestamp", timestamp)
	o.SortKeys(sort.Strings)
	o.Set("sign", sign)

	return o.MarshalJSON()
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, res interface{}) ([]byte, error) {
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("body:%s, err:%s", body, err)
	}
	return body, nil
}
