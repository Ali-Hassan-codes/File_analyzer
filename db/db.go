package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=file_analysis sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Error opening DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Error connecting to DB: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")
	return db
}
