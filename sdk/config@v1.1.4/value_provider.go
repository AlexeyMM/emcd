package config

import (
	"net/url"

	"github.com/caarlos0/env/v10"
)

type ValueProvider interface {
	Value(p *url.URL) (string, error)
	Scheme() string
}

func ValuesProviderToOnSetFn(providers ...ValueProvider) env.ParserFunc {
	registry := make(map[string]ValueProvider, len(providers))
	for _, p := range providers {
		registry[p.Scheme()] = p
	}
	return func(s string) (interface{}, error) {
		url, err := url.Parse(s)
		if err != nil {
			return s, nil
		}
		provider, ok := registry[url.Scheme]
		if !ok {
			return s, nil
		}
		v, err := provider.Value(url)
		if err != nil {
			return s, err
		}
		return v, err
	}
}
