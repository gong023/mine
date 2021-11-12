package env

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	BybitHost string `env:"BYBIT_HOST" envDefault:"https://api-testnet.bybit.com"`
	BybitKey  string `env:"BYBIT_KEY,required"`
	BybitSec  string `env:"BYBIT_SEC,required"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
