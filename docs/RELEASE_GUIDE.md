# Release Quick Reference

## Create a New Release

### Using the Release Script (Recommended)

```bash
# Patch release (bug fixes): 1.2.3 → 1.2.4
./scripts/release.sh patch

# Minor release (new features): 1.2.3 → 1.3.0
./scripts/release.sh minor

# Major release (breaking changes): 1.2.3 → 2.0.0
./scripts/release.sh major

# Specific version
./scripts/release.sh 1.5.0
```

### Manual Release

```bash
# 1. Update VERSION file
echo "1.2.3" > VERSION

# 2. Update CHANGELOG.md
# Add your changes under ## [1.2.3] - YYYY-MM-DD

# 3. Commit
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v1.2.3"

# 4. Tag
git tag -a v1.2.3 -m "Release v1.2.3"

# 5. Push
git push origin main
git push origin v1.2.3
```

## Check Current Version

```bash
# From VERSION file
cat VERSION

# Latest git tag
git describe --tags --abbrev=0

# From running API
curl http://localhost:8080/version
```

## View Releases

```bash
# List all tags
git tag -l

# Show specific tag
git show v1.2.3

# List sorted by version
git tag -l --sort=-v:refname | head -10
```

## Pre-release Versions

```bash
# Alpha
git tag -a v1.3.0-alpha.1 -m "Alpha release"

# Beta  
git tag -a v1.3.0-beta.1 -m "Beta release"

# Release Candidate
git tag -a v1.3.0-rc.1 -m "Release candidate"
```

## What Happens After Tagging?

When you push a tag (e.g., `v1.2.3`), GitHub Actions automatically:

1. ✅ Runs all tests
2. 🐳 Builds Docker image
3. 📦 Publishes to GitHub Container Registry
4. 📝 Creates GitHub Release with changelog
5. 🏷️ Tags image with multiple versions

## Docker Images

After release, images are available:

```bash
# Pull specific version
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2.3

# Pull latest
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:latest

# Pull latest minor
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2

# Pull latest major
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1
```

## Hotfix Process

```bash
# 1. Create hotfix branch
git checkout -b hotfix/v1.2.4 main

# 2. Make fixes
git commit -m "fix: critical security issue"

# 3. Update VERSION and CHANGELOG
echo "1.2.4" > VERSION
# Update CHANGELOG.md

# 4. Commit and tag
git commit -am "chore: bump version to v1.2.4"
git tag -a v1.2.4 -m "Hotfix v1.2.4"

# 5. Push
git push origin hotfix/v1.2.4
git push origin v1.2.4

# 6. Merge back
git checkout main
git merge hotfix/v1.2.4
git push origin main
```

## Rollback a Release

```bash
# Delete local tag
git tag -d v1.2.3

# Delete remote tag
git push origin :refs/tags/v1.2.3

# Or rollback deployment
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2.2
```

## Common Issues

### "Tag already exists"
```bash
# Delete and recreate
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

### "Forgot to update CHANGELOG"
```bash
# Amend the release commit
git commit --amend
git push origin main --force
# Delete and recreate tag
git tag -d v1.2.3
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3 --force
```

### "Tests failed after tagging"
```bash
# Delete the tag
git push origin :refs/tags/v1.2.3
# Fix tests and try again
```

## Versioning Rules

- **MAJOR**: Breaking API changes (e.g., 1.x.x → 2.0.0)
- **MINOR**: New features, backward compatible (e.g., 1.2.x → 1.3.0)
- **PATCH**: Bug fixes, backward compatible (e.g., 1.2.3 → 1.2.4)

## Release Checklist

Before running `./scripts/release.sh`:

- [ ] All tests passing
- [ ] Pre-commit hooks passing
- [ ] CHANGELOG.md updated
- [ ] Documentation updated
- [ ] No uncommitted changes
- [ ] On main branch
- [ ] Reviewed all changes

## Help

- 📖 Full documentation: [docs/VERSIONING.md](docs/VERSIONING.md)
- 🤝 Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
- 📝 Changelog format: [keepachangelog.com](https://keepachangelog.com/)
- 🏷️ Semantic versioning: [semver.org](https://semver.org/)
