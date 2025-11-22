# Branch Protection Rules - Configuration Guide

## Quick Setup via GitHub UI

Navigate to: **Settings** → **Branches** → **Branch protection rules**

---

## 🛡️ Protection Rules for `main` Branch

### Basic Settings

- **Branch name pattern**: `main`
- **Require a pull request before merging**: ✅ Enabled
  - **Require approvals**: 1
  - **Dismiss stale pull request approvals when new commits are pushed**: ✅
  - **Require review from Code Owners**: ✅ (if you want stricter control)

### Status Checks

- **Require status checks to pass before merging**: ✅ Enabled
- **Require branches to be up to date before merging**: ✅ Enabled
- **Status checks that are required**:
  - `Validate Release`
  - `Lint`
  - `Full Test Suite`
  - `Integration Tests`
  - `Build and Test Docker Image`
  - `Security Scan`
  - `Pre-Release Checklist`

### Additional Rules

- **Require conversation resolution before merging**: ✅ Enabled
- **Require signed commits**: ⚠️ Optional (recommended for security)
- **Require linear history**: ⚠️ Optional (cleaner history)
- **Do not allow bypassing the above settings**: ✅ Enabled
- **Restrict who can push to matching branches**: ✅ Enabled
  - Add: Repository administrators only
- **Allow force pushes**: ❌ Disabled
- **Allow deletions**: ❌ Disabled

---

## 🔧 Protection Rules for `develop` Branch

### Basic Settings

- **Branch name pattern**: `develop`
- **Require a pull request before merging**: ✅ Enabled
  - **Require approvals**: 1
  - **Dismiss stale pull request approvals when new commits are pushed**: ✅
  - **Require review from Code Owners**: ⚠️ Optional

### Status Checks

- **Require status checks to pass before merging**: ✅ Enabled
- **Require branches to be up to date before merging**: ✅ Enabled
- **Status checks that are required**:
  - `Lint`
  - `Test`
  - `Validate Migrations`
  - `Build Docker Image`
  - `Security Scan`
  - `PR Validation`

### Additional Rules

- **Require conversation resolution before merging**: ✅ Enabled
- **Do not allow bypassing the above settings**: ⚠️ Optional (you might want flexibility on develop)
- **Allow force pushes**: ❌ Disabled
- **Allow deletions**: ❌ Disabled

---

## 🔐 Protection Rules for `release/*` and `hotfix/*` Branches

### Pattern-Based Rules

Create two separate rules:

#### Rule 1: `release/*`

- **Branch name pattern**: `release/*`
- **Require a pull request before merging**: ❌ Disabled (these branches are temporary)
- **Allow force pushes**: ❌ Disabled
- **Allow deletions**: ✅ Enabled (can be deleted after merge)

#### Rule 2: `hotfix/*`

- **Branch name pattern**: `hotfix/*`
- **Require a pull request before merging**: ❌ Disabled (these branches are temporary)
- **Allow force pushes**: ❌ Disabled
- **Allow deletions**: ✅ Enabled (can be deleted after merge)

---

## 📋 Branch Protection Checklist

After configuring, verify:

- [ ] `main` branch is locked down (no direct pushes)
- [ ] `develop` branch requires PR + CI to pass
- [ ] Only release/hotfix branches can PR to `main`
- [ ] Feature branches can only target `develop`
- [ ] CI status checks are enforced
- [ ] At least 1 approval required for merges

---

## 🧪 Testing Branch Protection

### Test 1: Direct Push to Main (Should Fail)

```bash
git checkout main
echo "test" >> README.md
git add README.md
git commit -m "test: direct push"
git push origin main
# Expected: ❌ Protected branch update failed
```

### Test 2: PR Without CI (Should Block)

1. Create feature branch
2. Push changes
3. Create PR to `develop`
4. Try to merge before CI completes
   - Expected: ❌ Merge button disabled

### Test 3: PR Without Approval (Should Block)

1. Create PR to `main` from `release/v1.0.0`
2. Wait for CI to pass
3. Try to merge without approval
   - Expected: ❌ Requires 1 approval

---

## 🎯 Recommended Configuration

### For Solo Developer (Current State)

If you're working alone initially:

- **main**: Full protection (approvals can be self-approved for now)
- **develop**: Full protection
- Consider using GitHub's "Allow specified actors to bypass required pull requests" temporarily

### For Small Team (2-3 Developers)

- **main**: Strict protection, require 1 approval (not from PR author)
- **develop**: Require 1 approval, can be from any team member
- Enable Code Owners for critical paths

### For Production Team (4+ Developers)

- **main**: Require 2 approvals, enforce Code Owners
- **develop**: Require 1 approval from Code Owners
- Enable signed commits
- Enable linear history
- Strict status check enforcement

---

## 🚨 Emergency Override Procedure

If you absolutely need to bypass protection (production emergency):

1. **Document the reason** in an issue
2. **Temporarily disable** branch protection
3. **Make the emergency fix**
4. **Re-enable** branch protection immediately
5. **Create post-mortem** document

⚠️ **Never** leave protection disabled longer than necessary!

---

## 📊 Branch Protection Summary

| Branch      | Direct Push | PR Required | Approvals | CI Required | Can Delete |
| ----------- | ----------- | ----------- | --------- | ----------- | ---------- |
| `main`      | ❌           | ✅           | 1+        | ✅           | ❌          |
| `develop`   | ❌           | ✅           | 1         | ✅           | ❌          |
| `release/*` | ✅           | ❌           | -         | -           | ✅          |
| `hotfix/*`  | ✅           | ❌           | -         | -           | ✅          |
| `feature/*` | ✅           | ❌           | -         | -           | ✅          |
| `bugfix/*`  | ✅           | ❌           | -         | -           | ✅          |

---

## 🔗 Useful Links

- GitHub Branch Protection: https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches
- Rulesets (New Feature): https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets
- Required Status Checks: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/collaborating-on-repositories-with-code-quality-features/about-status-checks

---

**Last Updated**: 22 November 2025
