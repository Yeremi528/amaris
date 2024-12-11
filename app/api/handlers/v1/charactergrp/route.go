package charactergrp

import (
	"dragonball/business/core/character"
	"dragonball/business/core/character/repository/characterdb"
	dragonball "dragonball/business/core/dragon-ball"
	"dragonball/foundation/logger"
	"dragonball/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log            *logger.Logger
	DB             *sqlx.DB
	DragonBallCore *dragonball.Core
}

func Routes(app *web.App, group string, cfg Config) {
	repositoryCharacter := characterdb.New(cfg.Log, cfg.DB)
	characterCore := character.New(repositoryCharacter, cfg.DragonBallCore)
	h := New(characterCore)

	app.Handle(http.MethodGet, group, "/character", h.character)
}
