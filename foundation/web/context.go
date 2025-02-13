package web

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type contextKey int

const ctxKey contextKey = 1
const defaultTraceID = "00000000-0000-000000000000"

// Values struct represents the state for each request.
type Values struct {
	TraceID       string
	Now           time.Time
	StatusCode    int
	Response      string
	RUT           string
	DeviceVersion string
	SecurityToken string
	DeviceID      string
	Token         string
	Origin        string
}

/*
	The following methods are an exception to the policy of no getters & setters.
*/

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return &Values{
			TraceID: defaultTraceID,
			Now:     time.Now().UTC(),
		}
	}

	return v
}

// GetTraceID returns the trace ID from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return defaultTraceID
	}

	return v.TraceID
}

// GetTime returns the time from the context.
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return time.Now()
	}

	return v.Now
}

func GetRut(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.RUT
}

func GetOrigin(ctx context.Context) string {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return ""
	}

	return v.Origin
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}

// SetResponse sets the status code back into the context.
func SetResponse(ctx context.Context, response string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.Response = response
}

// SetRut sets the user's RUT into the context.
func SetRut(ctx context.Context, rut string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.RUT = rut
}

func SetOrigin(ctx context.Context, origin string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.Origin = origin
}

// SetDeviceVersion sets the user's Device Version into the context.
func SetDeviceVersion(ctx context.Context, version string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.DeviceVersion = version
}

func SetSecurityToken(ctx context.Context, securityToken string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.SecurityToken = securityToken
}

func SetDeviceID(ctx context.Context, deviceID string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.DeviceID = deviceID
}

func SetToken(ctx context.Context, token string) {
	v, ok := ctx.Value(ctxKey).(*Values)
	if !ok {
		return
	}

	v.Token = token

}

// ////////////////////
func setTracer(ctx context.Context, tracer trace.Tracer) context.Context {
	return context.WithValue(ctx, ctxKey, tracer)
}
