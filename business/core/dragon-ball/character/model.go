package character

type Response struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Ki          string      `json:"ki"`
	MaxKi       string      `json:"maxKi"`
	Race        string      `json:"race"`
	Gender      string      `json:"gender"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Affiliation string      `json:"affiliation"`
	DeletedAt   interface{} `json:"deletedAt"`
}
