package repository

import (
	"database/sql"
	"errors"

	"github.com/DmitriyZhevnov/library/src/entities"
)

func FindAll(db *sql.DB) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book")
	return buildBooks(rows, err)
}

func FindById(db *sql.DB, id int) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book where id = ?", id)
	return buildBooks(rows, err)
}

func FindByName(db *sql.DB, name string) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book where name = ?", name)
	return buildBooks(rows, err)
}

func FilterByGenre(db *sql.DB, id int) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book where genre_id = ?", id)
	return buildBooks(rows, err)
}

func FilterByPrices(db *sql.DB, min, max float64) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book where price >= ? AND price <= ?", min, max)
	return buildBooks(rows, err)
}

func buildBooks(rows *sql.Rows, er error) (book []entities.Book, err error) {
	if er != nil {
		return nil, er
	} else {
		defer rows.Close()
		var books []entities.Book
		for rows.Next() {
			var b entities.Book
			if err := rows.Scan(&b.Id, &b.Name, &b.Price, &b.GenreId, &b.Amount); err != nil {
				return nil, err
			}
			if b.Amount > 0 {
				books = append(books, b)
			}
		}
		return books, nil
	}
}

func Create(db *sql.DB, book *entities.Book) error {
	if len(book.Name) > 100 || book.Price < 0 || book.Amount < 0 {
		err := errors.New("Entered data not valid")
		return err
	} else {
		result, err := db.Exec("insert into book(name, price, genre_id, amount) values (?, ?, ?, ?)",
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

func Update(db *sql.DB, id int64, book *entities.Book) (int64, error) {
	if len(book.Name) > 100 || book.Price < 0 || book.Amount < 0 {
		err := errors.New("Entered data not valid")
		return 0, err
	} else {
		result, err := db.Exec("update book set name = ?, price = ?, genre_id = ?, amount = ? where id = ?",
			book.Name, book.Price, book.GenreId, book.Amount, id)
		if err != nil {
			return 0, err
		} else {
			return result.RowsAffected()
		}
	}
}

func Delete(db *sql.DB, id int64) (int64, error) {
	result, err := db.Exec("delete from book where id = ?", id)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
