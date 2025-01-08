package main

import (
	"time"
)

type Config struct {
	Web        WebConfig
	Debug      DebugConfig
	Postgres   PostgresConfig
	DragonBall DragonBallConfig
	Otel       OtelConfig
}

type WebConfig struct {
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:180s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:0.0.0.0:8080"`
}
type DebugConfig struct {
	ReadTimeout     time.Duration `conf:"default:180s"`
	WriteTimeout    time.Duration `conf:"default:180s"`
	IdleTimeout     time.Duration `conf:"default:180s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:0.0.0.0:8081"`
}

type PostgresConfig struct {
	Host            string        `conf:"default:servicedb-go-ms-dragon-ball"`
	Port            string        `conf:"default:5432"`
	User            string        `conf:"default:postgres"`
	Password        string        `conf:"default:postgres"`
	Name            string        `conf:"default:postgres"`
	MaxIdleConns    int           `conf:"default:2"`
	ConnMaxIdleTime time.Duration `conf:"default:20s"`
	MaxOpenConns    int           `conf:"default:2"`
	EnableTLS       bool          `conf:"default:false"`
}

type DragonBallConfig struct {
	BaseURL          string        `conf:"default:https://dragonball-api.com/api"`
	RetryCount       int           `conf:"default:3"`
	RetryMaxWaitTime time.Duration `conf:"default:3s"`
	Timeout          time.Duration `conf:"default:10s"`
}

type OtelConfig struct {
	Host        string  `conf:"default:0.0.0.0:4312"`
	ServiceName string  `conf:"default:amaris"`
	Probability float64 `conf:"default:90"`
}
