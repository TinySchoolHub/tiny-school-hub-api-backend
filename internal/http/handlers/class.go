package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/config"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/repository"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/storage"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/pkg/log"
)

// ClassHandler handles class endpoints
type ClassHandler struct {
	classRepo  repository.ClassRepository
	memberRepo repository.ClassMemberRepository
	cfg        *config.Config
	logger     *log.Logger
}

func NewClassHandler(
	classRepo repository.ClassRepository,
	memberRepo repository.ClassMemberRepository,
	cfg *config.Config,
	logger *log.Logger,
) *ClassHandler {
	return &ClassHandler{classRepo: classRepo, memberRepo: memberRepo, cfg: cfg, logger: logger}
}

type createClassRequest struct {
	Name     string     `json:"name"`
	Grade    string     `json:"grade"`
	SchoolID *uuid.UUID `json:"school_id,omitempty"`
}

func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		writeError(w, "unauthorized", "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req createClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	class := &domain.Class{
		ID:        uuid.New(),
		Name:      req.Name,
		Grade:     req.Grade,
		SchoolID:  req.SchoolID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.classRepo.Create(ctx, class); err != nil {
		h.logger.WithError(err).Error("Failed to create class")
		writeError(w, "internal_error", "Failed to create class", http.StatusInternalServerError)
		return
	}

	// Add creator as teacher
	member := &domain.ClassMember{
		ID:          uuid.New(),
		UserID:      userID,
		ClassID:     class.ID,
		RoleInClass: domain.ClassRoleTeacher,
		CreatedAt:   time.Now(),
	}

	if err := h.memberRepo.Create(ctx, member); err != nil {
		h.logger.WithError(err).Error("Failed to add class member")
	}

	writeJSON(w, class, http.StatusCreated)
}

func (h *ClassHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		writeError(w, "unauthorized", "Unauthorized", http.StatusUnauthorized)
		return
	}

	classID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid_request", "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Check membership
	isMember, err := h.memberRepo.IsMember(ctx, userID, classID)
	if err != nil || !isMember {
		writeError(w, "forbidden", "Not a member of this class", http.StatusForbidden)
		return
	}

	class, err := h.classRepo.GetByID(ctx, classID)
	if err != nil {
		writeError(w, "not_found", "Class not found", http.StatusNotFound)
		return
	}

	writeJSON(w, class, http.StatusOK)
}

func (h *ClassHandler) ListMyClasses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		writeError(w, "unauthorized", "Unauthorized", http.StatusUnauthorized)
		return
	}

	classes, err := h.classRepo.GetUserClasses(ctx, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list classes")
		writeError(w, "internal_error", "Failed to list classes", http.StatusInternalServerError)
		return
	}

	writeJSON(w, classes, http.StatusOK)
}

func (h *ClassHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		writeError(w, "unauthorized", "Unauthorized", http.StatusUnauthorized)
		return
	}

	classID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid_request", "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Must be teacher to view members
	isTeacher, err := h.memberRepo.IsTeacher(ctx, userID, classID)
	if err != nil || !isTeacher {
		writeError(w, "forbidden", "Must be a teacher to view members", http.StatusForbidden)
		return
	}

	members, err := h.memberRepo.ListByClass(ctx, classID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list members")
		writeError(w, "internal_error", "Failed to list members", http.StatusInternalServerError)
		return
	}

	writeJSON(w, members, http.StatusOK)
}

// PhotoHandler handles photo endpoints
type PhotoHandler struct {
	photoRepo  repository.PhotoRepository
	memberRepo repository.ClassMemberRepository
	storage    *storage.Client
	cfg        *config.Config
	logger     *log.Logger
}

func NewPhotoHandler(
	photoRepo repository.PhotoRepository,
	memberRepo repository.ClassMemberRepository,
	storage *storage.Client,
	cfg *config.Config,
	logger *log.Logger,
) *PhotoHandler {
	return &PhotoHandler{photoRepo: photoRepo, memberRepo: memberRepo, storage: storage, cfg: cfg, logger: logger}
}

type createPhotoRequest struct {
	Caption     *string `json:"caption"`
	ContentType string  `json:"content_type"`
	FileSize    int     `json:"file_size"`
}

type photoUploadResponse struct {
	PhotoID   uuid.UUID `json:"photo_id"`
	UploadURL string    `json:"upload_url"`
	MediaKey  string    `json:"media_key"`
}

func (h *PhotoHandler) CreateUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		writeError(w, "unauthorized", "Unauthorized", http.StatusUnauthorized)
		return
	}

	classID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid_request", "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Must be teacher to upload
	isTeacher, err := h.memberRepo.IsTeacher(ctx, userID, classID)
	if err != nil || !isTeacher {
		writeError(w, "forbidden", "Must be a teacher to upload photos", http.StatusForbidden)
		return
	}

	var req createPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate
	if err := storage.ValidateContentType(req.ContentType); err != nil {
		writeError(w, "invalid_file_type", "Invalid file type", http.StatusBadRequest)
		return
	}

	if err := storage.ValidateFileSize(req.FileSize); err != nil {
		writeError(w, "file_too_large", "File too large (max 5MB)", http.StatusBadRequest)
		return
	}

	photoID := uuid.New()
	mediaKey := "photos/" + classID.String() + "/" + photoID.String()

	// Generate presigned URL
	uploadURL, err := h.storage.GeneratePresignedPutURL(ctx, mediaKey, req.ContentType)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate presigned URL")
		writeError(w, "internal_error", "Failed to generate upload URL", http.StatusInternalServerError)
		return
	}

	// Save photo metadata
	photo := &domain.Photo{
		ID:            photoID,
		ClassID:       classID,
		UploaderID:    userID,
		Caption:       req.Caption,
		MediaKey:      mediaKey,
		ContentType:   req.ContentType,
		FileSizeBytes: req.FileSize,
		CreatedAt:     time.Now(),
	}

	if err := h.photoRepo.Create(ctx, photo); err != nil {
		h.logger.WithError(err).Error("Failed to create photo")
		writeError(w, "internal_error", "Failed to create photo", http.StatusInternalServerError)
		return
	}

	writeJSON(w, photoUploadResponse{
		PhotoID:   photoID,
		UploadURL: uploadURL,
		MediaKey:  mediaKey,
	}, http.StatusCreated)
}

type photoResponse struct {
	*domain.Photo
	ViewURL string `json:"view_url"`
}

func (h *PhotoHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		writeError(w, "unauthorized", "Unauthorized", http.StatusUnauthorized)
		return
	}

	classID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid_request", "Invalid class ID", http.StatusBadRequest)
		return
	}

	// Must be member
	isMember, err := h.memberRepo.IsMember(ctx, userID, classID)
	if err != nil || !isMember {
		writeError(w, "forbidden", "Not a member of this class", http.StatusForbidden)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	photos, err := h.photoRepo.ListByClass(ctx, classID, limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list photos")
		writeError(w, "internal_error", "Failed to list photos", http.StatusInternalServerError)
		return
	}

	// Add presigned URLs
	response := make([]photoResponse, len(photos))
	for i, photo := range photos {
		viewURL, _ := h.storage.GeneratePresignedGetURL(ctx, photo.MediaKey)
		response[i] = photoResponse{
			Photo:   photo,
			ViewURL: viewURL,
		}
	}

	writeJSON(w, response, http.StatusOK)
}
