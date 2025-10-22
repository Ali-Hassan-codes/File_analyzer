package controllers

import (
	"bufio"
	"database/sql"
	"net/http"
	"strings"
	"sync"
	"time"
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

	// Read all lines from the uploaded file
	var lines []string
	scanner := bufio.NewScanner(uploadedFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// -----------------------------
	// 1️⃣ SIMPLE VERSION (no goroutines)
	// -----------------------------
	startSimple := time.Now()
	simpleStats := analyzeSimple(lines)
	simpleTime := time.Since(startSimple)

	// -----------------------------
	// 2️⃣ GOROUTINE VERSION (chunks)
	// -----------------------------
	startConcurrent := time.Now()
	concurrentStats := analyzeConcurrent(lines)
	concurrentTime := time.Since(startConcurrent)

	// -----------------------------
	// Save one result (e.g., goroutine version) to DB
	// -----------------------------
	_, err = db.Exec(`INSERT INTO file_stats 
		(para_count, line_count, word_count, char_count, alpha_count, digit_count, vowel_count, non_vowel_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		concurrentStats["para_count"], concurrentStats["line_count"], concurrentStats["word_count"],
		concurrentStats["char_count"], concurrentStats["alpha_count"], concurrentStats["digit_count"],
		concurrentStats["vowel_count"], concurrentStats["non_vowel_count"])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// -----------------------------
	// Return both results & their timings (in microseconds)
	// -----------------------------
	c.JSON(http.StatusOK, gin.H{
		"simple_result":                simpleStats,
		"simple_time_microseconds":     simpleTime.Milliseconds(),
		"goroutine_result":             concurrentStats,
		"goroutine_time_microseconds":  concurrentTime.Milliseconds(),
	})
}

// --------------------------------------------------
// 🔹 Simple Version (Normal Sequential Loop)
// --------------------------------------------------
func analyzeSimple(lines []string) map[string]int {
	paraCount, lineCount, wordCount := 0, 0, 0
	charCount, alphaCount, digitCount := 0, 0, 0
	vowelCount, nonVowelCount := 0, 0
	isParaCounted := false
	vowels := "aeiouAEIOU"

	for _, line := range lines {
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

	return map[string]int{
		"para_count":      paraCount,
		"line_count":      lineCount,
		"word_count":      wordCount,
		"char_count":      charCount,
		"alpha_count":     alphaCount,
		"digit_count":     digitCount,
		"vowel_count":     vowelCount,
		"non_vowel_count": nonVowelCount,
	}
}

// --------------------------------------------------
// 🔹 Goroutine Version (Chunked Concurrent Processing)
// --------------------------------------------------
func analyzeConcurrent(lines []string) map[string]int {
	chunkSize := 10 // lines per chunk (adjust as needed)
	numLines := len(lines)
	numChunks := (numLines + chunkSize - 1) / chunkSize

	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	total := map[string]int{}
	vowels := "aeiouAEIOU"

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > numLines {
			end = numLines
		}

		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			stats := map[string]int{}
			isParaCounted := false

			for _, line := range chunk {
				stats["line_count"]++
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					if !isParaCounted {
						stats["para_count"]++
						isParaCounted = true
					}
					words := strings.Fields(trimmed)
					stats["word_count"] += len(words)
					for _, char := range trimmed {
						stats["char_count"]++
						if unicode.IsLetter(char) {
							stats["alpha_count"]++
							if strings.ContainsRune(vowels, char) {
								stats["vowel_count"]++
							} else {
								stats["non_vowel_count"]++
							}
						} else if unicode.IsDigit(char) {
							stats["digit_count"]++
						}
					}
				} else {
					isParaCounted = false
				}
			}

			// Merge chunk results safely
			mutex.Lock()
			for k, v := range stats {
				total[k] += v
			}
			mutex.Unlock()
		}(lines[start:end])
	}

	wg.Wait()
	return total
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