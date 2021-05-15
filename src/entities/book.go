package entities

type Book struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	GenreId int     `db:"genre_id" json:"genre"`
	Amount  int     `json:"amount"`
}
