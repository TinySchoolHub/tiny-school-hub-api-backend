# 🚀 Development Workflow Guide

**Last Updated**: 22 November 2025

This guide explains how to work with the Tiny School Hub API backend repository, from daily development to releasing new versions.

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Daily Development Workflow](#daily-development-workflow)
3. [Creating a Release](#creating-a-release)
4. [Hotfix Workflow](#hotfix-workflow)
5. [What Happens Automatically](#what-happens-automatically)
6. [Common Commands Cheat Sheet](#common-commands-cheat-sheet)
7. [Troubleshooting](#troubleshooting)

---

## 🎯 Quick Start

### Branch Overview

```
main          → Production code (always stable)
  └─ develop  → Next release (integration branch)
       ├─ feature/xxx  → New features
       ├─ bugfix/xxx   → Bug fixes
       └─ renovate/*   → Dependency updates
```

### Where Does What Happen?

| Activity            | Location         | Tool              |
| ------------------- | ---------------- | ----------------- |
| Write code          | 💻 Your computer  | Your editor       |
| Create PR           | 🌐 GitHub.com     | Web UI            |
| Run tests           | ☁️ GitHub Actions | Automatic         |
| Review code         | 🌐 GitHub.com     | Web UI            |
| Merge PR            | 🌐 GitHub.com     | Web UI            |
| Create tags         | 💻 Your computer  | `git tag` command |
| Build Docker images | ☁️ GitHub Actions | Automatic         |

---

## 💻 Daily Development Workflow

### Step 1: Start a New Feature

```bash
# Make sure you're on the latest develop branch
git checkout develop
git pull origin develop

# Create your feature branch (use descriptive names!)
git checkout -b feature/add-parent-notifications

# Example branch names:
# ✅ feature/add-user-authentication
# ✅ feature/class-roster-view
# ✅ bugfix/fix-login-timeout
# ❌ feature/changes (too vague)
# ❌ my-branch (doesn't describe the work)
```

### Step 2: Work on Your Feature

```bash
# Make your changes
# ... edit files ...

# Check what you changed
git status
git diff

# Add your changes
git add .

# Commit with a good message (follows Conventional Commits)
git commit -m "feat: add email notifications for parents"

# More commit message examples:
# feat: add new API endpoint for absences
# fix: resolve timezone bug in attendance
# docs: update API documentation
# refactor: simplify authentication logic
# test: add tests for class roster
```

### Step 3: Push Your Branch

```bash
# Push to GitHub (first time)
git push -u origin feature/add-parent-notifications

# For subsequent pushes (after more commits)
git push
```

### Step 4: Create a Pull Request

1. Go to **GitHub.com** → Your repository
2. Click **"Pull requests"** tab
3. Click **"New pull request"**
4. **Base branch**: `develop`
5. **Compare branch**: `feature/add-parent-notifications`
6. Click **"Create pull request"**
7. Fill in the description (template will auto-fill)
8. Click **"Create pull request"**

### Step 5: Wait for CI/CD Checks

**Automatically runs:**
- ✅ Linting (code style check)
- ✅ Tests (all unit tests)
- ✅ Database migrations validation
- ✅ Docker build test
- ✅ Security scan

**You'll see:**
- 🟡 Yellow checks = Running
- ✅ Green checks = Passed
- ❌ Red checks = Failed (needs fixing)

### Step 6: Get Review & Merge

1. Wait for checks to pass ✅
2. Request review from team member (or self-review if solo)
3. Address any review comments
4. Once approved, click **"Merge pull request"**
5. Choose **"Squash and merge"** (recommended)
6. Confirm merge
7. **Delete the branch** (cleanup)

**Done!** ✨ Your code is now in `develop` branch.

---

## 🎉 Creating a Release

When you're ready to deploy to production, follow these steps:

### Step 1: Prepare the Release (On Your Computer)

```bash
# Make sure develop is up to date
git checkout develop
git pull origin develop

# Create a release branch
git checkout -b release/v1.2.0

# Run the release script
./scripts/release.sh 1.2.0

# The script will:
# ✅ Run all tests
# ✅ Update VERSION file
# ✅ Generate CHANGELOG
# ✅ Create commit
# ✅ Push the release branch

# Alternative: Let the script bump the version
./scripts/release.sh patch   # 1.1.5 → 1.1.6
./scripts/release.sh minor   # 1.1.5 → 1.2.0
./scripts/release.sh major   # 1.1.5 → 2.0.0
```

### Step 2: Create Release PR to Main

1. Go to **GitHub.com** → Your repository
2. You'll see a prompt: **"Compare & pull request"** for your release branch
3. Click it (or create PR manually)
4. **Base branch**: `main` ⚠️ (Important!)
5. **Compare branch**: `release/v1.2.0`
6. Review the changes (VERSION, CHANGELOG)
7. Click **"Create pull request"**

### Step 3: Wait for Release Validation

**GitHub Actions will automatically run:**
- ✅ Validate branch name (must be `release/*` or `hotfix/*`)
- ✅ Check VERSION file was updated
- ✅ Check CHANGELOG was updated
- ✅ Run full test suite
- ✅ Run integration tests
- ✅ Build Docker image
- ✅ Security scans (Gosec + Trivy)
- ✅ Generate release summary

**This takes ~5-10 minutes.**

### Step 4: Merge to Main

1. Wait for all checks to pass ✅
2. Get approval (if required)
3. Click **"Merge pull request"**
4. Choose **"Create a merge commit"** (preserves release history)
5. Confirm merge
6. **Don't delete the branch yet!**

### Step 5: Create and Push the Version Tag

⚠️ **IMPORTANT: This step is MANUAL!**

```bash
# Switch to main and pull the merged changes
git checkout main
git pull origin main

# Create an annotated tag
git tag -a v1.2.0 -m "Release v1.2.0"

# Push the tag
git push origin v1.2.0
```

### Step 6: Automatic Deployment

**The moment you push the tag, GitHub Actions automatically:**

1. 🧪 Runs all tests
2. 🐳 Builds Docker image
3. 📦 Publishes to GitHub Container Registry:
   - `ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2.0`
   - `ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2`
   - `ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1`
   - `ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:latest`
4. 📝 Creates GitHub Release with changelog
5. 📊 Updates Helm chart version

**Check progress at:**
- GitHub → Actions tab → "Release" workflow

### Step 7: Sync Changes Back to Develop

⚠️ **Don't forget this step!** Otherwise develop and main will drift apart.

**Option A: Create PR (Recommended)**

1. Go to **GitHub.com** → Pull requests
2. Create new PR
3. **Base branch**: `develop`
4. **Compare branch**: `main`
5. Title: `chore: sync main back to develop after v1.2.0 release`
6. Create and merge the PR

**Option B: Cherry-pick (Faster)**

```bash
# Get the commit hash of the version bump
git log main --oneline -n 5

# Find the commit like: "chore: bump version to v1.2.0"
# Copy its hash (e.g., abc1234)

git checkout develop
git pull origin develop
git cherry-pick abc1234
git push origin develop
```

### Step 8: Deploy to Production

Use the Docker image that was just published:

```bash
# Pull the new image
docker pull ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2.0

# Or update your Kubernetes deployment
kubectl set image deployment/tiny-school-hub \
  api=ghcr.io/tinyschoolhub/tiny-school-hub-api-backend:v1.2.0
```

**Done!** 🎉 Your release is live!

---

## 🚨 Hotfix Workflow

When you need to fix a critical bug in production **immediately**:

### Step 1: Create Hotfix Branch

```bash
# Start from production (main)
git checkout main
git pull origin main

# Create hotfix branch
git checkout -b hotfix/v1.2.1
```

### Step 2: Fix the Issue

```bash
# Make your fix
# ... edit files ...

# Commit the fix
git add .
git commit -m "fix: resolve critical security vulnerability"
```

### Step 3: Update Version

```bash
# Update VERSION file
echo "1.2.1" > VERSION

# Update CHANGELOG.md
# Add entry for v1.2.1 with the fix

# Commit version bump
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v1.2.1"

# Push the branch
git push -u origin hotfix/v1.2.1
```

### Step 4: Create TWO Pull Requests

**PR 1: Deploy to Production**
- Base: `main`
- Compare: `hotfix/v1.2.1`
- Merge this FIRST

**PR 2: Incorporate Fix into Development**
- Base: `develop`
- Compare: `hotfix/v1.2.1`
- Merge this AFTER PR 1

### Step 5: Tag and Deploy

```bash
# After PR to main is merged
git checkout main
git pull origin main

# Create tag
git tag -a v1.2.1 -m "Hotfix v1.2.1 - Security patch"
git push origin v1.2.1
```

**Automatic deployment happens immediately!**

---

## ⚙️ What Happens Automatically

### When You Create a PR to `develop`:

- ✅ **pr-develop.yml** workflow runs
  - Linting with golangci-lint
  - All tests with race detection
  - Database migration validation
  - Docker image build test
  - Security scan (Gosec)
  - PR title validation (Conventional Commits)
  - Adds helpful comment

### When You Create a PR to `main`:

- ✅ **pr-main.yml** workflow runs
  - Validates branch name (must be release/* or hotfix/*)
  - Checks VERSION file updated
  - Checks CHANGELOG updated
  - Full test suite
  - Integration tests
  - Docker build
  - Security scans (Gosec + Trivy)
  - Generates release summary

### When You Push a Tag:

- ✅ **release.yml** workflow runs
  - Runs all tests
  - Builds Docker image
  - Pushes to GitHub Container Registry
  - Creates GitHub Release
  - Updates Helm chart

### When Renovate Creates a PR:

- ✅ **pr-develop.yml** runs (normal validation)
- ✅ **renovate-automerge.yml** runs
  - Auto-approves if labeled "automerge"
  - Renovate then auto-merges if CI passes

---

## 📝 Common Commands Cheat Sheet

### Starting Work

```bash
# Update develop
git checkout develop && git pull

# Create feature
git checkout -b feature/my-feature

# Create bugfix
git checkout -b bugfix/fix-something
```

### Committing

```bash
# Check status
git status

# Add all changes
git add .

# Commit (use conventional format!)
git commit -m "feat: add new feature"
git commit -m "fix: resolve bug"
git commit -m "docs: update README"
```

### Pushing

```bash
# First push
git push -u origin feature/my-feature

# Subsequent pushes
git push
```

### Updating Your Branch

```bash
# Get latest changes from develop
git checkout develop
git pull

# Update your feature branch
git checkout feature/my-feature
git merge develop

# Or use rebase (cleaner history)
git rebase develop
```

### Releasing

```bash
# Quick release
git checkout develop && git pull
git checkout -b release/v1.2.0
./scripts/release.sh 1.2.0
# → Create PR to main on GitHub
# → After merge, create tag:
git checkout main && git pull
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### Checking Status

```bash
# Current branch
git branch --show-current

# Current version
cat VERSION

# Latest tag
git describe --tags --abbrev=0

# Uncommitted changes
git status

# Recent commits
git log --oneline -10
```

---

## 🔧 Troubleshooting

### Problem: CI Checks Failing

**Symptoms:** ❌ Red X on your PR

**Solutions:**

1. **Check the logs**
   - Click on the failed check
   - Read the error message

2. **Common issues:**
   - **Linting errors**: Run `make lint` locally
   - **Test failures**: Run `go test ./...` locally
   - **Build errors**: Run `docker build .` locally

3. **Fix and push**
   ```bash
   # Make fixes
   git add .
   git commit -m "fix: resolve CI issues"
   git push
   ```

### Problem: Can't Push to Main or Develop

**Error:** `remote: error: GH006: Protected branch update failed`

**Cause:** Branch protection is enabled (good!)

**Solution:** Always create PRs, never push directly.

### Problem: Forgot to Create Tag After Release

**Symptoms:** Release merged to main but no Docker images

**Solution:**
```bash
# Check what version you released
cat VERSION  # Shows: 1.2.0

# Create the tag now
git checkout main
git pull
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### Problem: Need to Redo a Tag

⚠️ **Not recommended, but if absolutely necessary:**

```bash
# Delete local tag
git tag -d v1.2.0

# Delete remote tag
git push origin :refs/tags/v1.2.0

# Recreate correctly
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### Problem: Merge Conflict

**When merging your PR or updating branch:**

```bash
# Update your branch
git checkout feature/my-feature
git fetch origin
git merge origin/develop

# If conflicts occur:
# 1. Open conflicted files (marked with <<<<<<<)
# 2. Resolve conflicts manually
# 3. Mark as resolved:
git add .
git commit -m "chore: resolve merge conflicts"
git push
```

### Problem: Forgot to Sync Main Back to Develop

**Symptoms:** Develop missing VERSION or CHANGELOG changes

**Solution:**
```bash
# Quick fix with cherry-pick
git checkout develop
git pull
git log main --oneline -n 10  # Find version bump commit
git cherry-pick <commit-hash>
git push origin develop
```

### Problem: Made Commit to Wrong Branch

**Solution:**
```bash
# If you committed to develop instead of a feature branch:
git checkout develop
git reset --soft HEAD~1  # Undo commit, keep changes
git stash  # Save changes

git checkout -b feature/my-feature
git stash pop  # Restore changes
git add .
git commit -m "feat: my changes"
```

---

## ✅ Best Practices

### Commit Messages

**Good Examples:**
```
feat: add parent notification system
fix: resolve login timeout issue
docs: update API documentation for v1.2
refactor: simplify database connection logic
test: add unit tests for authentication
chore: update dependencies
```

**Bad Examples:**
```
update stuff
fix bug
changes
WIP
asdfasdf
```

### Branch Names

**Good Examples:**
```
feature/add-parent-portal
feature/class-roster-view
bugfix/fix-login-error
hotfix/v1.2.1
release/v1.3.0
```

**Bad Examples:**
```
my-branch
test
fix
new-feature
changes
```

### Pull Request Titles

Follow the same format as commits:

```
feat: Add parent notification system
fix: Resolve login timeout issue
chore: Update dependencies to latest versions
```

### Before Creating a PR

- [ ] Run tests locally: `go test ./...`
- [ ] Run linter: `make lint` or `golangci-lint run`
- [ ] Build Docker image: `docker build .`
- [ ] Update documentation if needed
- [ ] Write clear PR description
- [ ] Link related issues

### Before Merging a PR

- [ ] All CI checks passed ✅
- [ ] Code reviewed and approved
- [ ] No merge conflicts
- [ ] Tested locally if possible
- [ ] Documentation updated

---

## 🆘 Getting Help

### Check Documentation

- `docs/CICD_ANALYSIS.md` - Complete CI/CD analysis
- `docs/BRANCH_PROTECTION_RULES.md` - Branch protection guide
- `docs/GITFLOW.md` - Detailed GitFlow explanation
- `docs/RELEASE_GUIDE.md` - Release process details

### Check Workflow Status

- GitHub → Actions tab
- See real-time progress of CI/CD
- Download logs if needed

### Common Links

- **Repository**: https://github.com/TinySchoolHub/tiny-school-hub-api-backend
- **Pull Requests**: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/pulls
- **Actions**: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/actions
- **Releases**: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/releases
- **Docker Images**: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/pkgs/container/tiny-school-hub-api-backend

---

## 📊 Workflow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        DEVELOPMENT CYCLE                        │
└─────────────────────────────────────────────────────────────────┘

    LOCAL                          GITHUB
    ─────                          ──────

 1. Create feature
    feature/xxx
         │
         ├─────────────────────▶  Push branch
         │                              │
         │                              ▼
         │                         Create PR to develop
         │                              │
         │                              ▼
         │                         pr-develop.yml runs
         │                         (lint, test, build)
         │                              │
         │                              ▼
         │                         Review & Approve
         │                              │
         │                              ▼
         │                         Merge to develop ✅
         │
         │
 2. Ready to release?
    release/v1.2.0
         │
         ├─────────────────────▶  Push release branch
         │                              │
         │                              ▼
         │                         Create PR to main
         │                              │
         │                              ▼
         │                         pr-main.yml runs
         │                         (full validation)
         │                              │
         │                              ▼
         │                         Review & Approve
         │                              │
         │                              ▼
         │                         Merge to main ✅
         │                              │
         ▼                              │
 3. Create tag                          │
    git tag v1.2.0 ────────────▶  Push tag
                                        │
                                        ▼
                                   release.yml runs
                                   (build & publish)
                                        │
                                        ├─▶ Docker images
                                        ├─▶ GitHub Release
                                        └─▶ Helm chart
                                        
 4. Sync back
    main ──────────────────────▶  Create PR to develop
                                        │
                                        ▼
                                   Merge ✅
```

---

**Remember:** 
- ✅ Always work in feature branches
- ✅ Always create PRs (never push directly to main/develop)
- ✅ Always wait for CI to pass
- ✅ Always sync main back to develop after releases

**Happy coding!** 🚀
