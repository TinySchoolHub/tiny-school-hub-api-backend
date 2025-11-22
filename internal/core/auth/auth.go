package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

// Argon2 parameters for password hashing
const (
	argon2Time        = 1
	argon2Memory      = 64 * 1024
	argon2Threads     = 4
	argon2KeyLen      = 32
	saltLength        = 16
	refreshTokenBytes = 32
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using argon2id
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Format: base64(salt):base64(hash)
	encoded := base64.RawStdEncoding.EncodeToString(salt) + ":" + base64.RawStdEncoding.EncodeToString(hash)
	return encoded, nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, encoded string) (bool, error) {
	// Parse salt and hash
	var salt, hash []byte
	var err error

	parts := []byte(encoded)
	colonIdx := -1
	for i, b := range parts {
		if b == ':' {
			colonIdx = i
			break
		}
	}
	if colonIdx == -1 {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err = base64.RawStdEncoding.DecodeString(string(parts[:colonIdx]))
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	hash, err = base64.RawStdEncoding.DecodeString(string(parts[colonIdx+1:]))
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the password with the same salt
	computedHash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Compare hashes
	if len(computedHash) != len(hash) {
		return false, nil
	}

	for i := range computedHash {
		if computedHash[i] != hash[i] {
			return false, nil
		}
	}

	return true, nil
}

// GenerateAccessToken generates a JWT access token
func GenerateAccessToken(userID, email, role, secret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a random refresh token
func GenerateRefreshToken() (string, error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HashRefreshToken hashes a refresh token for storage
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// ValidateAccessToken validates a JWT access token
func ValidateAccessToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
