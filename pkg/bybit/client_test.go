package bybit

import (
	"regexp"
	"strings"
	"testing"
)

func TestComposeQuery(t *testing.T) {
	client := NewClient(TestHost, "dummy_key", "dummy_sec")
	cases := map[string]struct {
		req  Request
		want *regexp.Regexp
	}{
		"simple": {
			req: &WalletBalanceReq{
				Coin: "BTC",
			},
			want: regexp.MustCompile(
				"api_key=dummy_key&coin=BTC&timestamp=.*",
			),
		},
		"omitempty": {
			req: &WalletBalanceReq{},
			want: regexp.MustCompile(
				"api_key=dummy_key&timestamp=.*",
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			vals, ts, err := client.composeQuery(c.req)
			if err != nil {
				t.Fatal(err)
			}
			if ts == 0 {
				t.Fatal("must not be 0")
			}
			if !c.want.MatchString(vals) {
				t.Fatalf("want to match: %s, got: %s", c.want, vals)
			}
		})
	}
}

func TestSign(t *testing.T) {
	client := NewClient(TestHost, "dummy_key", "dummy_sec")
	_, _, err := client.sign(&WalletBalanceReq{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestComposeBody(t *testing.T) {
	client := NewClient(TestHost, "dummy_key", "dummy_sec")
	cases := map[string]struct {
		req  PostRequest
		want *regexp.Regexp
	}{
		"simple": {
			req: &OrderCancelReq{
				Symbol:  SymbolBTCUSD,
				OrderID: "dummy_order_id",
			},
			want: regexp.MustCompile(
				`{"api_key":"dummy_key","order_id":"dummy_order_id","symbol":"BTCUSD","timestamp":.*,"sign":".*"}`,
			),
		},
		"omitempty": {
			req: &OrderCancelReq{
				Symbol: SymbolBTCUSD,
			},
			want: regexp.MustCompile(
				`{"api_key":"dummy_key","symbol":"BTCUSD","timestamp":.*,"sign":".*"}`,
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := client.composeBody(c.req)
			if err != nil {
				t.Fatal(err)
			}
			g := strings.ReplaceAll(string(got), "\n", "")
			if !c.want.MatchString(g) {
				t.Fatalf("want to match: %s, got: %s", c.want, g)
			}
		})
	}
}
