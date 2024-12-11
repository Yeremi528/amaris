package charactergrp

import (
	"context"
	"dragonball/business/core/character"
	"dragonball/foundation/web"
	"net/http"
)

type Handlers struct {
	CharacterCore *character.Core
}

func New(characterCore *character.Core) *Handlers {
	return &Handlers{
		CharacterCore: characterCore,
	}
}

func (h *Handlers) character(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get("name")
	character, err := h.CharacterCore.Query(ctx, name)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, character, http.StatusOK)
}
