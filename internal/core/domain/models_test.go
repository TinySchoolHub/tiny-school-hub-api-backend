package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserRole(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want string
	}{
		{"teacher role", RoleTeacher, "TEACHER"},
		{"parent role", RoleParent, "PARENT"},
		{"admin role", RoleAdmin, "ADMIN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.want {
				t.Errorf("Role = %v, want %v", tt.role, tt.want)
			}
		})
	}
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid teacher",
			user: User{
				ID:           uuid.New(),
				Email:        "teacher@school.com",
				PasswordHash: "hashed_password",
				Role:         RoleTeacher,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid parent",
			user: User{
				ID:           uuid.New(),
				Email:        "parent@example.com",
				PasswordHash: "hashed_password",
				Role:         RoleParent,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid admin",
			user: User{
				ID:           uuid.New(),
				Email:        "admin@school.com",
				PasswordHash: "hashed_password",
				Role:         RoleAdmin,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user.ID == uuid.Nil {
				t.Error("User ID should not be nil")
			}
			if tt.user.Email == "" {
				t.Error("User email should not be empty")
			}
			if tt.user.PasswordHash == "" {
				t.Error("User password hash should not be empty")
			}
			if tt.user.CreatedAt.IsZero() {
				t.Error("User CreatedAt should not be zero")
			}
		})
	}
}

func TestClassValidation(t *testing.T) {
	tests := []struct {
		name  string
		class Class
		valid bool
	}{
		{
			name: "valid class",
			class: Class{
				ID:        uuid.New(),
				Name:      "Math 101",
				Grade:     "5th Grade",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			valid: true,
		},
		{
			name: "class with school ID",
			class: Class{
				ID:        uuid.New(),
				Name:      "Science 202",
				Grade:     "6th Grade",
				SchoolID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.class.ID == uuid.Nil {
				t.Error("Class ID should not be nil")
			}
			if tt.class.Name == "" {
				t.Error("Class name should not be empty")
			}
			if tt.class.Grade == "" {
				t.Error("Class grade should not be empty")
			}
		})
	}
}

func TestClassMemberRoles(t *testing.T) {
	tests := []struct {
		name string
		role ClassRole
		want bool
	}{
		{"teacher role", ClassRoleTeacher, true},
		{"parent role", ClassRoleParent, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := ClassMember{
				ID:          uuid.New(),
				UserID:      uuid.New(),
				ClassID:     uuid.New(),
				RoleInClass: tt.role,
				CreatedAt:   time.Now(),
			}

			isValid := member.RoleInClass.IsValid()
			if isValid != tt.want {
				t.Errorf("Role validation = %v, want %v for role %s", isValid, tt.want, tt.role)
			}
		})
	}
}

func TestPhotoValidation(t *testing.T) {
	tests := []struct {
		name  string
		photo Photo
		valid bool
	}{
		{
			name: "valid photo with caption",
			photo: Photo{
				ID:            uuid.New(),
				ClassID:       uuid.New(),
				UploaderID:    uuid.New(),
				Caption:       func() *string { s := "School event"; return &s }(),
				MediaKey:      "photos/class123/photo456.jpg",
				ContentType:   "image/jpeg",
				FileSizeBytes: 1024000,
				CreatedAt:     time.Now(),
			},
			valid: true,
		},
		{
			name: "valid photo without caption",
			photo: Photo{
				ID:            uuid.New(),
				ClassID:       uuid.New(),
				UploaderID:    uuid.New(),
				Caption:       nil,
				MediaKey:      "photos/class456/photo789.png",
				ContentType:   "image/png",
				FileSizeBytes: 512000,
				CreatedAt:     time.Now(),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.photo.ID == uuid.Nil {
				t.Error("Photo ID should not be nil")
			}
			if tt.photo.MediaKey == "" {
				t.Error("Photo media key should not be empty")
			}
			if tt.photo.ContentType == "" {
				t.Error("Photo content type should not be empty")
			}
			if tt.photo.FileSizeBytes <= 0 {
				t.Error("Photo file size should be positive")
			}
		})
	}
}

func TestAbsenceValidation(t *testing.T) {
	absenceDate := time.Now()

	tests := []struct {
		name    string
		absence Absence
		valid   bool
	}{
		{
			name: "valid absence with reason",
			absence: Absence{
				ID:          uuid.New(),
				StudentName: "John Doe",
				ClassID:     uuid.New(),
				AbsenceDate: absenceDate,
				ReportedBy:  ReportedByTeacher,
				ReporterID:  uuid.New(),
				Reason:      func() *string { s := "Sick leave"; return &s }(),
				Status:      AbsenceStatusPending,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			valid: true,
		},
		{
			name: "valid absence without reason",
			absence: Absence{
				ID:          uuid.New(),
				StudentName: "Jane Smith",
				ClassID:     uuid.New(),
				AbsenceDate: absenceDate,
				ReportedBy:  ReportedByParent,
				ReporterID:  uuid.New(),
				Reason:      nil,
				Status:      AbsenceStatusAcked,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.absence.ID == uuid.Nil {
				t.Error("Absence ID should not be nil")
			}
			if tt.absence.StudentName == "" {
				t.Error("Absence student name should not be empty")
			}
			if tt.absence.AbsenceDate.IsZero() {
				t.Error("Absence date should not be zero")
			}
			if tt.absence.ReporterID == uuid.Nil {
				t.Error("Reporter ID should not be nil")
			}
		})
	}
}

func TestMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		valid   bool
	}{
		{
			name: "valid direct message",
			message: Message{
				ID:          uuid.New(),
				SenderID:    uuid.New(),
				RecipientID: func() *uuid.UUID { id := uuid.New(); return &id }(),
				ClassID:     nil,
				Body:        "Hello, this is a test message",
				ReadAt:      nil,
				CreatedAt:   time.Now(),
			},
			valid: true,
		},
		{
			name: "read class message",
			message: Message{
				ID:          uuid.New(),
				SenderID:    uuid.New(),
				RecipientID: nil,
				ClassID:     func() *uuid.UUID { id := uuid.New(); return &id }(),
				Body:        "Class announcement message",
				ReadAt:      func() *time.Time { t := time.Now(); return &t }(),
				CreatedAt:   time.Now().Add(-1 * time.Hour),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.message.ID == uuid.Nil {
				t.Error("Message ID should not be nil")
			}
			if tt.message.Body == "" {
				t.Error("Message body should not be empty")
			}
			if tt.message.SenderID == uuid.Nil {
				t.Error("Message sender ID should not be nil")
			}
		})
	}
}

func TestAnnouncementValidation(t *testing.T) {
	tests := []struct {
		name         string
		announcement Announcement
		valid        bool
	}{
		{
			name: "class announcement",
			announcement: Announcement{
				ID:        uuid.New(),
				AuthorID:  uuid.New(),
				ClassID:   func() *uuid.UUID { id := uuid.New(); return &id }(),
				Title:     "School Event",
				Body:      "We will have a school event next week",
				PublishAt: time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			valid: true,
		},
		{
			name: "global announcement",
			announcement: Announcement{
				ID:        uuid.New(),
				AuthorID:  uuid.New(),
				ClassID:   nil,
				Title:     "School Closure",
				Body:      "School will be closed tomorrow",
				PublishAt: time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.announcement.ID == uuid.Nil {
				t.Error("Announcement ID should not be nil")
			}
			if tt.announcement.Title == "" {
				t.Error("Announcement title should not be empty")
			}
			if tt.announcement.Body == "" {
				t.Error("Announcement body should not be empty")
			}
			if tt.announcement.AuthorID == uuid.Nil {
				t.Error("Announcement author ID should not be nil")
			}
		})
	}
}

func TestRefreshTokenValidation(t *testing.T) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	tests := []struct {
		name  string
		token RefreshToken
		valid bool
	}{
		{
			name: "valid active token",
			token: RefreshToken{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				TokenHash: "hashed_token_value",
				ExpiresAt: expiresAt,
				RevokedAt: nil,
				CreatedAt: time.Now(),
			},
			valid: true,
		},
		{
			name: "revoked token",
			token: RefreshToken{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				TokenHash: "hashed_token_value",
				ExpiresAt: expiresAt,
				RevokedAt: func() *time.Time { t := time.Now(); return &t }(),
				CreatedAt: time.Now(),
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.token.ID == uuid.Nil {
				t.Error("Token ID should not be nil")
			}
			if tt.token.TokenHash == "" {
				t.Error("Token hash should not be empty")
			}
			if tt.token.ExpiresAt.IsZero() {
				t.Error("Token expiry should not be zero")
			}

			isValid := tt.token.IsValid()
			if isValid != tt.valid {
				t.Errorf("Token validity = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestProfileValidation(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		valid   bool
	}{
		{
			name: "valid profile with all fields",
			profile: Profile{
				UserID:      uuid.New(),
				DisplayName: "John Doe",
				AvatarURL:   func() *string { s := "https://example.com/avatar.jpg"; return &s }(),
				ChildName:   func() *string { s := "Little John"; return &s }(),
				ClassID:     func() *uuid.UUID { id := uuid.New(); return &id }(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			valid: true,
		},
		{
			name: "minimal valid profile",
			profile: Profile{
				UserID:      uuid.New(),
				DisplayName: "Jane Smith",
				AvatarURL:   nil,
				ChildName:   nil,
				ClassID:     nil,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profile.UserID == uuid.Nil {
				t.Error("Profile user ID should not be nil")
			}
			if tt.profile.DisplayName == "" {
				t.Error("Profile display name should not be empty")
			}
		})
	}
}
