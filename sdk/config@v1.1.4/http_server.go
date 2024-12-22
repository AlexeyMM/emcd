package config

type HTTPServer struct {
	ListenAddr string `env:"HTTP_LISTEN_ADDR" envDefault:":8080"`
}
