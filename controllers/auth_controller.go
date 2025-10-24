package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var SecretKey = []byte("my_secret_key")

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ---------------- SIGNUP ----------------
func Signup(c *gin.Context, db *sql.DB) {
	var user User

	// Parse request body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if email exists
	var exists int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, user.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}
	if exists > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// ✅ FIXED: Correct PostgreSQL syntax
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3);`
	_, err = db.Exec(query, user.Username, user.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// ---------------- LOGIN ----------------
func Login(c *gin.Context, db *sql.DB) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var storedPassword, username string
	err := db.QueryRow(`SELECT username, password FROM users WHERE email = $1`, loginData.Email).
		Scan(&username, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    loginData.Email,
		"username": username,
		"exp":      jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
	})
	tokenString, _ := token.SignedString(SecretKey)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    tokenString,
		"username": username,
	})
}
