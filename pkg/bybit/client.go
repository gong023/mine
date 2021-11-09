package bybit

import (
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
)

const (
	Host     = "https://api.bybit.com"
	TestHost = "https://api-testnet.bybit.com"

	SymbolBTCUSD = "BTCUSD"

	RetCodeSuccess float32 = 0
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

func (c *Client) GetWalletBalance(ctx context.Context, req *WalletBalanceReq) (res WalletBalanceRes, err error) {
	b, err := c.doGet(ctx, req)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("body:%s, err:%s", b, err)
	}
	if res.Response.RetCode != RetCodeSuccess {
		return res, fmt.Errorf("failed GET request:%s, res:%s", req.Path(), b)
	}
	return res, nil
}

func (c *Client) doGet(ctx context.Context, req GetRequest) ([]byte, error) {
	param, err := c.toSortedParamString(req)
	if err != nil {
		return nil, err
	}
	sign, err := c.sign(req)
	if err != nil {
		return nil, err
	}
	param += "&sign=" + sign

	url := fmt.Sprintf("%s%s?%s", c.host, req.Path(), param)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r = r.WithContext(ctx)

	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	_ = response.Body.Close()
	return body, nil
}

func (c *Client) sign(r Request) (string, error) {
	q, err := c.toSortedParamString(r)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	if _, err := io.WriteString(h, q); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (c *Client) toSortedParamString(r Request) (string, error) {
	sb, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	var srcMap map[string]interface{}
	if err := json.Unmarshal(sb, &srcMap); err != nil {
		return "", err
	}

	srcMap["api_key"] = c.apiKey
	srcMap["timestamp"] = time.Now().UnixMilli()

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

	return dest[0 : len(dest)-1], nil
}
