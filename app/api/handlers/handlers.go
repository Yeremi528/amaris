// Package handlers manages the different versions of the API.
package handlers

import (
	"context"
	"dragonball/app/api/handlers/healthgrp"
	"dragonball/app/api/handlers/swaggergrp"
	v1 "dragonball/app/api/handlers/v1"
	dragonball "dragonball/business/core/dragon-ball"
	"dragonball/business/web/v1/mid"
	"dragonball/foundation/logger"
	"dragonball/foundation/otel"
	"dragonball/foundation/timecl"
	"dragonball/foundation/web"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

const (
	group = "/dragon-ball"
)

// APIMuxConfig contains all the values required by APIMux (http.Handler).
type APIMuxConfig struct {
	Enviroment        string
	Build             string
	Shutdown          chan os.Signal
	Log               *logger.Logger
	DefaultHTTPClient http.Client
	DB                *sqlx.DB
	DragonBallCore    *dragonball.Core
	Tracer            trace.Tracer
}

// APIMux constructs an http.Handler that contains the app routes.
func APIMux(cfg APIMuxConfig) http.Handler {
	cores := strconv.Itoa(runtime.GOMAXPROCS(0))
	startTime := timecl.Now()

	otel.AddSpan(context.Background(), "hello woooorl")
	app := web.NewApp(
		cfg.Shutdown, cfg.Tracer, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Panics(),
		mid.Otel(cfg.Tracer),
	)

	healthgrp.Routes(app, group, healthgrp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		Since: startTime,
		Cores: cores,
	})

	swaggergrp.Routes(app, group, swaggergrp.Config{
		Build: cfg.Enviroment,
	})

	v1.Routes(app, group, v1.Config{
		Log:            cfg.Log,
		DB:             cfg.DB,
		DragonBallCore: cfg.DragonBallCore,
	})
	return app
}
