package controllers

import (
	"bufio"
	"database/sql"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context, db *sql.DB) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open uploaded file"})
		return
	}
	defer uploadedFile.Close()

	paraCount, lineCount, wordCount := 0, 0, 0
	charCount, alphaCount, digitCount := 0, 0, 0
	vowelCount, nonVowelCount := 0, 0
	isParaCounted := false
	vowels := "aeiouAEIOU"

	scanner := bufio.NewScanner(uploadedFile)
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			if !isParaCounted {
				paraCount++
				isParaCounted = true
			}
			words := strings.Fields(trimmed)
			wordCount += len(words)
			for _, char := range trimmed {
				charCount++
				if unicode.IsLetter(char) {
					alphaCount++
					if strings.ContainsRune(vowels, char) {
						vowelCount++
					} else {
						nonVowelCount++
					}
				} else if unicode.IsDigit(char) {
					digitCount++
				}
			}
		} else {
			isParaCounted = false
		}
	}

	_, err = db.Exec(`INSERT INTO file_stats 
		(para_count, line_count, word_count, char_count, alpha_count, digit_count, vowel_count, non_vowel_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		paraCount, lineCount, wordCount, charCount, alphaCount, digitCount, vowelCount, nonVowelCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"para_count":      paraCount,
		"line_count":      lineCount,
		"word_count":      wordCount,
		"char_count":      charCount,
		"alpha_count":     alphaCount,
		"digit_count":     digitCount,
		"vowel_count":     vowelCount,
		"non_vowel_count": nonVowelCount,
	})
}

func DeleteFile(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM file_stats WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
