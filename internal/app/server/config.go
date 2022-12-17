package server

type ServerConfig struct {
	srvAddress string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL,required"`
}
