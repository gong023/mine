package env

import (
	"context"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		BybitKey string `env:"BYBIT_KEY"`
		BybitSec string `env:"BYBIT_SEC"`
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
