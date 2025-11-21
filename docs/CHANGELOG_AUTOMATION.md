# Automated Changelog Generation

## Overview

The changelog is automatically generated from your git commit messages using **Conventional Commits** format.

## How It Works

The `scripts/changelog.sh` script analyzes your git commits and automatically categorizes them into changelog sections based on the commit type.

## Usage

### Option 1: Automatic (Recommended)

The release script now automatically generates and inserts changelog entries:

```bash
./scripts/release.sh patch
```

When you run the release script, it will:
1. ✅ Generate changelog from commits since last tag
2. ✅ Show you the generated content
3. ✅ Ask if you want to automatically insert it
4. ✅ Insert it into CHANGELOG.md if you confirm

### Option 2: Manual Generation

Generate changelog without releasing:

```bash
# Changes since last tag
./scripts/changelog.sh

# Changes from specific tag to now
./scripts/changelog.sh v1.0.0

# Changes between two tags
./scripts/changelog.sh v1.0.0 v1.1.0
```

### Option 3: Copy to Clipboard

```bash
# macOS
./scripts/changelog.sh | pbcopy

# Linux (requires xclip)
./scripts/changelog.sh | xclip -selection clipboard

# Windows (Git Bash)
./scripts/changelog.sh | clip
```

## Commit Message Format

For automatic categorization, use **Conventional Commits**:

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Changelog Section | Description |
|------|-------------------|-------------|
| `feat` | **Added** | New features |
| `fix` | **Fixed** | Bug fixes |
| `docs` | **Documentation** | Documentation changes |
| `style` | **Other** | Code style (formatting, etc.) |
| `refactor` | **Changed** | Code refactoring |
| `perf` | **Changed** | Performance improvements |
| `test` | **Tests** | Adding/updating tests |
| `build` | **Build/CI** | Build system changes |
| `ci` | **Build/CI** | CI configuration |
| `chore` | **Other** | Maintenance tasks |

### Breaking Changes

Add `BREAKING CHANGE:` in the commit body:

```bash
git commit -m "feat: redesign authentication API

BREAKING CHANGE: Auth endpoints now require API version header"
```

This will create a **⚠️ BREAKING CHANGES** section at the top of the changelog.

## Examples

### Good Commit Messages

```bash
# Feature
git commit -m "feat: add user profile endpoint"
# → Added: add user profile endpoint

# Bug fix
git commit -m "fix: resolve JWT token expiration issue"
# → Fixed: resolve JWT token expiration issue

# With scope
git commit -m "feat(auth): implement OAuth2 support"
# → Added: implement OAuth2 support

# Performance improvement
git commit -m "perf: optimize database queries for class listing"
# → Changed: optimize database queries for class listing (performance)

# Breaking change
git commit -m "feat: upgrade to Go 1.24

BREAKING CHANGE: Minimum Go version is now 1.24"
# → ⚠️ BREAKING CHANGES: upgrade to Go 1.24
```

### Bad Commit Messages (Won't Be Categorized Well)

```bash
# Too vague
git commit -m "updates"

# Not following convention
git commit -m "Fixed a bug"

# No description
git commit -m "feat:"
```

## Generated Output Example

```markdown
## [1.2.0] - 2025-11-22

### ⚠️ BREAKING CHANGES

- Redesign authentication API (Auth endpoints now require API version header)

### Added

- User profile endpoint
- Photo gallery with pagination
- Email notification system

### Fixed

- JWT token expiration issue
- Rate limiting not working for proxied requests
- Photo upload failing for large files

### Changed

- Optimize database queries for class listing (performance)
- Refactor authentication middleware

### Documentation

- Update API documentation
- Add deployment guide

### Tests

- Add unit tests for auth middleware
- Add integration tests for photo upload
```

## Workflow

### Daily Development

1. **Write good commit messages** following Conventional Commits
2. Commits accumulate in git history
3. When ready to release, run `./scripts/release.sh`
4. Changelog is automatically generated and shown
5. Review and confirm to insert into CHANGELOG.md

### Manual Review and Edit

Even with auto-generation, you should:

1. **Review the generated content** - Check for accuracy
2. **Combine similar items** - Group related changes
3. **Add context** - Clarify technical changes for users
4. **Remove noise** - Delete WIP or internal commits
5. **Reorder** - Put most important changes first

### Example: Before and After

**Before (Auto-generated):**
```markdown
### Added
- add user endpoint
- implement authentication
- add tests for auth
```

**After (Manual edit):**
```markdown
### Added
- User authentication with JWT tokens
- User profile management with role-based access control
```

## Best Practices

### 1. Commit Often with Good Messages

```bash
# ✅ Good - Atomic commits with clear messages
git commit -m "feat: add user registration endpoint"
git commit -m "test: add unit tests for user registration"
git commit -m "docs: update API documentation for user endpoints"

# ❌ Bad - Large commits with vague messages
git commit -m "various updates"
```

### 2. Use Scopes for Context

```bash
# ✅ Good - Scopes add context
git commit -m "feat(api): add pagination to class listing"
git commit -m "fix(docker): resolve PostgreSQL connection timeout"

# ✅ Also good - Without scope when obvious
git commit -m "feat: add photo upload"
```

### 3. Write User-Focused Descriptions

```bash
# ✅ Good - User can understand the benefit
git commit -m "feat: add real-time notifications for new messages"

# ❌ Bad - Too technical
git commit -m "feat: implement WebSocket handler with goroutine pool"
```

### 4. Include Breaking Change Details

```bash
# ✅ Good - Clear migration path
git commit -m "feat: redesign authentication API

BREAKING CHANGE: Authentication endpoints have been redesigned.
- /auth/login is now /api/v1/auth/login
- Token format changed from JWT to OAuth2
- Migration guide: https://docs.example.com/migration"
```

## Configuration

The changelog generator can be customized by editing `scripts/changelog.sh`:

- **Date format**: Line ~90 - `CURRENT_DATE=$(date +%Y-%m-%d)`
- **Section order**: Lines ~160-240 - Reorder sections
- **Filtering**: Add commit filtering logic in the categorization loop

## Troubleshooting

### No commits found

```bash
# Check git history
git log --oneline

# If no tags exist
./scripts/changelog.sh $(git rev-list --max-parents=0 HEAD) HEAD
```

### Wrong commits included

```bash
# Generate for specific range
./scripts/changelog.sh v1.0.0 v1.1.0
```

### Commits not categorized

Ensure commits follow Conventional Commits format:
- Starts with type: `feat:`, `fix:`, etc.
- Has description after colon
- Optional scope: `feat(scope):`

## Integration with CI/CD

You can integrate this into your CI pipeline:

```yaml
# .github/workflows/changelog-check.yml
name: Changelog Check

on: [pull_request]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Check commit messages
        run: |
          ./scripts/changelog.sh > /tmp/changelog.md
          if ! grep -q "### Added\|### Fixed" /tmp/changelog.md; then
            echo "⚠️ No significant changes detected"
            echo "Make sure your commits follow Conventional Commits"
          fi
```

## Tips

1. **Preview before release**: Run `./scripts/changelog.sh` before releasing to see what will be generated
2. **Edit after generation**: Auto-generation is a starting point, always review and improve
3. **Keep Unreleased section**: Add breaking changes or important notes in the Unreleased section manually
4. **Link to issues**: Add issue references manually after generation (e.g., "Fixed #123")

## Related Documentation

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [RELEASE_GUIDE.md](RELEASE_GUIDE.md)
- [VERSIONING.md](VERSIONING.md)

## Questions?

See [CONTRIBUTING.md](../CONTRIBUTING.md) for more details on commit message conventions.
