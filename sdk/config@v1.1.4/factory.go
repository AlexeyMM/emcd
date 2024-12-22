package config

import (
	"fmt"
	"reflect"

	"github.com/caarlos0/env/v10"
	"github.com/pkg/errors"

	"code.emcdtech.com/emcd/sdk/config/vault"
)

func New[T any](opts ...Option) (T, error) {
	cfg := new(T)
	vaultProvider, err := vault.NewProvider()
	if err != nil {
		return *cfg, fmt.Errorf("new vault provider: %w", err)
	}

	envOpts := env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(""): ValuesProviderToOnSetFn(vaultProvider),
		},
	}
	for _, opt := range opts {
		opt(&envOpts)
	}

	err = env.ParseWithOptions(cfg, envOpts)
	if err != nil {
		return *cfg, errors.WithStack(err)
	}
	return *cfg, nil
}

type Option func(*env.Options)

func WithPrefix(prefix string) Option {
	return func(envOpts *env.Options) {
		envOpts.Prefix = prefix
	}
}
