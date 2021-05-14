package repository

import (
	"database/sql"
	"errors"

	"github.com/DmitriyZhevnov/library/src/entities"
)

type BookRepository struct {
	Db *sql.DB
}

func (bookModel BookRepository) FindAll() (book []entities.Book, err error) {
	rows, err := bookModel.Db.Query("select * from book")
	return buildBooks(rows, err)
}

func (bookModel BookRepository) FindAllByName(name string) (book []entities.Book, err error) {
	rows, err := bookModel.Db.Query("select * from book where name = ?", name)
	return buildBooks(rows, err)
}

func (bookModel BookRepository) FilterByPrices(min, max float64) (book []entities.Book, err error) {
	rows, err := bookModel.Db.Query("select * from book where price >= ? AND price <= ?", min, max)
	return buildBooks(rows, err)

}

func (bookModel BookRepository) FilterByGenre(id int) (book []entities.Book, err error) {
	rows, err := bookModel.Db.Query("select * from book where genre_id = ?", id)
	return buildBooks(rows, err)
}

func (bookModel BookRepository) FindById(id int) (book []entities.Book, err error) {
	rows, err := bookModel.Db.Query("select * from book where id = ?", id)
	return buildBooks(rows, err)
}

func buildBooks(rows *sql.Rows, er error) (book []entities.Book, err error) {
	if er != nil {
		return nil, err
	} else {
		var books []entities.Book
		for rows.Next() {
			var id int
			var name string
			var price float64
			var genre int
			var amount int
			err2 := rows.Scan(&id, &name, &price, &genre, &amount)
			if err2 != nil {
				return nil, err2
			} else {
				if amount > 0 {
					book := entities.Book{
						Id:      id,
						Name:    name,
						Price:   price,
						GenreId: genre,
						Amount:  amount,
					}
					books = append(books, book)
				}
			}
		}
		return books, nil
	}
}

func (bookModel BookRepository) Create(book *entities.Book) error {
	if len(book.Name) > 100 || book.Price < 0 || book.Amount < 0 {
		err := errors.New("Entered data not valid")
		return err
	} else {
		result, err := bookModel.Db.Exec("insert into book(name, price, genre_id, amount) values (?, ?, ?, ?)",
			book.Name, book.Price, book.GenreId, book.Amount)
		if err != nil {
			return err
		} else {
			id, _ := result.LastInsertId()
			book.Id = int(id)
			return nil
		}
	}
}

func (bookModel BookRepository) Update(id int64, book *entities.Book) (int64, error) {
	if len(book.Name) > 100 || book.Price < 0 || book.Amount < 0 {
		err := errors.New("Entered data not valid")
		return 0, err
	} else {
		result, err := bookModel.Db.Exec("update book set name = ?, price = ?, genre_id = ?, amount = ? where id = ?",
			book.Name, book.Price, book.GenreId, book.Amount, id)
		if err != nil {
			return 0, err
		} else {
			return result.RowsAffected()
		}
	}
}

func (bookModel BookRepository) Delete(id int64) (int64, error) {
	result, err := bookModel.Db.Exec("delete from book where id = ?", id)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
