package character

import "dragonball/business/core/dragon-ball/character"

type Character struct {
	ID    int
	Name  string
	Ki    string
	Image string
	Race  string
}

func formatCharacter(character character.Response) Character {
	return Character{
		ID:    character.ID,
		Name:  character.Name,
		Ki:    character.Ki,
		Image: character.Image,
		Race:  character.Race,
	}
}
