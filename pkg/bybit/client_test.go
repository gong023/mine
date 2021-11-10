package bybit

import (
	"regexp"
	"testing"
)

func TestToSortedURLValues(t *testing.T) {
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
			vals, ts, err := client.toSortedParamString(c.req)
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
	_, err := client.sign(&WalletBalanceReq{})
	if err != nil {
		t.Fatal(err)
	}
}
