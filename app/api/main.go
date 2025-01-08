package main

import (
	"context"
	"dragonball/app/api/handlers"
	dragonball "dragonball/business/core/dragon-ball"
	"dragonball/foundation/database/pgx"
	"dragonball/foundation/debug"
	"dragonball/foundation/logger"
	"dragonball/foundation/otel"
	"dragonball/foundation/web"
	"errors"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/conf/v3"
)

var build = "local"
var enviroment = os.Getenv("ENV")

// @title           Spending Line API
// @version         1.0
// @description     This is documentation for the SPENDING LINE API.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /spending-line

// @securityDefinitions.apikey  JWT
// @in                          header
// @name                        jwt-token
// @description                 JWT Token for API authentication

// @securityDefinitions.apikey  DeviceID
// @in                          header
// @name                        device-id
// @description                 Unique identifier for the device

// @securityDefinitions.apikey  SecurityToken
// @in                          header
// @name                        x-security-token
// @description                 Security Token for API authentication

// @security                    JWT, DeviceID, SecurityToken
func main() {
	var logLevel logger.Level
	level := os.Getenv("APP_LOG_LEVEL")
	switch level {
	case "INFO":
		logLevel = logger.LevelInfo
	case "DEBUG":
		logLevel = logger.LevelDebug
	case "WARN":
		logLevel = logger.LevelWarn
	case "ERROR":
		logLevel = logger.LevelError
	default:
		level = "INFO"
		logLevel = logger.LevelInfo
	}

	ctx := context.Background()
	traceFunc := func(ctx context.Context) []any {
		v := web.GetValues(ctx)

		fields := make([]any, 2, 4)
		fields[0], fields[1] = "traceID", v.TraceID
		return fields
	}
	log := logger.New(os.Stdout, logLevel, "go-ms-dragon-ball", traceFunc)

	log.Info(ctx, "startup - service version...", "version", build, "cores", runtime.GOMAXPROCS(0), "logLevelAt", level)

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "service error, shutting down", "errorDetails", err.Error())
		os.Exit(1)
	}

}

func run(ctx context.Context, log *logger.Logger) error {
	defer log.Info(ctx, "shutdown - complete")

	/*==========================================================================
		App Configuration
	==========================================================================*/
	var cfg Config

	const prefix = "APP"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup - set config to...", "config", out)

	/*==========================================================================
		Postgres Database Support
	==========================================================================*/

	db, err := pgx.Open(pgx.Config{
		User:            cfg.Postgres.User,
		Password:        cfg.Postgres.Password,
		Host:            cfg.Postgres.Host,
		Port:            cfg.Postgres.Port,
		Name:            cfg.Postgres.Name,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		IdleConnTimeout: cfg.Postgres.ConnMaxIdleTime,
		EnableTLS:       cfg.Postgres.EnableTLS,
		ApplicationName: "go-ms-dragon-ball",
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer func() {
		log.Info(ctx, "shutdown - stopping database support", "host", cfg.Postgres.Host, "port", cfg.Postgres.Port)
		if err := db.Close(); err != nil {
			log.Info(ctx, "shutdown - cannot stop database support gracefully", "host", cfg.Postgres.Host, "port", cfg.Postgres.Port, "error", err.Error())
			return
		}

		log.Info(ctx, "shutdown - database support stopped", "host", cfg.Postgres.Host, "port", cfg.Postgres.Port)
	}()
	/*==========================================================================
		Dragon Ball
	==========================================================================*/
	dragonBallCore, err := dragonball.NewCore(log, dragonball.Config{
		BaseURL:          cfg.DragonBall.BaseURL,
		RetryCount:       cfg.DragonBall.RetryCount,
		RetryMaxWaitTime: cfg.DragonBall.RetryMaxWaitTime,
		Timeout:          cfg.DragonBall.Timeout,
	})

	if err != nil {
		return fmt.Errorf("main.run: dragonBall.NewCore: %w", err)
	}

	/*==========================================================================
		Start Tracing Support
	==========================================================================*/

	traceProvider, teardown, err := otel.InitTracing(log, otel.Config{
		ServiceName: cfg.Otel.ServiceName,
		Host:        cfg.Otel.Host,
		Probability: cfg.Otel.Probability,
	})

	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}

	defer teardown(context.Background())

	tracer := traceProvider.Tracer(cfg.Debug.APIHost)

	/*==========================================================================
		Debug Service
	==========================================================================*/

	go func() {
		log.Info(ctx, "startup - debug router started", "host", cfg.Debug.APIHost)

		apiDebug := http.Server{
			Addr:         cfg.Debug.APIHost,
			Handler:      debug.Mux(),
			ReadTimeout:  cfg.Debug.ReadTimeout,
			WriteTimeout: cfg.Debug.WriteTimeout,
			IdleTimeout:  cfg.Debug.IdleTimeout,
		}
		if err := apiDebug.ListenAndServe(); err != nil {
			log.Info(ctx, "shutdown - debug v1 router closed", "host", cfg.Debug.APIHost, "error", err.Error())
		}
	}()

	/*==========================================================================
		API Service
	==========================================================================*/

	log.Info(ctx, "startup - initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(
		handlers.APIMuxConfig{
			Enviroment:     enviroment,
			Build:          build,
			Shutdown:       shutdown,
			Log:            log,
			DB:             db,
			DragonBallCore: dragonBallCore,
			Tracer:         tracer,
		},
	)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      otelhttp.NewHandler(apiMux, "http-server"),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup - API router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	/*==========================================================================
		Wait for shutdown signal or server error
	==========================================================================*/

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown - started", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("cannot stop server gracefully: %w", err)
		}
	}

	return nil

}
