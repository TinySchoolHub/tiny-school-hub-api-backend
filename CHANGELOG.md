# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup with Go backend
- User authentication with JWT tokens
- Class management endpoints
- Photo upload with S3/MinIO integration
- Profile management for parents
- Message system between teachers and parents
- Announcement system
- Absence tracking
- Comprehensive unit tests (85%+ coverage)
- Docker and docker-compose setup
- Swagger UI documentation
- Pre-commit hooks for code quality
- Rate limiting middleware
- CORS support
- PostgreSQL database with migrations

### Security
- Argon2id password hashing
- JWT token validation
- Rate limiting per IP
- Sensitive data detection in pre-commit

## [0.1.0] - 2025-11-21

### Added
- Initial alpha release
- Core API functionality
- Database schema
- Authentication system
- Basic CRUD operations

---

## Release Template

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security fixes or improvements
```

[Unreleased]: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/releases/tag/v0.1.0
