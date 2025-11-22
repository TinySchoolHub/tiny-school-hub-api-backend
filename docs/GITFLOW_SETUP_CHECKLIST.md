# GitFlow & Renovate Setup Checklist

Use this checklist to ensure your GitFlow workflow and Renovate are properly configured.

## ✅ Initial Setup

### 1. Run Setup Script

```bash
chmod +x scripts/setup-gitflow.sh
./scripts/setup-gitflow.sh
```

This script will:
- ✅ Create `develop` branch
- ✅ Push `develop` to remote
- ✅ Create `.github/CODEOWNERS`
- ✅ Create `.github/pull_request_template.md`
- ✅ Update `CONTRIBUTING.md`
- ✅ Commit and push changes

### 2. Configure Branch Protection (GitHub)

#### For `main` branch:

Go to: **Settings** → **Branches** → **Add rule**

- Branch name pattern: `main`
- [x] Require a pull request before merging
  - [x] Require approvals: **2**
  - [x] Dismiss stale pull request approvals when new commits are pushed
  - [x] Require review from Code Owners
- [x] Require status checks to pass before merging
  - [x] Require branches to be up to date before merging
  - Required status checks:
    - `Lint`
    - `Full Test Suite`
    - `Build and Test Docker Image`
    - `Security Scan`
    - `Pre-Release Checklist`
- [x] Require conversation resolution before merging
- [x] Require linear history
- [x] Do not allow bypassing the above settings
- [x] Restrict who can push to matching branches
  - Add: Repository maintainers only

#### For `develop` branch:

Go to: **Settings** → **Branches** → **Add rule**

- Branch name pattern: `develop`
- [x] Require a pull request before merging
  - [x] Require approvals: **1**
- [x] Require status checks to pass before merging
  - Required status checks:
    - `Lint`
    - `Test`
    - `Validate Migrations`
    - `Build Docker Image`
- [x] Require conversation resolution before merging
- [x] Allow force pushes
  - Specify who can force push: Maintainers only

### 3. Install Renovate

1. Go to: https://github.com/apps/renovate
2. Click **Install** or **Configure**
3. Select your repository: `TinySchoolHub/tiny-school-hub-api-backend`
4. Grant permissions
5. Renovate will automatically detect `renovate.json` configuration

**What Renovate will do:**
- Create dependency update PRs to `develop` branch
- Group Go dependencies together
- Run weekly (Monday before 5am)
- Auto-merge patch updates with `automerge` label
- Create a Dependency Dashboard issue

### 4. Configure GitHub Secrets (if needed)

Go to: **Settings** → **Secrets and variables** → **Actions**

If using auto-merge for Renovate:
- Add `PAT_TOKEN` (Personal Access Token with `repo` scope)

### 5. Configure CODEOWNERS

Edit `.github/CODEOWNERS` to add team members:

```
# Default owners
* @fabien @team-member-1 @team-member-2

# Specific paths
/internal/http/ @backend-team
/deploy/ @devops-team
```

## ✅ Verify Setup

### Test Branch Protection

```bash
# Try to push directly to main (should fail)
git checkout main
echo "test" >> test.txt
git add test.txt
git commit -m "test"
git push origin main
# Expected: Error - protected branch

# Clean up
git reset --hard HEAD~1
```

### Test PR Workflow

```bash
# Create test feature branch
git checkout develop
git pull origin develop
git checkout -b feature/test-gitflow

# Make a small change
echo "# Test" >> docs/TEST.md
git add docs/TEST.md
git commit -m "feat: test GitFlow setup"
git push -u origin feature/test-gitflow

# Go to GitHub and create PR: feature/test-gitflow → develop
# Verify:
# - PR template is populated
# - CI workflows run
# - PR requires approval
# - Can't merge without passing checks
```

### Test Renovate

After installing Renovate:

1. Check for Dependency Dashboard issue
2. Wait for first PR (or trigger manually)
3. Verify PR follows format and targets `develop`
4. Test auto-merge on patch update PRs

## ✅ Daily Workflow Verification

### Starting a Feature

```bash
git checkout develop
git pull origin develop
git checkout -b feature/my-feature

# Work on feature
git add .
git commit -m "feat: implement feature"
git push -u origin feature/my-feature

# Create PR on GitHub
```

### Creating a Release

```bash
git checkout develop
git pull origin develop
git checkout -b release/v1.2.0

# Update VERSION and CHANGELOG.md
echo "1.2.0" > VERSION
# Edit CHANGELOG.md

git add VERSION CHANGELOG.md
git commit -m "chore: prepare release v1.2.0"
git push -u origin release/v1.2.0

# Create PRs:
# 1. release/v1.2.0 → main
# 2. release/v1.2.0 → develop
```

### Handling a Hotfix

```bash
git checkout main
git pull origin main
git checkout -b hotfix/v1.2.1

# Fix critical issue
git add .
git commit -m "fix: critical security issue"

# Update version
echo "1.2.1" > VERSION
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v1.2.1"

git push -u origin hotfix/v1.2.1

# Create PRs:
# 1. hotfix/v1.2.1 → main
# 2. hotfix/v1.2.1 → develop
```

## ✅ Troubleshooting

### Branch Protection Not Working

- Verify you're a repository admin
- Check Settings → Branches → Branch protection rules
- Ensure status check names match workflow names exactly

### Renovate Not Creating PRs

- Check Renovate app is installed: https://github.com/apps/renovate
- Look at Renovate logs: Repository → Actions → Renovate
- Verify `renovate.json` is valid: https://docs.renovatebot.com/config-validation/
- Check Dependency Dashboard issue for errors

### CI Workflows Not Running

- Go to Actions tab and check workflow runs
- Verify `.github/workflows/` files are present
- Check workflow trigger conditions match your branch names
- Ensure GitHub Actions are enabled: Settings → Actions → General

### Can't Create PRs

- Ensure you have write access to repository
- Check branch exists on remote: `git push -u origin your-branch`
- Verify base branch exists (`develop` or `main`)

### Auto-merge Not Working

- Check PR has `automerge` label
- Verify all required checks are passing
- Ensure auto-merge is enabled in repository settings
- Check if branch protection requires manual approval

## 📋 Maintenance Tasks

### Weekly

- [ ] Review and merge Renovate PRs
- [ ] Check CI pipeline health
- [ ] Review open PRs

### Bi-weekly

- [ ] Create release from `develop` to `main`
- [ ] Deploy to staging/production
- [ ] Update documentation if needed

### Monthly

- [ ] Review and update branch protection rules
- [ ] Check Renovate configuration
- [ ] Audit security alerts
- [ ] Review CHANGELOG

### Quarterly

- [ ] Review GitFlow process effectiveness
- [ ] Update team documentation
- [ ] Review and update CODEOWNERS
- [ ] Optimize CI/CD pipelines

## 📚 References

- [GitFlow Documentation](GITFLOW.md)
- [Quick Reference](GITFLOW_QUICK_REFERENCE.md)
- [Release Guide](RELEASE_GUIDE.md)
- [Renovate Docs](https://docs.renovatebot.com/)
- [GitHub Branch Protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)

## 🎯 Success Criteria

Your GitFlow setup is complete when:

- ✅ `develop` branch exists and is protected
- ✅ Branch protection rules are configured for `main` and `develop`
- ✅ Renovate is installed and creating PRs
- ✅ CI workflows run on all PRs
- ✅ Team understands the workflow
- ✅ Documentation is up to date
- ✅ First test PR successfully merged

## 🚨 Emergency Procedures

### Bypass Branch Protection (Emergency Only)

If you absolutely must bypass protection:

1. Go to Settings → Branches → Edit rule
2. Temporarily disable "Do not allow bypassing"
3. Make critical change
4. Re-enable protection immediately

**Always document why and notify team!**

### Rollback a Bad Release

```bash
# Revert the merge commit
git checkout main
git pull origin main
git revert -m 1 <merge-commit-hash>
git push origin main

# Tag the revert
git tag -a v1.2.1-revert -m "Revert v1.2.1"
git push origin v1.2.1-revert

# Fix issue in develop
git checkout develop
# Fix and test thoroughly
# Create new release
```

### Fix Broken CI

If CI is blocking all PRs:

1. Create hotfix branch from `main`
2. Fix CI configuration
3. Test locally with act: `act -j test`
4. PR to `main` with emergency label
5. Merge and backport to `develop`

---

**Need Help?**
- Create an issue with `question` label
- Check documentation in `docs/` folder
- Contact repository maintainers
