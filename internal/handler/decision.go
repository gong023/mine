package handler

import (
	"fmt"

	"github.com/gong023/mine/pkg/bybit"
)

type (
	Conclusion string
	Trend      string
	Position   string

	Decision struct {
		webhook  []byte
		position *bybit.Position
	}
)

const (
	DoNothing        Conclusion = "DoNothing"
	Long                        = "Long"
	Short                       = "Short"
	ReleaseThenLong             = "ReleaseThenLong"
	ReleaseThenShort            = "ReleaseThenShort"

	trendUnknown Trend = "TrendUnknown"
	trendUp            = "TrendUp"
	trendDown          = "TrendDown"

	positionNone  Position = "PositionNone"
	positionLong           = "PositionLong"
	positionShort          = "PositionShort"
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
	return positionNone
}

func (d *Decision) ParseWebhook() Trend {
	return trendDown
}

func (d *Decision) String() string {
	return fmt.Sprintf(
		"trend:%s, position:%s, conclusion:%s",
		d.ParseWebhook(), d.GetPosition(), d.Make(),
	)
}
