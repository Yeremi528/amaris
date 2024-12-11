package character

import (
	"dragonball/business/core/dragon-ball/character"
	"fmt"
	"strings"
)

type Character struct {
	ID    int
	Name  string
	Ki    string
	Image string
	Race  string
}

func formatCharacter(characters []character.Response, name string) Character {
	for _, character := range characters {
		if character.Name == name {
			return Character{
				ID:    character.ID,
				Name:  character.Name,
				Ki:    character.Ki,
				Image: character.Image,
				Race:  character.Race,
			}
		}
	}

	return Character{}

}
func formatName(name string) string {
	name = strings.ReplaceAll(name, " ", "%20")
	fmt.Println(name, "name")
	return name
}
