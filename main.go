package main

import (
	"log"

	"file_analyzer/db"
	"file_analyzer/routes"

	"github.com/gin-gonic/gin"
)

func main () {
	r := gin.Default()
	dB := db.ConnectDB()
	routes.RegisterRoutes(r, dB)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}else {
		log.Println("✅ Server running on http://localhost:8080")
	}
}
