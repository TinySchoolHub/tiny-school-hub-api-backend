package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/config"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/auth"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key"
	userID := uuid.New().String()
	email := "test@example.com"
	role := string(domain.RoleTeacher)

	validToken, err := auth.GenerateAccessToken(userID, email, role, secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	expiredToken, err := auth.GenerateAccessToken(userID, email, role, secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "valid token",
			token:          "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "no authorization header",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid token format",
			token:          "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "expired token",
			token:          "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "wrong secret",
			token:          "Bearer " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test secret
			testSecret := secret
			if tt.name == "wrong secret" {
				testSecret = "wrong-secret"
			}

			// Create middleware with config
			cfg := &config.Config{
				JWT: config.JWTConfig{
					Secret: testSecret,
				},
			}
			middleware := AuthMiddleware(cfg)

			// Create test handler
			nextCalled := false
			var capturedReq *http.Request
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				capturedReq = r
				w.WriteHeader(http.StatusOK)
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute middleware
			handler := middleware(next)
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Status code = %v, want %v", rr.Code, tt.expectedStatus)
			}

			// Check if next handler was called
			if nextCalled != tt.expectNext {
				t.Errorf("Next handler called = %v, want %v", nextCalled, tt.expectNext)
			}

			// Check user context if next was called
			if nextCalled && capturedReq != nil {
				userID, ok := capturedReq.Context().Value(UserIDKey).(string)
				if !ok || userID == "" {
					t.Error("UserID not found in context")
				}
				userRole, ok := capturedReq.Context().Value(UserRoleKey).(string)
				if !ok || userRole == "" {
					t.Error("UserRole not found in context")
				}
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	secret := "test-secret-key"

	tests := []struct {
		name           string
		requiredRole   domain.Role
		userRole       domain.Role
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "teacher has teacher access",
			requiredRole:   domain.RoleTeacher,
			userRole:       domain.RoleTeacher,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "parent lacks teacher access",
			requiredRole:   domain.RoleTeacher,
			userRole:       domain.RoleParent,
			expectedStatus: http.StatusForbidden,
			expectNext:     false,
		},
		{
			name:           "admin has admin access",
			requiredRole:   domain.RoleAdmin,
			userRole:       domain.RoleAdmin,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "teacher lacks admin access",
			requiredRole:   domain.RoleAdmin,
			userRole:       domain.RoleTeacher,
			expectedStatus: http.StatusForbidden,
			expectNext:     false,
		},
		{
			name:           "parent has parent access",
			requiredRole:   domain.RoleParent,
			userRole:       domain.RoleParent,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token with user role
			userID := uuid.New().String()
			email := "test@example.com"
			token, err := auth.GenerateAccessToken(userID, email, string(tt.userRole), secret, 15*time.Minute)
			if err != nil {
				t.Fatalf("Failed to generate token: %v", err)
			}

			// Create middleware chain
			cfg := &config.Config{
				JWT: config.JWTConfig{
					Secret: secret,
				},
			}
			authMW := AuthMiddleware(cfg)
			roleMW := RequireRole(string(tt.requiredRole))

			// Create test handler
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Create request with auth token
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute middleware chain
			handler := authMW(roleMW(next))
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Status code = %v, want %v", rr.Code, tt.expectedStatus)
			}

			// Check if next handler was called
			if nextCalled != tt.expectNext {
				t.Errorf("Next handler called = %v, want %v", nextCalled, tt.expectNext)
			}
		})
	}
}

func TestRequireRole_NoAuth(t *testing.T) {
	// Test RequireRole without authentication
	roleMW := RequireRole(string(domain.RoleTeacher))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler := roleMW(next)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Status code = %v, want %v", rr.Code, http.StatusForbidden)
	}
}

func TestGetUserID(t *testing.T) {
	userID := uuid.New()

	// Test with UserID in context
	ctx := context.WithValue(context.Background(), UserIDKey, userID.String())

	gotID, ok := GetUserID(ctx)
	if !ok {
		t.Error("GetUserID() ok = false, want true")
	}
	if gotID != userID {
		t.Errorf("GetUserID() = %v, want %v", gotID, userID)
	}

	// Test without UserID in context
	ctx = context.Background()
	gotID, ok = GetUserID(ctx)
	if ok {
		t.Error("GetUserID() without userID ok = true, want false")
	}
	if gotID != uuid.Nil {
		t.Errorf("GetUserID() without userID = %v, want Nil", gotID)
	}

	// Test with invalid UserID format
	ctx = context.WithValue(context.Background(), UserIDKey, "invalid-uuid")
	gotID, ok = GetUserID(ctx)
	if ok {
		t.Error("GetUserID() with invalid uuid ok = true, want false")
	}
}

func TestGetUserRole(t *testing.T) {
	role := string(domain.RoleTeacher)

	// Test with UserRole in context
	ctx := context.WithValue(context.Background(), UserRoleKey, role)

	gotRole, ok := GetUserRole(ctx)
	if !ok {
		t.Error("GetUserRole() ok = false, want true")
	}
	if gotRole != role {
		t.Errorf("GetUserRole() = %v, want %v", gotRole, role)
	}

	// Test without UserRole in context
	ctx = context.Background()
	gotRole, ok = GetUserRole(ctx)
	if ok {
		t.Error("GetUserRole() without role ok = true, want false")
	}
	if gotRole != "" {
		t.Errorf("GetUserRole() without role = %v, want empty", gotRole)
	}
}

func TestRequestID(t *testing.T) {
	// Test with existing request ID
	existingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	rr := httptest.NewRecorder()

	var capturedReq *http.Request
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
	})

	handler := RequestID(next)
	handler.ServeHTTP(rr, req)

	if capturedReq == nil {
		t.Fatal("Handler was not called")
	}

	gotID, ok := capturedReq.Context().Value(RequestIDKey).(string)
	if !ok {
		t.Error("RequestID not found in context")
	}
	if gotID != existingID {
		t.Errorf("RequestID = %v, want %v", gotID, existingID)
	}
	if rr.Header().Get("X-Request-ID") != existingID {
		t.Errorf("Response header X-Request-ID = %v, want %v", rr.Header().Get("X-Request-ID"), existingID)
	}

	// Test without existing request ID
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	gotID, ok = capturedReq.Context().Value(RequestIDKey).(string)
	if !ok {
		t.Error("RequestID not found in context")
	}
	if gotID == "" {
		t.Error("RequestID should be generated")
	}
	if rr.Header().Get("X-Request-ID") != gotID {
		t.Error("Response header X-Request-ID should match context value")
	}
}

func BenchmarkAuthMiddleware(b *testing.B) {
	secret := "test-secret"
	token, _ := auth.GenerateAccessToken(uuid.New().String(), "test@example.com", "TEACHER", secret, 15*time.Minute)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: secret,
		},
	}
	middleware := AuthMiddleware(cfg)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := middleware(next)
		handler.ServeHTTP(rr, req)
	}
}
