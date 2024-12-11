package characterdb

var (
	queryByName = `
	SELECT
		"ID", name, ki, race, image
	FROM
		"dragon_ball"."characters"
	WHERE
		name = :name
	`
	querySave = `
	INSERT INTO dragon_ball.characters
		("ID", name, ki, race, image)
	VALUES
		(:ID, :name, :ki, :race, :image);`
)
