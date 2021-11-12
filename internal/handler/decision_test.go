package handler

import (
	"testing"

	"github.com/gong023/mine/pkg/bybit"

	"github.com/gong023/mine/internal/env"
)

func TestDecision_Make(t *testing.T) {
	cases := map[string]struct {
		webhook      string
		positionSide string
		want         Conclusion
	}{
		"unknown_trend": {
			webhook:      "",
			positionSide: bybit.SideSell,
			want:         DoNothing,
		},
		"unknown_position": {
			webhook:      "up",
			positionSide: "",
			want:         DoNothing,
		},
		"up_long": {
			webhook:      "up",
			positionSide: bybit.SideBuy,
			want:         DoNothing,
		},
		"up_short": {
			webhook:      "up",
			positionSide: bybit.SideSell,
			want:         ReleaseThenLong,
		},
		"up_none": {
			webhook:      "up",
			positionSide: bybit.SideNone,
			want:         Long,
		},
		"down_long": {
			webhook:      "down",
			positionSide: bybit.SideBuy,
			want:         ReleaseThenShort,
		},
		"down_short": {
			webhook:      "down",
			positionSide: bybit.SideSell,
			want:         DoNothing,
		},
		"down_none": {
			webhook:      "down",
			positionSide: bybit.SideNone,
			want:         Short,
		},
	}

	e := env.Config{TrendUp: "up", TrendDown: "down"}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			d := &Decision{
				config:   e,
				webhook:  []byte(c.webhook),
				position: bybit.Position{Side: c.positionSide},
			}
			if got := d.Make(); c.want != got {
				t.Fatalf("want: %s, got: %s", c.want, got)
			}
		})
	}
}

func TestDecision_String(t *testing.T) {
	d := &Decision{config: env.Config{TrendUp: "up", TrendDown: "down"}}
	got := d.String()
	want := "trend:TrendUnknown, position:PositionUnknown, conclusion:DoNothing"
	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
