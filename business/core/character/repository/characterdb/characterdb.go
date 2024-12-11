package characterdb

import (
	"context"
	"database/sql"
	"dragonball/business/core/character"
	"dragonball/foundation/database/pgx"
	"dragonball/foundation/logger"

	"github.com/jmoiron/sqlx"
)

type Storer struct {
	log *logger.Logger
	db  *sqlx.DB
}

func New(log *logger.Logger, db *sqlx.DB) *Storer {
	return &Storer{
		log: log,
		db:  db,
	}
}

func (s *Storer) QueryByName(ctx context.Context, name string) (character.Character, error) {
	data := map[string]any{
		"name": name,
	}

	query := pgx.ParseQuery(queryByName, data)

	var dbbalance dbCharacter
	err := pgx.RunQuery(ctx, s.db, query, &dbbalance)
	switch {
	case err == sql.ErrNoRows:
		return character.Character{}, nil
	case err != nil:
		return character.Character{}, err
	}

	usrs := toCoreCharacter(dbbalance)

	return usrs, nil
}

func (s *Storer) Save(ctx context.Context, character character.Character) error {
	data := map[string]any{
		"ID":    character.ID,
		"name":  character.Name,
		"ki":    character.Ki,
		"race":  character.Race,
		"image": character.Image,
	}

	query := pgx.ParseQuery(querySave, data)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}
