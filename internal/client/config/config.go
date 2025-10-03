package config

import (
	h "github.com/ktigay/goph-keeper/internal/config"
)

const (
	defaultServerGRPCHost    = ":5001"
	defaultLogLevel          = "debug"
	defaultLogFile           = "./client.log"
	defaultSrvSyncToInterval = 2000
	defaultSrvSyncTimeout    = 300
)

// Config конфигурация.
type Config struct {
	ServerGRPCHost    string `env:"GRPC_ADDRESS" json:"server_grpc_host" arg:"-a" help:"server host"`
	LogLevel          string `env:"LOG_LEVEL" json:"log_level" arg:"-l" help:"log level"`
	LogFile           string `env:"CONFIG" json:"log_file" arg:"-c" help:"log file path"`
	SrvSyncToInterval int64  `env:"SRV_SYNC_INTERVAL" json:"srv_sync_to_interval" arg:"-i" help:"server sync interval"`
	SrvRequestTimeout int64  `env:"SRV_REQUEST_TIMEOUT" json:"srv_request_timeout" arg:"-t" help:"server request timeout"`
	Version           bool   `arg:"-v" help:"show version"`
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
	c.ServerGRPCHost = defaultServerGRPCHost
	c.LogFile = defaultLogFile
	c.SrvSyncToInterval = defaultSrvSyncToInterval
	c.SrvRequestTimeout = defaultSrvSyncTimeout

	return d.next.Handle(c)
}

// NewDefaultHandler конструктор.
func NewDefaultHandler(next h.Handler[Config]) *DefaultHandler {
	return &DefaultHandler{
		next: next,
	}
}
