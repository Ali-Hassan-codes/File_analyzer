package main

import (
	"log"

	"file_analyzer/db"
	"file_analyzer/routes"
	"file_analyzer/controllers" // ✅ add this for Init functions

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ✅ Connect to database
	dB := db.ConnectDB()

	// ✅ Initialize tables (create if not exists)
	controllers.InitUsersTable(dB)
	controllers.InitFileStatsTable(dB)

	// ✅ Register routes
	routes.RegisterRoutes(r, dB)

	// ✅ Run server
	if err := r.Run(":8001"); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	} else {
		log.Println("✅ Server running on http://localhost:8001")
	}
}
