package domain

import "errors"

// Common domain errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden - insufficient permissions")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrNotAMember         = errors.New("not a member of this class")
	ErrInvalidFileType    = errors.New("invalid file type")
	ErrFileTooLarge       = errors.New("file too large")
)
