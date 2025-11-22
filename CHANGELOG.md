# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.4] - 2025-11-22

### Changed
- Retrigger release workflow

## [0.1.3] - 2025-11-22

### Added

- basic user crud
- bootstrap
- add unit testing
- add api to docker-compose

### Fixed

- bug in autorelease scripts
- bug in autorelease scripts

### Documentation

- document version and changelog generator

### Other

- fmt
- unit test
- unit test

[0;32m═══════════════════════════════════════════════════════[0m

[1;33mInstructions:[0m
1. Copy the above content
2. Replace [0.1.3] with your actual version (e.g., 1.2.3)
3. Paste it at the top of CHANGELOG.md under the '## [Unreleased]' section
4. Review and edit as needed
5. Remove any irrelevant entries (e.g., WIP commits)

[0;36mTip: You can also pipe to clipboard:[0m
  ./scripts/changelog.sh | pbcopy    # macOS
  ./scripts/changelog.sh | xclip     # Linux

[0;36mOr append directly to CHANGELOG.md:[0m
  ./scripts/changelog.sh >> CHANGELOG_TEMP.md

[0;34mStatistics:[0m
  Total commits: 10
  Features: 4
  Fixes: 2
  Breaking: 0


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
