package repository

import (
	"database/sql"
	"fmt"

	"github.com/DmitriyZhevnov/library/src/entities"
)

func FindAll(db *sql.DB) (book []entities.Book, err error) {
	rows, err := db.Query("select * from book")
	return buildBooks(rows, err)
}

func FindById(db *sql.DB, id int) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where id = '%d'", id)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func FindByName(db *sql.DB, name string) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where name = '%s'", name)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func FilterByGenre(db *sql.DB, id int) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where genre_id = '%d'", id)
	rows, err := db.Query(sqlRequest)
	return buildBooks(rows, err)
}

func FilterByPrices(db *sql.DB, min, max float64) (book []entities.Book, err error) {
	sqlRequest := fmt.Sprintf("select * from book where price >= '%f' AND price <= '%f'", min, max)
	rows, err := db.Query(sqlRequest)
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
	sqlRequest := fmt.Sprintf("insert into book(name, price, genre_id, amount) values ('%s', '%f', '%d', '%d')",
		book.Name, book.Price, book.GenreId, book.Amount)
	result, err := db.Exec(sqlRequest)
	if err != nil {
		return err
	} else {
		id, _ := result.LastInsertId()
		book.Id = int(id)
		return nil
	}

}

func Update(db *sql.DB, id int64, book *entities.Book) (int64, error) {
	sqlRequest := fmt.Sprintf("update book set name = '%s', price = '%f', genre_id = '%d', amount = '%d' where id = '%d'",
		book.Name, book.Price, book.GenreId, book.Amount, id)
	result, err := db.Exec(sqlRequest)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func Delete(db *sql.DB, id int64) (int64, error) {
	sqlRequest := fmt.Sprintf("delete from book where id = '%d'", id)
	result, err := db.Exec(sqlRequest)
	if err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
