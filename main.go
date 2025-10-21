package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"file_analyzer/routes"
)

func main() {
	// Database connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/file_analysis")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Gin setup
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r, db)

	// Run server
	r.Run(":8001")
}
