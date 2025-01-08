package mid

import (
	"context"
	"dragonball/foundation/otel"
	"dragonball/foundation/web"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

func Otel(tracer trace.Tracer) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = otel.InjectTracing(ctx, tracer)

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
