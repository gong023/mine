package env

import (
	"context"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		BybitHost string `env:"BYBIT_HOST" envDefault:"https://api-testnet.bybit.com"`
		BybitKey  string `env:"BYBIT_KEY,required"`
		BybitSec  string `env:"BYBIT_SEC,required"`
	}

	key struct{}
)

func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func MustValue(ctx context.Context) Config {
	cfg, _ := ctx.Value(key{}).(Config)
	return cfg
}

func With(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, key{}, cfg)
}
