package mid

import (
	"context"
	"dragonball/foundation/logger"
	"dragonball/foundation/web"
	"net/http"
	"strings"
	"time"
)

// Logger is a middleware that logs information about incoming HTTP requests and outgoing responses using the provided logger.
// It wraps the given handler and adds logging functionality to it.
func Logger(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			arrURI := strings.Split(r.RequestURI, "?")
			URI := arrURI[0]

			log.Info(ctx, "request started "+URI, "protocol", r.Proto, "method", r.Method, "URL", URI, "userAgent", r.UserAgent(),
				"remoteAddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed "+URI, "protocol", r.Proto, "method", r.Method, "URL", URI, "userAgent", r.UserAgent(),
				"remoteAddr", r.RemoteAddr, "statusCode", v.StatusCode, "since", time.Since(v.Now).String())

			return err
		}

		return h
	}

	return m
}
