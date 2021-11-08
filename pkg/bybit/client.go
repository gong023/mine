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
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/gong023/mine/internal/env"
)

const (
	Host     = "https://api.bybit.com"
	TestHost = "https://api-testnet.bybit.com"

	SymbolBTCUSD = "BTCUSD"
)

var client = &http.Client{}

type Client struct {
	host      string
	apiKey    string
	apiSecret string
}

func NewClient(cfg env.Config) *Client {
	return &Client{
		host:      cfg.BybitHost,
		apiKey:    cfg.BybitKey,
		apiSecret: cfg.BybitSec,
	}
}

func (c *Client) GetWalletBalance(ctx context.Context, req *WalletBalanceReq) (res WalletBalanceRes, err error) {
	b, err := c.doGet(req)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("body:%s, err:%s", b, err)
	}
	return res, nil
}

func (c *Client) doGet(req GetRequest) ([]byte, error) {
	vals, err := c.toSortedURLValues(req)
	if err != nil {
		return nil, err
	}
	sign, err := c.sign(req)
	if err != nil {
		return nil, err
	}
	vals.Add("sign", sign)

	url := fmt.Sprintf("%s%s?%s", c.host, req.Path(), vals.Encode())
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")

	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) sign(r Request) (string, error) {
	q, err := c.toSortedURLValues(r)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	if _, err := io.WriteString(h, q.Encode()); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (c *Client) toSortedURLValues(r Request) (url.Values, error) {
	sb, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var srcMap map[string]interface{}
	if err := json.Unmarshal(sb, &srcMap); err != nil {
		return nil, err
	}

	srcMap["api_key"] = c.apiKey
	srcMap["timestamp"] = time.Now().Unix()

	var keys []string
	for k := range srcMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	dest := url.Values{}
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
		dest.Add(k, val)
	}

	return dest, nil
}
