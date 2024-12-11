package charactergrp

import (
	"context"
	"dragonball/business/core/character"
	v1 "dragonball/business/web/v1"
	"dragonball/foundation/web"
	"errors"
	"net/http"
)

const empty = ""

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
	if name == empty {
		return v1.NewRequestError(errors.New("name is required for the search"), http.StatusBadRequest)
	}
	character, err := h.CharacterCore.Query(ctx, name)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, character, http.StatusOK)
}
