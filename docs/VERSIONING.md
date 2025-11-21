# Versioning and Release Process

## Version Format

This project follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., `v1.2.3`)
  - **MAJOR**: Incompatible API changes
  - **MINOR**: New functionality (backward compatible)
  - **PATCH**: Bug fixes (backward compatible)

## Current Version

Current version is tracked in:
- Git tags (e.g., `v1.0.0`)
- `VERSION` file in repository root
- Docker image tags

## Release Process

### 1. Prepare Release

```bash
# Update version in VERSION file
echo "1.2.3" > VERSION

# Update CHANGELOG.md
# Add release notes under new version heading

# Commit changes
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v1.2.3"
```

### 2. Create Git Tag

```bash
# Create annotated tag
git tag -a v1.2.3 -m "Release v1.2.3

- Feature: Added user profile management
- Fix: Resolved authentication timeout issue
- Improvement: Enhanced rate limiting performance
"

# Push tag to remote
git push origin v1.2.3
```

### 3. Automated Release (GitHub Actions)

When you push a tag, GitHub Actions will automatically:
- Run all tests
- Build Docker images
- Create GitHub Release with changelog
- Tag Docker images with version
- Deploy to staging (if configured)

### 4. Manual Release (Alternative)

```bash
# Use the release script
./scripts/release.sh 1.2.3
```

## Release Types

### Patch Release (Bug Fixes)
```bash
./scripts/release.sh patch  # Increments 1.2.3 → 1.2.4
```

### Minor Release (New Features)
```bash
./scripts/release.sh minor  # Increments 1.2.3 → 1.3.0
```

### Major Release (Breaking Changes)
```bash
./scripts/release.sh major  # Increments 1.2.3 → 2.0.0
```

## Pre-release Versions

For beta/alpha releases:

```bash
# Alpha
git tag -a v1.3.0-alpha.1 -m "Alpha release"

# Beta
git tag -a v1.3.0-beta.1 -m "Beta release"

# Release Candidate
git tag -a v1.3.0-rc.1 -m "Release candidate"
```

## Hotfix Process

For urgent production fixes:

```bash
# Create hotfix branch from main
git checkout -b hotfix/v1.2.4 main

# Make fixes
git commit -m "fix: critical security patch"

# Tag and release
git tag -a v1.2.4 -m "Hotfix: Security patch"
git push origin v1.2.4

# Merge back to main
git checkout main
git merge hotfix/v1.2.4
git push origin main
```

## Docker Image Tagging

Docker images are tagged with:
- `latest` - Latest stable release
- `v1.2.3` - Specific version
- `v1.2` - Latest patch version of minor
- `v1` - Latest minor version of major
- `develop` - Development branch builds

## Changelog Guidelines

Follow [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
## [1.2.3] - 2025-11-21

### Added
- User profile management endpoints
- Photo upload with S3 integration

### Changed
- Improved JWT token validation performance
- Updated database connection pooling

### Fixed
- Authentication timeout after 15 minutes
- Rate limiting not working for X-Forwarded-For

### Security
- Fixed SQL injection vulnerability in class queries
```

## Best Practices

1. **Never delete or modify tags** - Tags are immutable references
2. **Always use annotated tags** (`-a` flag) - Include release notes
3. **Test before tagging** - Ensure all tests pass
4. **Update CHANGELOG.md** - Document all changes
5. **Follow semantic versioning** - Be consistent
6. **Use release branches** - For major versions (v1.x, v2.x)
7. **Sign tags** - Use GPG signing for security (`-s` flag)

## Viewing Releases

```bash
# List all tags
git tag -l

# Show tag details
git show v1.2.3

# List releases sorted
git tag -l --sort=-v:refname

# Find latest tag
git describe --tags --abbrev=0
```

## Rollback

If you need to rollback a release:

```bash
# Deploy previous version
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api:v1.2.2

# Or revert in Kubernetes/Helm
helm rollback tiny-school-hub 1
```

## Version in Code

The application exposes version via:
- `/health` endpoint includes version
- Logs include version on startup
- Docker image labels

## Release Checklist

Before creating a release:

- [ ] All tests passing
- [ ] Pre-commit hooks passing
- [ ] CHANGELOG.md updated
- [ ] VERSION file updated
- [ ] Documentation updated
- [ ] Migration files tested
- [ ] Breaking changes documented
- [ ] Security audit completed (for major releases)
- [ ] Performance benchmarks acceptable
- [ ] API documentation up to date

## Support Policy

- **Major versions**: 12 months support
- **Minor versions**: 6 months support
- **Patch versions**: Until next patch

## Questions?

See [CONTRIBUTING.md](../CONTRIBUTING.md) or open an issue.
