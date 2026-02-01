package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT token for a user and family.
func GenerateToken(userID, familyID uint, secret []byte, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"family_id": familyID,
		"exp":       time.Now().Add(duration).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
