package swaggergrp

import (
	"dragonball/foundation/web"
	"path"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Build string
}

func Routes(app *web.App, group string, cfg Config) {
	if cfg.Build == "local" || cfg.Build == "dev" {
		pathSwagger := path.Join("/", group, "api-docs/*")
		pathSwaggerDocs := path.Join("/", group, "/api-docs/doc.json")
		app.HandleFunc(pathSwagger, httpSwagger.Handler(
			httpSwagger.URL(pathSwaggerDocs),
		))
	}
	return
}
