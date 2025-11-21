package repository

import (
	"context"
	"time"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProfileRepository defines the interface for profile persistence
type ProfileRepository interface {
	Create(ctx context.Context, profile *domain.Profile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error)
	Update(ctx context.Context, profile *domain.Profile) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

// ClassRepository defines the interface for class persistence
type ClassRepository interface {
	Create(ctx context.Context, class *domain.Class) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Class, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Class, error)
	Update(ctx context.Context, class *domain.Class) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUserClasses(ctx context.Context, userID uuid.UUID) ([]*domain.Class, error)
}

// ClassMemberRepository defines the interface for class membership persistence
type ClassMemberRepository interface {
	Create(ctx context.Context, member *domain.ClassMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ClassMember, error)
	GetByUserAndClass(ctx context.Context, userID, classID uuid.UUID) (*domain.ClassMember, error)
	ListByClass(ctx context.Context, classID uuid.UUID) ([]*domain.ClassMember, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.ClassMember, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IsMember(ctx context.Context, userID, classID uuid.UUID) (bool, error)
	IsTeacher(ctx context.Context, userID, classID uuid.UUID) (bool, error)
}

// PhotoRepository defines the interface for photo persistence
type PhotoRepository interface {
	Create(ctx context.Context, photo *domain.Photo) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Photo, error)
	ListByClass(ctx context.Context, classID uuid.UUID, limit, offset int) ([]*domain.Photo, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// AbsenceRepository defines the interface for absence persistence
type AbsenceRepository interface {
	Create(ctx context.Context, absence *domain.Absence) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Absence, error)
	ListByClass(ctx context.Context, classID uuid.UUID, limit, offset int) ([]*domain.Absence, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AbsenceStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// MessageRepository defines the interface for message persistence
type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Message, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, readAt time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AnnouncementRepository defines the interface for announcement persistence
type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *domain.Announcement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Announcement, error)
	ListByClass(ctx context.Context, classID *uuid.UUID, limit, offset int) ([]*domain.Announcement, error)
	ListGlobal(ctx context.Context, limit, offset int) ([]*domain.Announcement, error)
	Update(ctx context.Context, announcement *domain.Announcement) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// RefreshTokenRepository defines the interface for refresh token persistence
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error
	DeleteExpired(ctx context.Context) error
}
