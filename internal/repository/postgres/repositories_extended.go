package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/repository"
)

// PhotoRepo implements repository.PhotoRepository
type PhotoRepo struct {
	db *DB
}

func NewPhotoRepo(db *DB) repository.PhotoRepository {
	return &PhotoRepo{db: db}
}

func (r *PhotoRepo) Create(ctx context.Context, photo *domain.Photo) error {
	query := `INSERT INTO photos (id, class_id, uploader_id, caption, media_key, content_type, file_size_bytes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, photo.ID, photo.ClassID, photo.UploaderID, photo.Caption, photo.MediaKey, photo.ContentType, photo.FileSizeBytes, photo.CreatedAt)
	return err
}

func (r *PhotoRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Photo, error) {
	query := `SELECT id, class_id, uploader_id, caption, media_key, content_type, file_size_bytes, created_at FROM photos WHERE id = $1`
	photo := &domain.Photo{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&photo.ID, &photo.ClassID, &photo.UploaderID, &photo.Caption, &photo.MediaKey, &photo.ContentType, &photo.FileSizeBytes, &photo.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return photo, err
}

func (r *PhotoRepo) ListByClass(ctx context.Context, classID uuid.UUID, limit, offset int) ([]*domain.Photo, error) {
	query := `SELECT id, class_id, uploader_id, caption, media_key, content_type, file_size_bytes, created_at 
		FROM photos WHERE class_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, classID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*domain.Photo
	for rows.Next() {
		photo := &domain.Photo{}
		if err := rows.Scan(&photo.ID, &photo.ClassID, &photo.UploaderID, &photo.Caption, &photo.MediaKey, &photo.ContentType, &photo.FileSizeBytes, &photo.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, rows.Err()
}

func (r *PhotoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM photos WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// AbsenceRepo implements repository.AbsenceRepository
type AbsenceRepo struct {
	db *DB
}

func NewAbsenceRepo(db *DB) repository.AbsenceRepository {
	return &AbsenceRepo{db: db}
}

func (r *AbsenceRepo) Create(ctx context.Context, absence *domain.Absence) error {
	query := `INSERT INTO absences (id, student_name, class_id, absence_date, reported_by, reporter_id, reason, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, absence.ID, absence.StudentName, absence.ClassID, absence.AbsenceDate, absence.ReportedBy, absence.ReporterID, absence.Reason, absence.Status, absence.CreatedAt, absence.UpdatedAt)
	return err
}

func (r *AbsenceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Absence, error) {
	query := `SELECT id, student_name, class_id, absence_date, reported_by, reporter_id, reason, status, created_at, updated_at FROM absences WHERE id = $1`
	absence := &domain.Absence{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&absence.ID, &absence.StudentName, &absence.ClassID, &absence.AbsenceDate, &absence.ReportedBy, &absence.ReporterID, &absence.Reason, &absence.Status, &absence.CreatedAt, &absence.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return absence, err
}

func (r *AbsenceRepo) ListByClass(ctx context.Context, classID uuid.UUID, limit, offset int) ([]*domain.Absence, error) {
	query := `SELECT id, student_name, class_id, absence_date, reported_by, reporter_id, reason, status, created_at, updated_at 
		FROM absences WHERE class_id = $1 ORDER BY absence_date DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, classID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var absences []*domain.Absence
	for rows.Next() {
		absence := &domain.Absence{}
		if err := rows.Scan(&absence.ID, &absence.StudentName, &absence.ClassID, &absence.AbsenceDate, &absence.ReportedBy, &absence.ReporterID, &absence.Reason, &absence.Status, &absence.CreatedAt, &absence.UpdatedAt); err != nil {
			return nil, err
		}
		absences = append(absences, absence)
	}
	return absences, rows.Err()
}

func (r *AbsenceRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AbsenceStatus) error {
	query := `UPDATE absences SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *AbsenceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM absences WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MessageRepo implements repository.MessageRepository
type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) repository.MessageRepository {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Create(ctx context.Context, message *domain.Message) error {
	query := `INSERT INTO messages (id, sender_id, recipient_id, class_id, body, read_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, message.ID, message.SenderID, message.RecipientID, message.ClassID, message.Body, message.ReadAt, message.CreatedAt)
	return err
}

func (r *MessageRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	query := `SELECT id, sender_id, recipient_id, class_id, body, read_at, created_at FROM messages WHERE id = $1`
	message := &domain.Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&message.ID, &message.SenderID, &message.RecipientID, &message.ClassID, &message.Body, &message.ReadAt, &message.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return message, err
}

func (r *MessageRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Message, error) {
	query := `SELECT id, sender_id, recipient_id, class_id, body, read_at, created_at 
		FROM messages WHERE recipient_id = $1 OR sender_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		message := &domain.Message{}
		if err := rows.Scan(&message.ID, &message.SenderID, &message.RecipientID, &message.ClassID, &message.Body, &message.ReadAt, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *MessageRepo) MarkAsRead(ctx context.Context, id uuid.UUID, readAt time.Time) error {
	query := `UPDATE messages SET read_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, readAt, id)
	return err
}

func (r *MessageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// AnnouncementRepo implements repository.AnnouncementRepository
type AnnouncementRepo struct {
	db *DB
}

func NewAnnouncementRepo(db *DB) repository.AnnouncementRepository {
	return &AnnouncementRepo{db: db}
}

func (r *AnnouncementRepo) Create(ctx context.Context, announcement *domain.Announcement) error {
	query := `INSERT INTO announcements (id, class_id, author_id, title, body, publish_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, announcement.ID, announcement.ClassID, announcement.AuthorID, announcement.Title, announcement.Body, announcement.PublishAt, announcement.CreatedAt, announcement.UpdatedAt)
	return err
}

func (r *AnnouncementRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Announcement, error) {
	query := `SELECT id, class_id, author_id, title, body, publish_at, created_at, updated_at FROM announcements WHERE id = $1`
	announcement := &domain.Announcement{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&announcement.ID, &announcement.ClassID, &announcement.AuthorID, &announcement.Title, &announcement.Body, &announcement.PublishAt, &announcement.CreatedAt, &announcement.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return announcement, err
}

func (r *AnnouncementRepo) ListByClass(ctx context.Context, classID *uuid.UUID, limit, offset int) ([]*domain.Announcement, error) {
	query := `SELECT id, class_id, author_id, title, body, publish_at, created_at, updated_at 
		FROM announcements WHERE class_id = $1 AND publish_at <= NOW() ORDER BY publish_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, classID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []*domain.Announcement
	for rows.Next() {
		announcement := &domain.Announcement{}
		if err := rows.Scan(&announcement.ID, &announcement.ClassID, &announcement.AuthorID, &announcement.Title, &announcement.Body, &announcement.PublishAt, &announcement.CreatedAt, &announcement.UpdatedAt); err != nil {
			return nil, err
		}
		announcements = append(announcements, announcement)
	}
	return announcements, rows.Err()
}

func (r *AnnouncementRepo) ListGlobal(ctx context.Context, limit, offset int) ([]*domain.Announcement, error) {
	query := `SELECT id, class_id, author_id, title, body, publish_at, created_at, updated_at 
		FROM announcements WHERE class_id IS NULL AND publish_at <= NOW() ORDER BY publish_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []*domain.Announcement
	for rows.Next() {
		announcement := &domain.Announcement{}
		if err := rows.Scan(&announcement.ID, &announcement.ClassID, &announcement.AuthorID, &announcement.Title, &announcement.Body, &announcement.PublishAt, &announcement.CreatedAt, &announcement.UpdatedAt); err != nil {
			return nil, err
		}
		announcements = append(announcements, announcement)
	}
	return announcements, rows.Err()
}

func (r *AnnouncementRepo) Update(ctx context.Context, announcement *domain.Announcement) error {
	query := `UPDATE announcements SET title = $1, body = $2, publish_at = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, announcement.Title, announcement.Body, announcement.PublishAt, announcement.UpdatedAt, announcement.ID)
	return err
}

func (r *AnnouncementRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM announcements WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// RefreshTokenRepo implements repository.RefreshTokenRepository
type RefreshTokenRepo struct {
	db *DB
}

func NewRefreshTokenRepo(db *DB) repository.RefreshTokenRepository {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.RevokedAt, token.CreatedAt)
	return err
}

func (r *RefreshTokenRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE token_hash = $1`
	token := &domain.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.RevokedAt, &token.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return token, err
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2`
	_, err := r.db.ExecContext(ctx, query, revokedAt, tokenHash)
	return err
}

func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, revokedAt, userID)
	return err
}

func (r *RefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
