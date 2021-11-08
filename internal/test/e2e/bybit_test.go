// +build e2e

package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/pkg/bybit"
)

func TestGetWalletBalance(t *testing.T) {
	cfg, err := env.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	cli := bybit.NewClient(cfg)
	res, err := cli.GetWalletBalance(ctx, &bybit.WalletBalanceReq{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v", res)
}
