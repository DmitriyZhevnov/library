package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db, err := GetDB()
	if err != nil {
		fmt.Println("Not connected to db")
	} else {
		db.Query(`DROP TABLE IF EXISTS book;`)
		db.Query(`DROP TABLE IF EXISTS genre;`)
		db.Query(`CREATE TABLE genre (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(45) NOT NULL,
			PRIMARY KEY (id));`)
		db.Query(`INSERT INTO library.genre (name) VALUES ('Adventure');`)
		db.Query(`INSERT INTO library.genre (name) VALUES ('Classics');`)
		db.Query(`INSERT INTO library.genre (name) VALUES ('Fantasy');`)
		db.Query(`CREATE TABLE book (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL,
			price DOUBLE NOT NULL,
			genre_id INT NOT NULL,
			amount INT NOT NULL,
			PRIMARY KEY (id),
			CONSTRAINT fk_book_genre
			  FOREIGN KEY (genre_id)
			  REFERENCES library.genre (id)
			  ON DELETE NO ACTION
			  ON UPDATE NO ACTION);`)
		db.Query(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('book1', '10', '1', '50');`)
		db.Query(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('book2', '11', '2', '1');`)
		db.Query(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('book3', '20.6', '3', '3');`)
		db.Query(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('book4', '25', '1', '4');`)
		db.Query(`INSERT INTO library.book (name, price, genre_id, amount) VALUES ('book5', '30.5', '2', '2');`)

	}
}

func GetDB() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "950621"
	dbName := "library"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return
}
