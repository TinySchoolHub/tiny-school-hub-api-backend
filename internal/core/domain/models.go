package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents user roles in the system
type Role string

const (
	RoleTeacher Role = "TEACHER"
	RoleParent  Role = "PARENT"
	RoleAdmin   Role = "ADMIN"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleTeacher, RoleParent, RoleAdmin:
		return true
	}
	return false
}

// User represents a system user
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never serialize password
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Profile represents user profile information
type Profile struct {
	UserID      uuid.UUID  `json:"user_id"`
	DisplayName string     `json:"display_name"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	ChildName   *string    `json:"child_name,omitempty"`
	ClassID     *uuid.UUID `json:"class_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Class represents a school class
type Class struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Grade     string     `json:"grade"`
	SchoolID  *uuid.UUID `json:"school_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ClassRole represents a role within a class
type ClassRole string

const (
	ClassRoleTeacher ClassRole = "TEACHER"
	ClassRoleParent  ClassRole = "PARENT"
)

// IsValid checks if the class role is valid
func (cr ClassRole) IsValid() bool {
	switch cr {
	case ClassRoleTeacher, ClassRoleParent:
		return true
	}
	return false
}

// ClassMember represents a user's membership in a class
type ClassMember struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ClassID     uuid.UUID `json:"class_id"`
	RoleInClass ClassRole `json:"role_in_class"`
	CreatedAt   time.Time `json:"created_at"`
}

// Photo represents a photo uploaded to a class
type Photo struct {
	ID            uuid.UUID `json:"id"`
	ClassID       uuid.UUID `json:"class_id"`
	UploaderID    uuid.UUID `json:"uploader_id"`
	Caption       *string   `json:"caption,omitempty"`
	MediaKey      string    `json:"media_key"`
	ContentType   string    `json:"content_type"`
	FileSizeBytes int       `json:"file_size_bytes"`
	CreatedAt     time.Time `json:"created_at"`
}

// ReportedBy represents who reported an absence
type ReportedBy string

const (
	ReportedByTeacher ReportedBy = "TEACHER"
	ReportedByParent  ReportedBy = "PARENT"
)

// AbsenceStatus represents the status of an absence
type AbsenceStatus string

const (
	AbsenceStatusPending AbsenceStatus = "PENDING"
	AbsenceStatusAcked   AbsenceStatus = "ACKED"
)

// Absence represents a student absence record
type Absence struct {
	ID           uuid.UUID     `json:"id"`
	StudentName  string        `json:"student_name"`
	ClassID      uuid.UUID     `json:"class_id"`
	AbsenceDate  time.Time     `json:"absence_date"`
	ReportedBy   ReportedBy    `json:"reported_by"`
	ReporterID   uuid.UUID     `json:"reporter_id"`
	Reason       *string       `json:"reason,omitempty"`
	Status       AbsenceStatus `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// Message represents a message between users
type Message struct {
	ID          uuid.UUID  `json:"id"`
	SenderID    uuid.UUID  `json:"sender_id"`
	RecipientID *uuid.UUID `json:"recipient_id,omitempty"`
	ClassID     *uuid.UUID `json:"class_id,omitempty"`
	Body        string     `json:"body"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Announcement represents a class or global announcement
type Announcement struct {
	ID        uuid.UUID  `json:"id"`
	ClassID   *uuid.UUID `json:"class_id,omitempty"`
	AuthorID  uuid.UUID  `json:"author_id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	PublishAt time.Time  `json:"publish_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"` // Never serialize token
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// IsRevoked checks if the token has been revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsExpired checks if the token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the token is valid (not revoked and not expired)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked() && !rt.IsExpired()
}
