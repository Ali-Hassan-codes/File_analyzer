package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	var db *sql.DB
	var err error

	// Retry loop: try 5 times with 3 seconds interval
	for i := 0; i < 5; i++ {
		connStr := "host=db port=5432 user=ali password=1234 dbname=mydb sslmode=disable"
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("❌ Error opening DB: %v", err)
		} else {
			err = db.Ping()
			if err == nil {
				fmt.Println("✅ Connected to PostgreSQL successfully!")
				return db
			}
			log.Printf("❌ Error connecting to DB, retrying... (%d/5)", i+1)
		}
		time.Sleep(3 * time.Second) // wait 3 seconds before retry
	}

	log.Fatalf("❌ Could not connect to DB after multiple attempts: %v", err)
	return nil
}
