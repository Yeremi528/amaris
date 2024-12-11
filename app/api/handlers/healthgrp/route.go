package healthgrp

import (
	"dragonball/business/web/v1/mid"
	"dragonball/foundation/logger"
	"dragonball/foundation/web"
	"net/http"
	"time"
)

type Config struct {
	Since time.Time
	Build string
	Cores string
	Log   *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, group string, cfg Config) {

	hgh := New(cfg.Build, cfg.Since, cfg.Cores)

	app.CustomHandle(http.MethodGet, group, "/health", hgh.Health, mid.Errors(cfg.Log), mid.Panics())
}
