package handler

import (
	"fmt"
	"strings"

	"github.com/gong023/mine/internal/env"
	"github.com/gong023/mine/pkg/bybit"
)

type (
	Conclusion string
	Trend      string
	Position   string

	Decision struct {
		config   env.Config
		webhook  []byte
		position bybit.Position
	}
)

const (
	DoNothing        Conclusion = "DoNothing"
	Long             Conclusion = "Long"
	Short            Conclusion = "Short"
	ReleaseThenLong  Conclusion = "ReleaseThenLong"
	ReleaseThenShort Conclusion = "ReleaseThenShort"

	trendUnknown Trend = "TrendUnknown"
	trendUp      Trend = "TrendUp"
	trendDown    Trend = "TrendDown"

	positionUnknown Position = "PositionUnknown"
	positionNone    Position = "PositionNone"
	positionLong    Position = "PositionLong"
	positionShort   Position = "PositionShort"
)

func (d *Decision) Make() Conclusion {
	switch d.ParseWebhook() {
	case trendUp:
		switch d.GetPosition() {
		case positionLong:
			return DoNothing
		case positionShort:
			return ReleaseThenLong
		case positionNone:
			return Long
		}
	case trendDown:
		switch d.GetPosition() {
		case positionLong:
			return ReleaseThenShort
		case positionShort:
			return DoNothing
		case positionNone:
			return Short
		}
	}

	return DoNothing
}

func (d *Decision) GetPosition() Position {
	if d.position.Side == bybit.SideBuy {
		return positionLong
	}
	if d.position.Side == bybit.SideSell {
		return positionShort
	}
	if d.position.Side == bybit.SideNone {
		return positionNone
	}
	return positionUnknown
}

func (d *Decision) ParseWebhook() Trend {
	w := string(d.webhook)
	if strings.Contains(w, d.config.TrendUp) {
		return trendUp
	}
	if strings.Contains(w, d.config.TrendDown) {
		return trendDown
	}
	return trendUnknown
}

func (d *Decision) String() string {
	return fmt.Sprintf(
		"trend:%s, position:%s, conclusion:%s",
		d.ParseWebhook(), d.GetPosition(), d.Make(),
	)
}
