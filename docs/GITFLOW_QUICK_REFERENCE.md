# GitFlow Quick Reference

## Common Commands

### Daily Development

```bash
# Start working on a new feature
git checkout develop
git pull origin develop
git checkout -b feature/my-feature

# Make changes and commit
git add .
git commit -m "feat: add new feature"
git push -u origin feature/my-feature

# Create PR on GitHub: feature/my-feature → develop
```

### Bug Fixes

```bash
# Fix a bug in develop
git checkout develop
git pull origin develop
git checkout -b bugfix/fix-something

# Make changes
git add .
git commit -m "fix: resolve issue"
git push -u origin bugfix/fix-something

# Create PR: bugfix/fix-something → develop
```

### Preparing a Release

```bash
# Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.2.0

# Bump version
echo "1.2.0" > VERSION

# Update CHANGELOG.md
# Add release notes for v1.2.0

# Commit
git add VERSION CHANGELOG.md
git commit -m "chore: prepare release v1.2.0"
git push -u origin release/v1.2.0

# Create PRs:
# 1. release/v1.2.0 → main (for release)
# 2. release/v1.2.0 → develop (to sync back)
```

### Hotfix Production

```bash
# Critical fix needed in production
git checkout main
git pull origin main
git checkout -b hotfix/v1.2.1

# Fix the issue
git add .
git commit -m "fix: critical bug"

# Bump version
echo "1.2.1" > VERSION
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v1.2.1"
git push -u origin hotfix/v1.2.1

# Create PRs:
# 1. hotfix/v1.2.1 → main (deploy now)
# 2. hotfix/v1.2.1 → develop (incorporate fix)
```

## Branch Naming

| Type    | Pattern               | Example                       |
| ------- | --------------------- | ----------------------------- |
| Feature | `feature/description` | `feature/user-authentication` |
| Bugfix  | `bugfix/description`  | `bugfix/login-error`          |
| Hotfix  | `hotfix/version`      | `hotfix/v1.2.1`               |
| Release | `release/version`     | `release/v1.3.0`              |

## Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting, missing semicolons, etc)
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `test`: Adding or updating tests
- `build`: Build system changes
- `ci`: CI configuration changes
- `chore`: Other changes (dependencies, etc)
- `revert`: Revert a previous commit

### Examples

```bash
# Simple feature
git commit -m "feat: add user profile endpoint"

# Bug fix with reference
git commit -m "fix: resolve null pointer in auth handler

Fixes #123"

# Breaking change
git commit -m "feat!: change API response format

BREAKING CHANGE: API now returns data in a different format"
```

## PR Workflow

### 1. Create Feature Branch

```bash
git checkout develop
git pull origin develop
git checkout -b feature/my-feature
```

### 2. Develop & Commit

```bash
# Make changes
git add .
git commit -m "feat: implement feature"

# Keep up to date with develop
git fetch origin
git rebase origin/develop

# Push
git push origin feature/my-feature
```

### 3. Create Pull Request

- Go to GitHub
- Create PR: `feature/my-feature` → `develop`
- Fill out PR template
- Request review

### 4. Address Review Comments

```bash
# Make requested changes
git add .
git commit -m "fix: address review comments"
git push origin feature/my-feature
```

### 5. Merge & Clean Up

After PR is approved and merged:

```bash
# Switch to develop
git checkout develop
git pull origin develop

# Delete feature branch
git branch -d feature/my-feature
git push origin --delete feature/my-feature
```

## Resolving Conflicts

```bash
# Update your branch with latest develop
git checkout feature/my-feature
git fetch origin
git rebase origin/develop

# If conflicts occur, resolve them
# Edit conflicting files
git add <resolved-files>
git rebase --continue

# Push (force with lease for safety)
git push --force-with-lease origin feature/my-feature
```

## Checking Status

```bash
# See current branch
git branch --show-current

# See all branches
git branch -a

# See commits ahead/behind
git status

# See recent commits
git log --oneline -10

# See branch history
git log --oneline --graph --all --decorate
```

## Renovate Dependency Updates

Renovate automatically creates PRs for dependency updates.

### Auto-merge (Patch Updates)

Patch updates with `automerge` label will auto-merge if CI passes.

### Manual Review (Minor/Major)

1. Check the PR created by Renovate
2. Review the changelog of the updated dependency
3. Run tests locally if needed
4. Approve and merge to `develop`

### Configuration

Edit `renovate.json` to customize:
- Update schedule
- Auto-merge rules
- Grouping rules
- Ignored dependencies

## Release Process

### Regular Release

1. **Prepare** (from `develop`):
   ```bash
   git checkout -b release/v1.2.0 develop
   echo "1.2.0" > VERSION
   # Update CHANGELOG.md
   git commit -am "chore: prepare release v1.2.0"
   ```

2. **PR to main**:
   - Create PR: `release/v1.2.0` → `main`
   - Wait for all checks to pass
   - Get 2 approvals
   - Merge

3. **Tag** (after merge):
   ```bash
   git checkout main
   git pull
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```

4. **Sync back**:
   - Create PR: `release/v1.2.0` → `develop`
   - Merge to keep branches in sync

### Using Release Script

```bash
# From release branch
./scripts/release.sh patch   # 1.2.3 → 1.2.4
./scripts/release.sh minor   # 1.2.3 → 1.3.0
./scripts/release.sh major   # 1.2.3 → 2.0.0
./scripts/release.sh 1.5.0   # Specific version
```

## CI/CD Workflows

### PR to `develop`

- ✅ Lint code
- ✅ Run tests
- ✅ Validate migrations
- ✅ Build Docker image
- ✅ Security scan

### PR to `main`

- ✅ Validate release branch
- ✅ Check VERSION updated
- ✅ Run full test suite
- ✅ Integration tests
- ✅ Build Docker image
- ✅ Security scans
- ✅ Generate release summary

### Tag Push (`v*.*.*`)

- ✅ Run tests
- 🐳 Build production Docker image
- 📦 Push to GHCR
- 📝 Create GitHub Release

## Tips & Best Practices

### Keep Branches Updated

```bash
# Before starting work
git checkout develop
git pull origin develop

# During development
git fetch origin
git rebase origin/develop
```

### Small, Focused PRs

- One feature per PR
- Keep changes under 400 lines when possible
- Write clear descriptions

### Test Before Pushing

```bash
# Run tests
go test ./...

# Run linting
golangci-lint run

# Build Docker image
docker build -t test .
```

### Clean Commit History

```bash
# Squash multiple commits
git rebase -i HEAD~3

# Amend last commit
git commit --amend

# Interactive rebase to clean up
git rebase -i origin/develop
```

## Getting Help

- **GitFlow Guide**: [docs/GITFLOW.md](GITFLOW.md)
- **Release Guide**: [docs/RELEASE_GUIDE.md](RELEASE_GUIDE.md)
- **Contributing**: [CONTRIBUTING.md](../CONTRIBUTING.md)
- **Issues**: Create a GitHub issue

## Useful Git Aliases

Add to your `~/.gitconfig`:

```ini
[alias]
    co = checkout
    br = branch
    ci = commit
    st = status
    lg = log --oneline --graph --all --decorate
    unstage = reset HEAD --
    last = log -1 HEAD
    visual = log --oneline --graph --all --decorate
```

Usage:
```bash
git co develop      # checkout develop
git br -a           # list all branches
git lg              # pretty log
```
