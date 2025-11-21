package auth

import (
	"strings"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "mySecurePassword123",
			wantErr:  false,
		},
		{
			name:     "short password",
			password: "abc",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 100),
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify hash format (salt:hash)
				parts := strings.Split(hash, ":")
				if len(parts) != 2 {
					t.Errorf("HashPassword() hash format incorrect, got %d parts, want 2", len(parts))
				}

				// Verify hash is not empty
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "correctPassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "wrongPassword",
			hash:     hash,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "invalid hash format",
			password: password,
			hash:     "invalid",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyPassword(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samePassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Hashes should be different due to unique salts
	if hash1 == hash2 {
		t.Error("HashPassword() should produce unique hashes for same password")
	}

	// But both should verify correctly
	valid1, _ := VerifyPassword(password, hash1)
	valid2, _ := VerifyPassword(password, hash2)

	if !valid1 || !valid2 {
		t.Error("Both hashes should verify correctly")
	}
}

func TestGenerateAccessToken(t *testing.T) {
	secret := "test-secret-key"

	tests := []struct {
		name    string
		userID  string
		email   string
		role    string
		expiry  time.Duration
		wantErr bool
	}{
		{
			name:    "valid token",
			userID:  "user-123",
			email:   "test@example.com",
			role:    "TEACHER",
			expiry:  15 * time.Minute,
			wantErr: false,
		},
		{
			name:    "short expiry",
			userID:  "user-456",
			email:   "admin@example.com",
			role:    "ADMIN",
			expiry:  1 * time.Second,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateAccessToken(tt.userID, tt.email, tt.role, secret, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateAccessToken() returned empty token")
				}

				// Verify token can be parsed
				claims, err := ValidateAccessToken(token, secret)
				if err != nil {
					t.Errorf("Generated token is invalid: %v", err)
					return
				}

				if claims.UserID != tt.userID {
					t.Errorf("Token UserID = %v, want %v", claims.UserID, tt.userID)
				}
				if claims.Email != tt.email {
					t.Errorf("Token Email = %v, want %v", claims.Email, tt.email)
				}
				if claims.Role != tt.role {
					t.Errorf("Token Role = %v, want %v", claims.Role, tt.role)
				}
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	secret := "test-secret-key"
	userID := "user-123"
	email := "test@example.com"
	role := "PARENT"

	validToken, err := GenerateAccessToken(userID, email, role, secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	expiredToken, err := GenerateAccessToken(userID, email, role, secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		secret  string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			secret:  secret,
			wantErr: false,
		},
		{
			name:    "expired token",
			token:   expiredToken,
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.here",
			secret:  secret,
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateAccessToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims.UserID != userID {
					t.Errorf("Claims UserID = %v, want %v", claims.UserID, userID)
				}
				if claims.Email != email {
					t.Errorf("Claims Email = %v, want %v", claims.Email, email)
				}
				if claims.Role != role {
					t.Errorf("Claims Role = %v, want %v", claims.Role, role)
				}
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "generate token",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateRefreshToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateRefreshToken() returned empty token")
				}

				// Verify token length (should be base64 encoded)
				if len(token) < 32 {
					t.Errorf("GenerateRefreshToken() token too short: %d bytes", len(token))
				}
			}
		})
	}
}

func TestGenerateRefreshTokenUniqueness(t *testing.T) {
	token1, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	token2, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token1 == token2 {
		t.Error("GenerateRefreshToken() should generate unique tokens")
	}
}

func TestHashRefreshToken(t *testing.T) {
	token := "test-refresh-token"

	hash := HashRefreshToken(token)

	if hash == "" {
		t.Error("HashRefreshToken() returned empty hash")
	}

	if hash == token {
		t.Error("HashRefreshToken() should not return the original token")
	}

	// Verify same token produces same hash
	hash2 := HashRefreshToken(token)
	if hash != hash2 {
		t.Error("HashRefreshToken() should be deterministic")
	}

	// Verify different tokens produce different hashes
	differentToken := "different-token"
	differentHash := HashRefreshToken(differentToken)
	if hash == differentHash {
		t.Error("Different tokens should produce different hashes")
	}
}

func TestClaimsValidation(t *testing.T) {
	secret := "test-secret"

	t.Run("token with future expiry is valid", func(t *testing.T) {
		token, err := GenerateAccessToken("user-1", "user@test.com", "TEACHER", secret, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		claims, err := ValidateAccessToken(token, secret)
		if err != nil {
			t.Errorf("Token should be valid: %v", err)
		}

		if claims == nil {
			t.Error("Claims should not be nil")
		}
	})

	t.Run("token with past expiry is invalid", func(t *testing.T) {
		token, err := GenerateAccessToken("user-1", "user@test.com", "TEACHER", secret, -1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		_, err = ValidateAccessToken(token, secret)
		if err == nil {
			t.Error("Expired token should be invalid")
		}
	})
}

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyPassword(password, hash)
	}
}

func BenchmarkGenerateAccessToken(b *testing.B) {
	secret := "test-secret"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateAccessToken("user-123", "test@example.com", "TEACHER", secret, 15*time.Minute)
	}
}

func BenchmarkValidateAccessToken(b *testing.B) {
	secret := "test-secret"
	token, _ := GenerateAccessToken("user-123", "test@example.com", "TEACHER", secret, 15*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateAccessToken(token, secret)
	}
}
