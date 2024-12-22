package config

type Environment struct {
	Name string `env:"ENVIRONMENT,required"`
}

func (e Environment) IsProduction() bool {
	return e.Name == "production"
}
