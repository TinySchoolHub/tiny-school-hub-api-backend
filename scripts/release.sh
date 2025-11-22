#!/bin/bash

# Release script for Tiny School Hub API
# Usage: ./scripts/release.sh [major|minor|patch|VERSION]
# Examples:
#   ./scripts/release.sh patch    # Bump patch version
#   ./scripts/release.sh minor    # Bump minor version
#   ./scripts/release.sh major    # Bump major version
#   ./scripts/release.sh 1.2.3    # Set specific version

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get current version
CURRENT_VERSION=$(cat VERSION 2>/dev/null || echo "0.0.0")
echo -e "${BLUE}Current version: ${CURRENT_VERSION}${NC}"

# Parse version
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Determine new version
if [ "$1" == "major" ]; then
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
elif [ "$1" == "minor" ]; then
    MINOR=$((MINOR + 1))
    PATCH=0
    NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
elif [ "$1" == "patch" ]; then
    PATCH=$((PATCH + 1))
    NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
elif [ -n "$1" ]; then
    # Validate version format
    if [[ ! "$1" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
        echo -e "${RED}Error: Invalid version format. Use X.Y.Z or X.Y.Z-suffix${NC}"
        exit 1
    fi
    NEW_VERSION="$1"
else
    echo -e "${RED}Error: Please specify version bump type (major/minor/patch) or exact version${NC}"
    echo "Usage: $0 [major|minor|patch|X.Y.Z]"
    exit 1
fi

echo -e "${YELLOW}New version: ${NEW_VERSION}${NC}"

# Confirm
read -p "Create release v${NEW_VERSION}? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Release cancelled${NC}"
    exit 0
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: You have uncommitted changes${NC}"
    echo "Please commit or stash them before creating a release"
    exit 1
fi

# Check if on main or release branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "main" && ! "$CURRENT_BRANCH" =~ ^release/ ]]; then
    echo -e "${YELLOW}Warning: You are on branch '${CURRENT_BRANCH}'${NC}"
    echo -e "${YELLOW}Releases should typically be created from 'main' or 'release/*' branches${NC}"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Check if tag already exists
if git rev-parse "v${NEW_VERSION}" >/dev/null 2>&1; then
    echo -e "${RED}Error: Tag v${NEW_VERSION} already exists${NC}"
    exit 1
fi

echo -e "\n${BLUE}Running pre-release checks...${NC}"

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
if ! go test ./...; then
    echo -e "${RED}Tests failed! Fix them before releasing${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Tests passed${NC}"

# Run pre-commit checks
if [ -f "scripts/pre-commit.sh" ]; then
    echo -e "${YELLOW}Running pre-commit checks...${NC}"
    if ! bash scripts/pre-commit.sh; then
        echo -e "${RED}Pre-commit checks failed!${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Pre-commit checks passed${NC}"
fi

# Update VERSION file
echo "$NEW_VERSION" > VERSION
echo -e "${GREEN}✓ Updated VERSION file${NC}"

# Update version in main.go if it exists
if grep -q "version.*=" cmd/api/main.go 2>/dev/null; then
    sed -i.bak "s/version = \".*\"/version = \"${NEW_VERSION}\"/" cmd/api/main.go
    rm -f cmd/api/main.go.bak
    echo -e "${GREEN}✓ Updated version in main.go${NC}"
fi

# Generate changelog automatically
echo -e "\n${YELLOW}Generating changelog from git commits...${NC}"
if [ -f "scripts/changelog.sh" ]; then
    # Generate changelog and save to temp file
    ./scripts/changelog.sh > /tmp/changelog_entry.md 2>/dev/null || true
    
    # Show generated changelog
    if [ -f /tmp/changelog_entry.md ] && [ -s /tmp/changelog_entry.md ]; then
        echo -e "${GREEN}✓ Changelog generated!${NC}"
        echo
        cat /tmp/changelog_entry.md
        echo
        
        read -p "Would you like to automatically insert this into CHANGELOG.md? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            # Extract just the changelog section (skip the colored output and instructions)
            # Use sed to remove last line instead of head -n -1 (which doesn't work on macOS)
            CHANGELOG_CONTENT=$(./scripts/changelog.sh 2>/dev/null | sed -n '/^## \[VERSION\]/,/^═══/p' | sed '$d')
            
            # Replace [VERSION] with actual version
            CHANGELOG_CONTENT=$(echo "$CHANGELOG_CONTENT" | sed "s/\[VERSION\]/[${NEW_VERSION}]/g")
            
            # Create temporary file with new content
            {
                # Keep everything up to and including [Unreleased] section
                sed -n '1,/## \[Unreleased\]/p' CHANGELOG.md
                echo
                # Add new version
                echo "$CHANGELOG_CONTENT"
                echo
                # Add rest of file after [Unreleased] section
                sed '1,/## \[Unreleased\]/d' CHANGELOG.md
            } > /tmp/changelog_new.md
            
            mv /tmp/changelog_new.md CHANGELOG.md
            echo -e "${GREEN}✓ CHANGELOG.md updated automatically${NC}"
        else
            echo -e "${YELLOW}Please update CHANGELOG.md manually${NC}"
            read -p "Press Enter when done..."
        fi
    else
        echo -e "${YELLOW}Could not generate changelog automatically${NC}"
        echo "Please update CHANGELOG.md manually"
        read -p "Press Enter when done..."
    fi
else
    # Fallback to manual update
    echo -e "\n${YELLOW}Please update CHANGELOG.md before continuing${NC}"
    echo "Add release notes for v${NEW_VERSION}"
    read -p "Press Enter when ready to continue..."
fi

# Check if CHANGELOG was updated
if ! grep -q "\[${NEW_VERSION}\]" CHANGELOG.md; then
    echo -e "${YELLOW}Warning: CHANGELOG.md doesn't contain [${NEW_VERSION}]${NC}"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        # Revert VERSION file
        echo "$CURRENT_VERSION" > VERSION
        exit 0
    fi
fi

# Commit version bump
git add VERSION CHANGELOG.md cmd/api/main.go 2>/dev/null || git add VERSION CHANGELOG.md
git commit -m "chore: bump version to v${NEW_VERSION}"
echo -e "${GREEN}✓ Committed version bump${NC}"

# Create annotated tag
echo -e "\n${YELLOW}Creating git tag v${NEW_VERSION}...${NC}"

# Extract changelog for this version
CHANGELOG_ENTRY=$(sed -n "/## \[${NEW_VERSION}\]/,/## \[/p" CHANGELOG.md | sed '$d')

if [ -z "$CHANGELOG_ENTRY" ]; then
    CHANGELOG_ENTRY="Release v${NEW_VERSION}"
fi

git tag -a "v${NEW_VERSION}" -m "${CHANGELOG_ENTRY}"
echo -e "${GREEN}✓ Created tag v${NEW_VERSION}${NC}"

# Push changes and tag
echo -e "\n${YELLOW}Pushing to remote...${NC}"
git push origin "$CURRENT_BRANCH"
git push origin "v${NEW_VERSION}"
echo -e "${GREEN}✓ Pushed changes and tag${NC}"

# Success
echo -e "\n${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ Release v${NEW_VERSION} created successfully!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "\nNext steps:"
echo -e "  1. GitHub Actions will automatically build and publish Docker images"
echo -e "  2. Create GitHub Release at: https://github.com/TinySchoolHub/tiny-school-hub-api-backend/releases/new?tag=v${NEW_VERSION}"
echo -e "  3. Deploy to staging/production if needed"
echo -e "\nDocker image will be available at:"
echo -e "  ${BLUE}ghcr.io/tinyschoolhub/tiny-school-hub-api:v${NEW_VERSION}${NC}"
echo -e "  ${BLUE}ghcr.io/tinyschoolhub/tiny-school-hub-api:latest${NC}"

exit 0
