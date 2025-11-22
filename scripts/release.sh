#!/bin/bash

# GitFlow Setup Script for Tiny School Hub API
# This script initializes the GitFlow workflow structure

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}   GitFlow Setup for Tiny School Hub    ${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo

# Check if we're in a git repository
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${BLUE}Current branch: ${CURRENT_BRANCH}${NC}"

# Check if we're on main
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}Warning: You are not on the main branch${NC}"
    read -p "Do you want to continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: You have uncommitted changes${NC}"
    echo "Please commit or stash them before running this script"
    exit 1
fi

echo -e "\n${YELLOW}Step 1: Creating develop branch${NC}"

# Check if develop branch already exists
if git show-ref --verify --quiet refs/heads/develop; then
    echo -e "${YELLOW}develop branch already exists locally${NC}"
    git checkout develop
    git pull origin develop 2>/dev/null || echo "No remote develop branch yet"
else
    # Check if develop exists on remote
    if git ls-remote --heads origin develop | grep -q develop; then
        echo -e "${YELLOW}develop branch exists on remote, checking out${NC}"
        git checkout -b develop origin/develop
    else
        echo -e "${GREEN}Creating new develop branch from main${NC}"
        git checkout -b develop
    fi
fi

echo -e "${GREEN}✓ develop branch ready${NC}"

echo -e "\n${YELLOW}Step 2: Pushing develop to remote${NC}"

if git ls-remote --heads origin develop | grep -q develop; then
    echo -e "${YELLOW}develop branch already exists on remote${NC}"
else
    echo -e "${GREEN}Pushing develop branch to remote${NC}"
    git push -u origin develop
    echo -e "${GREEN}✓ develop branch pushed${NC}"
fi

echo -e "\n${YELLOW}Step 3: Branch protection setup${NC}"
echo -e "${BLUE}Please configure branch protection rules on GitHub:${NC}"
echo
echo -e "${YELLOW}For 'main' branch:${NC}"
echo "  1. Go to: Settings > Branches > Add rule"
echo "  2. Branch name pattern: main"
echo "  3. Enable: ✅ Require a pull request before merging"
echo "  4. Enable: ✅ Require approvals (2)"
echo "  5. Enable: ✅ Dismiss stale pull request approvals when new commits are pushed"
echo "  6. Enable: ✅ Require review from Code Owners"
echo "  7. Enable: ✅ Require status checks to pass before merging"
echo "  8. Enable: ✅ Require branches to be up to date before merging"
echo "  9. Enable: ✅ Require conversation resolution before merging"
echo "  10. Enable: ✅ Require linear history"
echo "  11. Enable: ✅ Do not allow bypassing the above settings"
echo "  12. Enable: ✅ Restrict who can push to matching branches"
echo
echo -e "${YELLOW}For 'develop' branch:${NC}"
echo "  1. Go to: Settings > Branches > Add rule"
echo "  2. Branch name pattern: develop"
echo "  3. Enable: ✅ Require a pull request before merging"
echo "  4. Enable: ✅ Require approvals (1)"
echo "  5. Enable: ✅ Require status checks to pass before merging"
echo "  6. Enable: ✅ Require conversation resolution before merging"
echo "  7. Enable: ✅ Allow force pushes (for maintainers only)"
echo
read -p "Press Enter after configuring branch protection rules..."

echo -e "\n${YELLOW}Step 4: Renovate setup${NC}"
echo -e "${BLUE}To enable Renovate:${NC}"
echo "  1. Visit: https://github.com/apps/renovate"
echo "  2. Click 'Install' or 'Configure'"
echo "  3. Select repository: TinySchoolHub/tiny-school-hub-api-backend"
echo "  4. Grant permissions"
echo
echo "  Renovate will:"
echo "  - Create PRs for dependency updates to 'develop' branch"
echo "  - Auto-merge patch updates with 'automerge' label"
echo "  - Group Go dependencies together"
echo "  - Run weekly on Monday mornings"
echo
read -p "Press Enter after setting up Renovate..."

echo -e "\n${YELLOW}Step 5: Creating .github/CODEOWNERS${NC}"

if [ ! -f ".github/CODEOWNERS" ]; then
    mkdir -p .github
    cat > .github/CODEOWNERS << 'EOF'
# Code Owners for Tiny School Hub API

# Default owners for everything
* @fabienchevalier

# API and handlers
/internal/http/ @fabienchevalier

# Core business logic
/internal/core/ @fabienchevalier

# Database and migrations
/internal/repository/ @fabienchevalier
/migrations/ @fabienchevalier

# Configuration and deployment
/deploy/ @fabienchevalier
/docker-compose.yml @fabienchevalier
/Dockerfile @fabienchevalier

# Documentation
/docs/ @fabienchevalier
*.md @fabienchevalier

# CI/CD
/.github/ @fabienchevalier
EOF
    echo -e "${GREEN}✓ Created .github/CODEOWNERS${NC}"
    git add .github/CODEOWNERS
else
    echo -e "${YELLOW}.github/CODEOWNERS already exists${NC}"
fi

echo -e "\n${YELLOW}Step 6: Creating .github/pull_request_template.md${NC}"

if [ ! -f ".github/pull_request_template.md" ]; then
    cat > .github/pull_request_template.md << 'EOF'
## Description

<!-- Provide a brief description of the changes in this PR -->

## Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📝 Documentation update
- [ ] 🔧 Configuration change
- [ ] 🧪 Test update
- [ ] ♻️ Refactoring (no functional changes)

## Related Issues

<!-- Link to related issues using #issue_number -->

Fixes #
Related to #

## Changes Made

<!-- List the main changes made in this PR -->

-
-
-

## Testing

<!-- Describe how you tested these changes -->

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] All tests pass locally

### Test Commands

```bash
# Commands used to test
go test ./...
```

## Screenshots (if applicable)

<!-- Add screenshots for UI changes -->

## Checklist

<!-- Mark completed items with an "x" -->

- [ ] Code follows the project's style guidelines
- [ ] Self-review of code performed
- [ ] Comments added for complex logic
- [ ] Documentation updated (if needed)
- [ ] No new warnings generated
- [ ] Tests added that prove the fix/feature works
- [ ] Dependent changes merged and published
- [ ] CHANGELOG.md updated (for features/fixes)

## Deployment Notes

<!-- Any special considerations for deployment? -->

## Rollback Plan

<!-- How can this change be rolled back if needed? -->

---

**For Reviewers:**
- Code quality and best practices
- Test coverage
- Security considerations
- Performance implications
EOF
    echo -e "${GREEN}✓ Created .github/pull_request_template.md${NC}"
    git add .github/pull_request_template.md
else
    echo -e "${YELLOW}.github/pull_request_template.md already exists${NC}"
fi

echo -e "\n${YELLOW}Step 7: Creating CONTRIBUTING.md updates${NC}"

if [ -f "CONTRIBUTING.md" ]; then
    if ! grep -q "GitFlow" CONTRIBUTING.md; then
        cat >> CONTRIBUTING.md << 'EOF'

## GitFlow Workflow

This project uses GitFlow. Please read [docs/GITFLOW.md](docs/GITFLOW.md) for detailed instructions.

### Quick Start

1. **Feature development**: Branch from `develop`
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/my-feature
   ```

2. **Create Pull Request**: `feature/my-feature` → `develop`

3. **After approval**: Merge to `develop`, delete feature branch

4. **Release**: Create `release/vX.Y.Z` from `develop`, then merge to `main`

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add new feature`
- `fix: resolve bug`
- `docs: update documentation`
- `chore: update dependencies`
- `refactor: improve code structure`
- `test: add tests`
EOF
        echo -e "${GREEN}✓ Updated CONTRIBUTING.md${NC}"
        git add CONTRIBUTING.md
    else
        echo -e "${YELLOW}CONTRIBUTING.md already mentions GitFlow${NC}"
    fi
fi

echo -e "\n${YELLOW}Step 8: Committing GitFlow setup${NC}"

if git diff --cached --quiet; then
    echo -e "${YELLOW}No changes to commit${NC}"
else
    git commit -m "chore: setup GitFlow workflow with Renovate

- Add GitFlow documentation
- Configure Renovate for dependency management
- Add PR workflows for develop and main branches
- Add CODEOWNERS and PR template
    - Update CONTRIBUTING.md with GitFlow instructions"
    
    echo -e "${GREEN}✓ Changes committed${NC}"
    
    read -p "Push changes to remote? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git push origin develop
        echo -e "${GREEN}✓ Changes pushed to develop${NC}"
    fi
fi

echo -e "\n${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}   ✓ GitFlow Setup Complete!            ${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo
echo -e "${BLUE}Next Steps:${NC}"
echo "  1. Configure branch protection on GitHub"
echo "  2. Install Renovate app on your repository"
echo "  3. Review docs/GITFLOW.md for workflow details"
echo "  4. Start creating feature branches from develop"
echo "  5. Create your first PR to develop"
echo
echo -e "${YELLOW}Helpful Commands:${NC}"
echo "  • Create feature: git checkout -b feature/my-feature develop"
echo "  • Create bugfix:  git checkout -b bugfix/fix-bug develop"
echo "  • Create release: git checkout -b release/v1.2.0 develop"
echo "  • Create hotfix:  git checkout -b hotfix/v1.2.1 main"
echo
echo -e "${BLUE}Documentation:${NC}"
echo "  • GitFlow Guide:  docs/GITFLOW.md"
echo "  • Release Guide:  docs/RELEASE_GUIDE.md"
echo "  • Renovate Config: renovate.json"
echo
echo -e "${GREEN}Happy coding! 🚀${NC}"

exit 0
