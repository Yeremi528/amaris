package character

import (
	"context"
	dragonball "dragonball/business/core/dragon-ball"
	characterService "dragonball/business/core/dragon-ball/character"
	v1 "dragonball/business/web/v1"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	var character Character

	character, err := c.repository.QueryByName(ctx, name)
	if err != nil {
		return Character{}, fmt.Errorf("character.Query: Query.QueryByName: %w", err)
	}

	name = strings.TrimSpace(name)
	nameURL := formatName(name)

	if character.Name != empty {
		return character, nil
	}

	characterArr, err := characterService.Run(ctx, c.dragonBallCore, nameURL)
	if err != nil {
		fmt.Println(err, "error")
		return Character{}, fmt.Errorf("character.Query: character.Run: %w", err)
	}

	if characterArr == nil {
		return Character{}, v1.NewRequestError(errors.New("character not found"), http.StatusBadRequest)
	}

	characterFormated := formatCharacter(characterArr, name)

	if characterFormated.Name == empty {
		return Character{}, fmt.Errorf("character not found: %w", err)
	}

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
