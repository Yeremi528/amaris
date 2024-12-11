package dragonball

import (
	"context"
	"crypto/tls"
	"dragonball/foundation/logger"
	"dragonball/foundation/modelvalidator"
	"dragonball/foundation/simplehttp"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Core struct {
	Client *simplehttp.Client
	log    *logger.Logger
}

type Config struct {
	BaseURL          string        `validate:"required"`
	RetryCount       int           `validate:"required"`
	RetryMaxWaitTime time.Duration `validate:"required"`
	Timeout          time.Duration `validate:"required"`
}

func NewCore(log *logger.Logger, cfg Config) (*Core, error) {
	if err := modelvalidator.Check(cfg, false); err != nil {
		return nil, err
	}

	fmt.Println(cfg.RetryCount, "retry ")
	fmt.Println(cfg.RetryMaxWaitTime, "RetryMaxWaitTime")
	cli := simplehttp.New(simplehttp.Config{
		BaseURL:          cfg.BaseURL,
		RetryCount:       cfg.RetryCount,
		RetryMaxWaitTime: cfg.RetryMaxWaitTime,
		Timeout:          cfg.Timeout,
	})

	cli.SetHeader("Content-Type", "application/json")
	cli.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	cli.EnableTrace()

	return &Core{
		log:    log,
		Client: cli,
	}, nil
}

// OnAfterResponse returns a function that logs the request and response, and generates an error if it did not succeed (200 <= code  <= 299)
func (c *Core) OnAfterResponse(ctx context.Context, service string, masked ...string) resty.ResponseMiddleware {
	fn := func(rc *resty.Client, r *resty.Response) error {

		req, err := json.Marshal(r.Request.Body)
		if err != nil {
			c.log.Error(ctx, "Failed to json marshal request", "error", err)
		}

		res, err := json.Marshal(r.Result())
		if err != nil {
			c.log.Error(ctx, "Failed to json marshal response", "error", err)
		}

		msg := fmt.Sprintf("%s - PMC Integration %s - ", r.Request.Method, service)
		fti := simplehttp.FormatTraceInfo(r.Request.TraceInfo())
		args := []any{
			"URL", r.Request.URL,
			"status", r.StatusCode(),
			"request", string(req),
			"response", string(res),
			"since", fti.TotalTime,
			"traceInfo", fti}

		if r.IsSuccess() {
			c.log.Info(ctx, msg+"OK", args...)
		} else {
			c.log.Warn(ctx, msg+"NOK", args...)
		}

		return nil
	}

	return fn
}

// OnError logs the attempted request when an error occurs while invoking a service.
func (c *Core) OnError(ctx context.Context, service string, masked ...string) resty.ErrorHook {
	fn := func(r *resty.Request, nativeErr error) {
		msg := fmt.Sprintf("%s - Dragon Ball %s - ERROR", r.Method, service)
		traceInfo := simplehttp.FormatTraceInfo(r.TraceInfo())
		traceInfo.TotalTime = "0s"

		args := []any{
			"error", nativeErr,
			"URL", r.URL,
			"request", r.Result,
			"since", traceInfo.TotalTime,
			"traceInfo", traceInfo}

		c.log.Error(ctx, msg, args...)
	}

	return fn
}
