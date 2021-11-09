package bybit

import (
	"regexp"
	"testing"

	"github.com/gong023/mine/internal/env"
)

func TestToSortedURLValues(t *testing.T) {
	client := NewClient(env.Config{
		BybitHost: TestHost,
		BybitKey:  "dummy_key",
		BybitSec:  "dummy_sec",
	})
	cases := map[string]struct {
		req  Request
		want *regexp.Regexp
	}{
		"doc_sample": {
			req: &WalletBalanceReq{
				Coin: "BTC",
			},
			want: regexp.MustCompile(
				"api_key=dummy_key&coin=BTC&timestamp=.*",
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			vals, err := client.toSortedParamString(c.req)
			if err != nil {
				t.Fatal(err)
			}
			if !c.want.MatchString(vals) {
				t.Fatalf("want to match: %s, got: %s", c.want, vals)
			}
		})
	}
}

func TestSign(t *testing.T) {
	client := NewClient(env.Config{
		BybitHost: TestHost,
		BybitKey:  "dummy_key",
		BybitSec:  "dummy_sec",
	})
	_, err := client.sign(&WalletBalanceReq{})
	if err != nil {
		t.Fatal(err)
	}
}
