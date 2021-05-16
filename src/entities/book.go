package entities

//`json:"id,omitemty"`

type Book struct {
	Id      int     `json:"id"`
	Name    string  `json:"name" validate:"required,max=100"`
	Price   float64 `json:"price" validate:"required,min=0"`
	GenreId int     `db:"genre_id" json:"genre" validate:"required"`
	Amount  int     `json:"amount" validate:"required,min=0"`
}
