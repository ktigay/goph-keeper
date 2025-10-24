package config

import (
	h "github.com/ktigay/goph-keeper/internal/config"
)

const (
	defaultLogLevel = "debug"
)

// Config конфигурация.
type Config struct {
	ServerGRPCHost string `env:"GRPC_ADDRESS" arg:"-a" json:"grpc_host" help:"gRPC server address"`
	LogLevel       string `env:"LOG_LEVEL" arg:"-l" json:"log_level" help:"log level"`
	DatabaseDSN    string `env:"DATABASE_URI" arg:"-d" json:"database_uri" help:"database URI"`
	AuthSecret     string `env:"JWT_SECRET" arg:"-s" json:"jwt_secret" help:"jwt secret"`
	ConfigFile     string `env:"CONFIG" arg:"-c" help:"JSON config file path"`
}

// New конструктор.
func New(arguments []string) (*Config, error) {
	config := &Config{}

	handler := NewDefaultHandler(
		h.NewFileHandler(
			h.NewEnvHandler(
				h.NewArgumentsHandler[Config](
					nil,
					arguments,
				)),
			arguments,
		),
	)

	return handler.Handle(config)
}

// DefaultHandler дефолтные значения.
type DefaultHandler struct {
	next h.Handler[Config]
}

// Handle обработчик.
func (d *DefaultHandler) Handle(c *Config) (*Config, error) {
	c.LogLevel = defaultLogLevel

	return d.next.Handle(c)
}

// NewDefaultHandler конструктор.
func NewDefaultHandler(next h.Handler[Config]) *DefaultHandler {
	return &DefaultHandler{
		next: next,
	}
}
