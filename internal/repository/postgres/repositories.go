package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/core/domain"
	"github.com/TinySchoolHub/tiny-school-hub-api-backend/internal/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// DB wraps sql.DB with additional methods
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// UserRepo implements repository.UserRepository
type UserRepo struct {
	db *DB
}

// NewUserRepo creates a new user repository
func NewUserRepo(db *DB) repository.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2, role = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.Role, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ProfileRepo implements repository.ProfileRepository
type ProfileRepo struct {
	db *DB
}

func NewProfileRepo(db *DB) repository.ProfileRepository {
	return &ProfileRepo{db: db}
}

func (r *ProfileRepo) Create(ctx context.Context, profile *domain.Profile) error {
	query := `INSERT INTO profiles (user_id, display_name, avatar_url, child_name, class_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, profile.UserID, profile.DisplayName, profile.AvatarURL, profile.ChildName, profile.ClassID, profile.CreatedAt, profile.UpdatedAt)
	return err
}

func (r *ProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	query := `SELECT user_id, display_name, avatar_url, child_name, class_id, created_at, updated_at FROM profiles WHERE user_id = $1`
	profile := &domain.Profile{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&profile.UserID, &profile.DisplayName, &profile.AvatarURL, &profile.ChildName, &profile.ClassID, &profile.CreatedAt, &profile.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return profile, err
}

func (r *ProfileRepo) Update(ctx context.Context, profile *domain.Profile) error {
	query := `UPDATE profiles SET display_name = $1, avatar_url = $2, child_name = $3, class_id = $4, updated_at = $5 WHERE user_id = $6`
	_, err := r.db.ExecContext(ctx, query, profile.DisplayName, profile.AvatarURL, profile.ChildName, profile.ClassID, profile.UpdatedAt, profile.UserID)
	return err
}

func (r *ProfileRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM profiles WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// ClassRepo implements repository.ClassRepository
type ClassRepo struct {
	db *DB
}

func NewClassRepo(db *DB) repository.ClassRepository {
	return &ClassRepo{db: db}
}

func (r *ClassRepo) Create(ctx context.Context, class *domain.Class) error {
	query := `INSERT INTO classes (id, name, grade, school_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, class.ID, class.Name, class.Grade, class.SchoolID, class.CreatedAt, class.UpdatedAt)
	return err
}

func (r *ClassRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Class, error) {
	query := `SELECT id, name, grade, school_id, created_at, updated_at FROM classes WHERE id = $1`
	class := &domain.Class{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&class.ID, &class.Name, &class.Grade, &class.SchoolID, &class.CreatedAt, &class.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return class, err
}

func (r *ClassRepo) List(ctx context.Context, limit, offset int) ([]*domain.Class, error) {
	query := `SELECT id, name, grade, school_id, created_at, updated_at FROM classes ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*domain.Class
	for rows.Next() {
		class := &domain.Class{}
		if err := rows.Scan(&class.ID, &class.Name, &class.Grade, &class.SchoolID, &class.CreatedAt, &class.UpdatedAt); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, rows.Err()
}

func (r *ClassRepo) GetUserClasses(ctx context.Context, userID uuid.UUID) ([]*domain.Class, error) {
	query := `SELECT c.id, c.name, c.grade, c.school_id, c.created_at, c.updated_at 
		FROM classes c INNER JOIN class_members cm ON c.id = cm.class_id WHERE cm.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*domain.Class
	for rows.Next() {
		class := &domain.Class{}
		if err := rows.Scan(&class.ID, &class.Name, &class.Grade, &class.SchoolID, &class.CreatedAt, &class.UpdatedAt); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}
	return classes, rows.Err()
}

func (r *ClassRepo) Update(ctx context.Context, class *domain.Class) error {
	query := `UPDATE classes SET name = $1, grade = $2, school_id = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, class.Name, class.Grade, class.SchoolID, class.UpdatedAt, class.ID)
	return err
}

func (r *ClassRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM classes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ClassMemberRepo implements repository.ClassMemberRepository
type ClassMemberRepo struct {
	db *DB
}

func NewClassMemberRepo(db *DB) repository.ClassMemberRepository {
	return &ClassMemberRepo{db: db}
}

func (r *ClassMemberRepo) Create(ctx context.Context, member *domain.ClassMember) error {
	query := `INSERT INTO class_members (id, user_id, class_id, role_in_class, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, member.ID, member.UserID, member.ClassID, member.RoleInClass, member.CreatedAt)
	return err
}

func (r *ClassMemberRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ClassMember, error) {
	query := `SELECT id, user_id, class_id, role_in_class, created_at FROM class_members WHERE id = $1`
	member := &domain.ClassMember{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&member.ID, &member.UserID, &member.ClassID, &member.RoleInClass, &member.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return member, err
}

func (r *ClassMemberRepo) GetByUserAndClass(ctx context.Context, userID, classID uuid.UUID) (*domain.ClassMember, error) {
	query := `SELECT id, user_id, class_id, role_in_class, created_at FROM class_members WHERE user_id = $1 AND class_id = $2`
	member := &domain.ClassMember{}
	err := r.db.QueryRowContext(ctx, query, userID, classID).Scan(&member.ID, &member.UserID, &member.ClassID, &member.RoleInClass, &member.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return member, err
}

func (r *ClassMemberRepo) ListByClass(ctx context.Context, classID uuid.UUID) ([]*domain.ClassMember, error) {
	query := `SELECT id, user_id, class_id, role_in_class, created_at FROM class_members WHERE class_id = $1`
	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.ClassMember
	for rows.Next() {
		member := &domain.ClassMember{}
		if err := rows.Scan(&member.ID, &member.UserID, &member.ClassID, &member.RoleInClass, &member.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *ClassMemberRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.ClassMember, error) {
	query := `SELECT id, user_id, class_id, role_in_class, created_at FROM class_members WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.ClassMember
	for rows.Next() {
		member := &domain.ClassMember{}
		if err := rows.Scan(&member.ID, &member.UserID, &member.ClassID, &member.RoleInClass, &member.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *ClassMemberRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM class_members WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ClassMemberRepo) IsMember(ctx context.Context, userID, classID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM class_members WHERE user_id = $1 AND class_id = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, classID).Scan(&exists)
	return exists, err
}

func (r *ClassMemberRepo) IsTeacher(ctx context.Context, userID, classID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM class_members WHERE user_id = $1 AND class_id = $2 AND role_in_class = 'TEACHER')`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, classID).Scan(&exists)
	return exists, err
}
