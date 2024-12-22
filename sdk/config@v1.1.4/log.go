package config

type Log struct {
	Level string `env:"LOG_LEVEL,required"`
}
