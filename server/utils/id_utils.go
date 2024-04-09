package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetIDFromRequest Get the ID from the request
func GetIDFromRequest(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("ID not provided")
		return "", fmt.Errorf("ID not provided")
	}

	return id, nil
}

// Function to generate a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// GenerateRandomID Function to generate custom 8 character ID
func GenerateRandomID() string {
	// Generate random string
	randomString := generateRandomString(10)

	// Generate timestamp as salt
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Concatenate random string and timestamp
	concatenated := randomString + timestamp

	// Hash the concatenated string using SHA-256
	hash := sha256.New()
	hash.Write([]byte(concatenated))
	hashed := hash.Sum(nil)

	// Convert hashed bytes to hexadecimal string
	customID64 := fmt.Sprintf("%x", hashed)

	// Uses builder to get every 8th character from customID64, resulting in only 8 characters
	var customIDBuilder strings.Builder
	for i := 0; i < len(customID64); i++ {
		if i%8 == 0 {
			customIDBuilder.WriteRune(rune(customID64[i]))
		}
	}

	return customIDBuilder.String()
}
