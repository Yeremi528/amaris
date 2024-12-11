package v1

import (
	"dragonball/app/api/handlers/v1/charactergrp"
	dragonball "dragonball/business/core/dragon-ball"
	"dragonball/foundation/logger"
	"dragonball/foundation/web"
	"os"

	"github.com/jmoiron/sqlx"
)

const (
	version = "/v1"
)

type Config struct {
	Build          string
	Shutdown       chan os.Signal
	Log            *logger.Logger
	DB             *sqlx.DB
	DragonBallCore *dragonball.Core
}

func Routes(app *web.App, group string, cfg Config) {
	charactergrp.Routes(app, group+version, charactergrp.Config{
		Log:            cfg.Log,
		DB:             cfg.DB,
		DragonBallCore: cfg.DragonBallCore,
	})
}
