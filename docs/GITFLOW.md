# GitFlow Workflow Guide

## Overview

This project uses **GitFlow** branching strategy to manage development and releases.

## Branch Structure

```
main          ← Production-ready code (stable releases)
  └─ develop  ← Integration branch (next release)
       ├─ feature/xxx  ← New features
       ├─ bugfix/xxx   ← Bug fixes
       └─ hotfix/xxx   ← Urgent production fixes
```

### Main Branches

- **`main`**: Production branch. Only contains stable, released code.
  - Protected branch
  - Only accepts merges from `develop` (releases) or `hotfix/*` branches
  - Every commit is tagged with a version number

- **`develop`**: Integration branch for the next release.
  - Protected branch
  - All features and bugfixes are merged here first
  - Should always be in a releasable state

### Supporting Branches

- **`feature/*`**: New features
  - Branch from: `develop`
  - Merge back to: `develop`
  - Naming: `feature/add-user-authentication`, `feature/class-roster`

- **`bugfix/*`**: Bug fixes during development
  - Branch from: `develop`
  - Merge back to: `develop`
  - Naming: `bugfix/fix-login-error`, `bugfix/correct-timezone`

- **`hotfix/*`**: Urgent fixes for production
  - Branch from: `main`
  - Merge to: `main` AND `develop`
  - Naming: `hotfix/security-patch`, `hotfix/critical-bug`

- **`release/*`**: Prepare for production release
  - Branch from: `develop`
  - Merge to: `main` AND `develop`
  - Naming: `release/v1.2.0`

## Workflow

### 1. Starting a New Feature

```bash
# Make sure develop is up to date
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/my-awesome-feature

# Work on your feature
git add .
git commit -m "feat: add awesome feature"

# Push to remote
git push -u origin feature/my-awesome-feature
```

**Create Pull Request**: `feature/my-awesome-feature` → `develop`

### 2. Fixing a Bug (Non-Critical)

```bash
# Start from develop
git checkout develop
git pull origin develop

# Create bugfix branch
git checkout -b bugfix/fix-issue-123

# Fix the bug
git add .
git commit -m "fix: resolve issue #123"

# Push to remote
git push -u origin bugfix/fix-issue-123
```

**Create Pull Request**: `bugfix/fix-issue-123` → `develop`

### 3. Creating a Release

When `develop` has enough features for a release:

```bash
# Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.2.0

# Update version and changelog
echo "1.2.0" > VERSION

# Commit version bump
git add VERSION CHANGELOG.md
git commit -m "chore: prepare release v1.2.0"

# Push release branch
git push -u origin release/v1.2.0
```

**Create Pull Requests**:
1. `release/v1.2.0` → `main` (will trigger release workflow)
2. `release/v1.2.0` → `develop` (to sync back any release changes)

After merging to `main`:
```bash
# Tag the release
git checkout main
git pull origin main
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### 4. Hotfix for Production

Critical bug in production:

```bash
# Branch from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-security-fix

# Fix the issue
git add .
git commit -m "fix: critical security vulnerability"

# Update version (patch bump)
echo "1.2.1" > VERSION
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v1.2.1"

# Push hotfix branch
git push -u origin hotfix/critical-security-fix
```

**Create Pull Requests**:
1. `hotfix/critical-security-fix` → `main` (deploy immediately)
2. `hotfix/critical-security-fix` → `develop` (incorporate fix)

After merging to `main`:
```bash
git checkout main
git pull origin main
git tag -a v1.2.1 -m "Hotfix v1.2.1"
git push origin v1.2.1
```

## Automated Workflows

### On Feature/Bugfix PR to `develop`:
- ✅ Run tests
- ✅ Run linting
- ✅ Validate migrations
- ✅ Build Docker image

### On Release PR to `main`:
- ✅ Run all tests
- ✅ Build and test Docker image
- ⏸️  Manual approval required

### On Tag Push (v*.*.*):
- ✅ Run full test suite
- 🐳 Build production Docker image
- 📦 Push to GitHub Container Registry
- 📝 Create GitHub Release
- 🚀 Ready for deployment

## Pull Request Guidelines

### PR Title Format
Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add user profile page`
- `fix: correct date formatting`
- `docs: update API documentation`
- `chore: update dependencies`
- `refactor: simplify auth logic`
- `test: add tests for login flow`

### PR Description
Include:
- **What**: Description of changes
- **Why**: Reason for changes
- **How**: Implementation approach
- **Testing**: How you tested the changes
- **Screenshots**: For UI changes (if applicable)

### Code Review
- At least 1 approval required for `develop`
- At least 2 approvals required for `main`
- All CI checks must pass
- No merge conflicts

## Branch Protection Rules

### `main` branch:
- ✅ Require pull request before merging
- ✅ Require 2 approvals
- ✅ Require status checks to pass
- ✅ Require conversation resolution
- ✅ Require linear history
- ✅ Do not allow bypassing settings
- ✅ Restrict who can push (only maintainers)

### `develop` branch:
- ✅ Require pull request before merging
- ✅ Require 1 approval
- ✅ Require status checks to pass
- ✅ Require conversation resolution
- ✅ Allow force pushes (for maintainers only)

## Best Practices

### Commits
- Make small, focused commits
- Write clear commit messages
- Use conventional commit format
- Reference issues when applicable

### Feature Development
- Keep features small and focused
- Merge to `develop` frequently
- Rebase on `develop` regularly to avoid conflicts
- Delete feature branches after merge

### Testing
- Write tests for new features
- Ensure all tests pass locally before pushing
- Add integration tests when needed
- Update test documentation

### Dependencies
- Renovate will automatically create PRs for dependency updates
- Review and merge dependency updates weekly
- Test thoroughly after major version updates
- Keep security updates high priority

## Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes
- **MINOR** (v1.3.0): New features, backwards compatible
- **PATCH** (v1.2.1): Bug fixes, backwards compatible

### When to Bump

- **Patch**: Bug fixes, security patches, minor improvements
- **Minor**: New features, non-breaking API changes
- **Major**: Breaking API changes, major refactoring

## Release Schedule

- **Weekly**: Merge to `develop` continuously
- **Bi-weekly**: Release to production (from `develop` to `main`)
- **As needed**: Hotfixes for critical issues

## Commands Quick Reference

```bash
# Setup
git clone <repository>
git checkout develop

# New feature
git checkout -b feature/my-feature develop
git push -u origin feature/my-feature

# Update from develop
git checkout develop
git pull origin develop
git checkout feature/my-feature
git rebase develop

# Ready to merge
# Create PR: feature/my-feature → develop

# Release
git checkout -b release/v1.2.0 develop
# Update VERSION, CHANGELOG
git push -u origin release/v1.2.0
# Create PR: release/v1.2.0 → main

# Hotfix
git checkout -b hotfix/v1.2.1 main
# Fix issue, update VERSION
git push -u origin hotfix/v1.2.1
# Create PR: hotfix/v1.2.1 → main
# Create PR: hotfix/v1.2.1 → develop
```

## Troubleshooting

### Merge Conflicts
```bash
# Update your branch with latest develop
git checkout develop
git pull origin develop
git checkout your-branch
git rebase develop

# Resolve conflicts
# Edit conflicting files
git add .
git rebase --continue
git push --force-with-lease
```

### Accidental Commit to Wrong Branch
```bash
# Move commits to correct branch
git checkout correct-branch
git cherry-pick <commit-hash>

# Remove from wrong branch
git checkout wrong-branch
git reset --hard HEAD~1
git push --force-with-lease
```

## Resources

- [GitFlow Original Article](https://nvie.com/posts/a-successful-git-branching-model/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Flow vs GitFlow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)
