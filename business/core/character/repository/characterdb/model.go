package characterdb

import (
	"dragonball/business/core/character"
)

type dbCharacter struct {
	ID    int    `db:"ID"`
	Name  string `db:"name"`
	Ki    string `db:"ki"`
	Race  string `db:"race"`
	Image string `db:"image"`
}

func toCoreCharacter(dbCharacter dbCharacter) character.Character {
	return character.Character{
		ID:    dbCharacter.ID,
		Name:  dbCharacter.Name,
		Ki:    dbCharacter.Ki,
		Image: dbCharacter.Image,
		Race:  dbCharacter.Race,
	}
}
