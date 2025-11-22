package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/config"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/auth"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/http/middleware"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/repository"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/pkg/log"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
	tokenRepo   repository.RefreshTokenRepository
	cfg         *config.Config
	logger      *log.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	tokenRepo repository.RefreshTokenRepository,
	cfg *config.Config,
	logger *log.Logger,
) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		tokenRepo:   tokenRepo,
		cfg:         cfg,
		logger:      logger,
	}
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		writeError(w, "invalid_input", "Email, password, and display name are required", http.StatusBadRequest)
		return
	}

	// Default to PARENT role if not specified or invalid
	role := domain.Role(req.Role)
	if !role.IsValid() {
		role = domain.RoleParent
	}

	// Check if user already exists
	existing, _ := h.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		writeError(w, "already_exists", "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.WithError(err).Error("Failed to hash password")
		writeError(w, "internal_error", "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		h.logger.WithError(err).Error("Failed to create user")
		writeError(w, "internal_error", "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Create profile
	profile := &domain.Profile{
		UserID:      user.ID,
		DisplayName: req.DisplayName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.profileRepo.Create(ctx, profile); err != nil {
		h.logger.WithError(err).Error("Failed to create profile")
		// Don't fail registration if profile creation fails
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID.String(), user.Email, string(user.Role), h.cfg.JWT.Secret, h.cfg.JWT.AccessExpiry)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate access token")
		writeError(w, "internal_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate refresh token")
		writeError(w, "internal_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Store refresh token
	tokenHash := auth.HashRefreshToken(refreshToken)
	refreshTokenModel := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(h.cfg.JWT.RefreshExpiry),
		CreatedAt: time.Now(),
	}

	if err := h.tokenRepo.Create(ctx, refreshTokenModel); err != nil {
		h.logger.WithError(err).Error("Failed to store refresh token")
		// Don't fail if refresh token storage fails
	}

	writeJSON(w, authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, http.StatusCreated)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		writeError(w, "invalid_credentials", "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	valid, err := auth.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !valid {
		writeError(w, "invalid_credentials", "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID.String(), user.Email, string(user.Role), h.cfg.JWT.Secret, h.cfg.JWT.AccessExpiry)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate access token")
		writeError(w, "internal_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate refresh token")
		writeError(w, "internal_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Store refresh token
	tokenHash := auth.HashRefreshToken(refreshToken)
	refreshTokenModel := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(h.cfg.JWT.RefreshExpiry),
		CreatedAt: time.Now(),
	}

	if err := h.tokenRepo.Create(ctx, refreshTokenModel); err != nil {
		h.logger.WithError(err).Error("Failed to store refresh token")
	}

	writeJSON(w, authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, http.StatusOK)
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenHash := auth.HashRefreshToken(req.RefreshToken)
	storedToken, err := h.tokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		writeError(w, "invalid_token", "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	if !storedToken.IsValid() {
		writeError(w, "invalid_token", "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		writeError(w, "invalid_token", "Invalid user", http.StatusUnauthorized)
		return
	}

	// Revoke old token
	if err := h.tokenRepo.Revoke(ctx, tokenHash, time.Now()); err != nil {
		h.logger.WithError(err).Error("Failed to revoke old token")
	}

	// Generate new tokens
	accessToken, err := auth.GenerateAccessToken(user.ID.String(), user.Email, string(user.Role), h.cfg.JWT.Secret, h.cfg.JWT.AccessExpiry)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate access token")
		writeError(w, "internal_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate refresh token")
		writeError(w, "internal_error", "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Store new refresh token
	newTokenHash := auth.HashRefreshToken(newRefreshToken)
	newTokenModel := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(h.cfg.JWT.RefreshExpiry),
		CreatedAt: time.Now(),
	}

	if err := h.tokenRepo.Create(ctx, newTokenModel); err != nil {
		h.logger.WithError(err).Error("Failed to store refresh token")
	}

	writeJSON(w, authResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, http.StatusOK)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenHash := auth.HashRefreshToken(req.RefreshToken)
	if err := h.tokenRepo.Revoke(ctx, tokenHash, time.Now()); err != nil {
		h.logger.WithError(err).Error("Failed to revoke token")
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return uuid.Nil, domain.ErrUnauthorized
	}
	return userID, nil
}
