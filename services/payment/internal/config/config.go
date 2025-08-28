package config

import "time"

const AppName = "Payment"

type PaymentConfiguration struct {
	LogLevel string `env:"LOG_LEVEL,default=INFO"`
	Port     string `env:"PORT,default=8000"`
	DB       DBConfig
}

type DBConfig struct {
	DNS             string        `env:"DB_DNS,required"`
	MaxConns        int32         `env:"DB_MAX_CONNS,default=10"`
	MinConns        int32         `env:"DB_MIN_CONNS,default=0"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME,default=30m"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME,default=1h"`
}
