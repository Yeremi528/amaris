package character

import (
	"context"
	dragonball "dragonball/business/core/dragon-ball"
	"dragonball/business/core/dragon-ball/character"
	v1 "dragonball/business/web/v1"
	"errors"
	"fmt"
	"net/http"
)

const empty = ""

type Core struct {
	repository     Repository
	dragonBallCore *dragonball.Core
}

type Repository interface {
	QueryByName(ctx context.Context, name string) (Character, error)
	Save(ctx context.Context, character Character) error
}

func New(repository Repository, dbCore *dragonball.Core) *Core {

	return &Core{
		repository:     repository,
		dragonBallCore: dbCore,
	}
}

func (c *Core) QueryByName(ctx context.Context, name string) (Character, error) {
	characterDB, err := c.repository.QueryByName(ctx, name)
	if err != nil {
		return Character{}, fmt.Errorf("character.Query: Query.QueryByName: %w", err)
	}

	if characterDB.Name != empty {
		return characterDB, nil
	}

	characterSv, err := character.Run(ctx, c.dragonBallCore, name)
	if err != nil {
		return Character{}, fmt.Errorf("character.Query: character.Run: %w", err)
	}

	if characterSv[0].Name == empty {
		return Character{}, v1.NewRequestError(errors.New("character not found"), http.StatusBadRequest)
	}

	characterFormated := formatCharacter(characterSv[0])

	if err := c.repository.Save(ctx, characterFormated); err != nil {
		return Character{}, fmt.Errorf("character.Query: repository.Save: %w", err)
	}

	return characterFormated, nil
}

func (c *Core) Query(ctx context.Context, name string) (Character, error) {
	character, err := c.QueryByName(ctx, name)
	if err != nil {
		return Character{}, err
	}

	return character, nil
}
