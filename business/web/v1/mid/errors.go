package mid

import (
	"context"
	v1 "dragonball/business/web/v1"
	"dragonball/business/web/v1/response"
	"dragonball/foundation/logger"
	"dragonball/foundation/modelvalidator"
	"dragonball/foundation/web"
	"net/http"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				v := web.GetValues(ctx)
				var er response.Error
				var status int
				switch {

				// ERROR 400 - MODEL VALIDATOR
				case modelvalidator.IsFieldErrors(err):
					status = http.StatusBadRequest
					modErr := modelvalidator.GetFieldErrors(err)

					er.Message = modErr.Error()
					er.TraceID = v.TraceID
					er.AdditionalInfo = modErr.Fields()

				// ERROR 400 - BAD REQUEST
				case v1.IsRequestError(err):

					reqErr := v1.GetRequestError(err)

					status = reqErr.Status

					er.Message = http.StatusText(reqErr.Status)
					er.TraceID = v.TraceID

				// ERROR 500 - Internal Server error
				default:

					status = http.StatusInternalServerError
					er.Message = http.StatusText(http.StatusInternalServerError)
					er.TraceID = v.TraceID

				}

				log.Errorc(ctx, 4, "handling error coming out of the call chain",
					"code", status, "errMessage", er.Message, "errDetails", err.Error())

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
