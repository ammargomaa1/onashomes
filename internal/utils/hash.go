package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomSKU generates a random SKU string in the format SKU-XXXXXXXX
func GenerateRandomSKU() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback if crypto/rand fails
		return fmt.Sprintf("SKU-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("SKU-%s", strings.ToUpper(hex.EncodeToString(bytes)))
}
