package otel

import (
	"context"
	"dragonball/foundation/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	ServiceName    string
	Host           string
	ExcludedRoutes map[string]struct{}
	Probability    float64
}

func InitTracing(log *logger.Logger, cfg Config) (trace.TracerProvider, func(ctx context.Context), error) {

	// TO DO
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("otel-collector:4312"),
		),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("creating new exporter: %w", err)
	}

	var traceProvider trace.TracerProvider

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(newEndpointExcluder(cfg.ExcludedRoutes, cfg.Probability)),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
			),
		),
	)

	teardown := func(ctx context.Context) {
		tp.Shutdown(ctx)
	}

	traceProvider = tp

	otel.SetTracerProvider(traceProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return traceProvider, teardown, nil
}

func InjectTracing(ctx context.Context, tracer trace.Tracer) context.Context {
	ctx = setTracer(ctx, tracer)

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	if traceID == "00000000000000000000000" {
		traceID = uuid.NewString()
	}

	ctx = setTraceID(ctx, traceID)

	return ctx
}

func AddSpan(ctx context.Context, spanName string, keyValues ...attribute.KeyValue) (context.Context, trace.Span) {
	v, ok := ctx.Value(tracerKey).(trace.Tracer)
	if !ok || v == nil {
		return ctx, trace.SpanFromContext(ctx)
	}

	ctx, span := v.Start(ctx, spanName)
	for _, kv := range keyValues {
		span.SetAttributes(kv)
	}

	return ctx, span
}

func AddTraceToRequest(ctx context.Context, r *http.Request) {
	hc := propagation.HeaderCarrier(r.Header)
	otel.GetTextMapPropagator().Inject(ctx, hc)
}
