package routes

import (
	"database/sql"
	"file_analyzer/controllers"
	"file_analyzer/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	// Public routes
	r.POST("/signup", func(c *gin.Context) {
		controllers.Signup(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db)
	})

	// Protected routes group (with middleware)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/upload", func(c *gin.Context) {
			controllers.UploadFile(c, db)
		})

		auth.DELETE("/files/:id", func(c *gin.Context) {
			controllers.DeleteFile(c, db)
		})
	}
}
