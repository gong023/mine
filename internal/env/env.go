package env

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	BybitHost  string  `env:"BYBIT_HOST" envDefault:"https://api-testnet.bybit.com"`
	BybitKey   string  `env:"BYBIT_KEY,required"`
	BybitSec   string  `env:"BYBIT_SEC,required"`
	TrendUp    string  `env:"TREND_UP" envDefault:"上昇"`
	TrendDown  string  `env:"TREND_UP" envDefault:"下降"`
	Leverage   float64 `env:"leverage" envDefault:"30.0"`
	Symbol     string  `env:"symbol" envDefault:"BTCUSD"`
	UseBalance float64 `env:"use_balance" envDefault:"0.5"`
}

func New() (cfg Config, err error) {
	err = env.Parse(&cfg)
	return cfg, err
}
