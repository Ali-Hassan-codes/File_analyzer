package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	url := "root:@tcp(127.0.0.1:3306)/file_analysis"
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}

	createTableUser := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50),
		email VARCHAR(50) UNIQUE,
		password VARCHAR(100)
	);`
	_, err = db.Exec(createTableUser)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS file_stats (
		id INT AUTO_INCREMENT PRIMARY KEY,
		para_count INT,
		line_count INT,
		word_count INT,
		char_count INT,
		alpha_count INT,
		digit_count INT,
		vowel_count INT,
		non_vowel_count INT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
