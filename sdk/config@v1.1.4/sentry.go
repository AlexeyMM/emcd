package config

type Sentry struct {
	DNS string `env:"SENTRY_DNS,required"`
}
