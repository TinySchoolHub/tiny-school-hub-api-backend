#!/bin/bash

# Changelog generator for Tiny School Hub API
# Generates changelog entries from git commits using Conventional Commits
# Usage: ./scripts/changelog.sh [from_tag] [to_tag]
# Examples:
#   ./scripts/changelog.sh                    # Changes since last tag
#   ./scripts/changelog.sh v1.0.0 HEAD        # Changes from v1.0.0 to now
#   ./scripts/changelog.sh v1.0.0 v1.1.0      # Changes between two tags

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Get the latest tag
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

# Determine range
if [ -n "$1" ] && [ -n "$2" ]; then
    FROM_REF="$1"
    TO_REF="$2"
elif [ -n "$1" ]; then
    FROM_REF="$1"
    TO_REF="HEAD"
elif [ -n "$LATEST_TAG" ]; then
    FROM_REF="$LATEST_TAG"
    TO_REF="HEAD"
else
    # No tags, get all commits
    FROM_REF=$(git rev-list --max-parents=0 HEAD)
    TO_REF="HEAD"
fi

echo -e "${BLUE}Generating changelog from ${CYAN}${FROM_REF}${BLUE} to ${CYAN}${TO_REF}${NC}"
echo

# Arrays to store commits by type
declare -a FEAT_COMMITS
declare -a FIX_COMMITS
declare -a DOCS_COMMITS
declare -a STYLE_COMMITS
declare -a REFACTOR_COMMITS
declare -a PERF_COMMITS
declare -a TEST_COMMITS
declare -a CHORE_COMMITS
declare -a BUILD_COMMITS
declare -a CI_COMMITS
declare -a BREAKING_COMMITS
declare -a OTHER_COMMITS

# Get commits
while IFS= read -r commit; do
    # Get commit message (first line only)
    MESSAGE=$(git log --format=%s -n 1 "$commit")
    
    # Check for breaking changes
    BREAKING=$(git log --format=%b -n 1 "$commit" | grep -i "BREAKING CHANGE" || true)
    
    if [ -n "$BREAKING" ]; then
        BREAKING_COMMITS+=("$MESSAGE")
        continue
    fi
    
    # Categorize by conventional commit type
    if [[ "$MESSAGE" =~ ^feat(\(.+\))?: ]]; then
        FEAT_COMMITS+=("${MESSAGE#feat*: }")
    elif [[ "$MESSAGE" =~ ^fix(\(.+\))?: ]]; then
        FIX_COMMITS+=("${MESSAGE#fix*: }")
    elif [[ "$MESSAGE" =~ ^docs(\(.+\))?: ]]; then
        DOCS_COMMITS+=("${MESSAGE#docs*: }")
    elif [[ "$MESSAGE" =~ ^style(\(.+\))?: ]]; then
        STYLE_COMMITS+=("${MESSAGE#style*: }")
    elif [[ "$MESSAGE" =~ ^refactor(\(.+\))?: ]]; then
        REFACTOR_COMMITS+=("${MESSAGE#refactor*: }")
    elif [[ "$MESSAGE" =~ ^perf(\(.+\))?: ]]; then
        PERF_COMMITS+=("${MESSAGE#perf*: }")
    elif [[ "$MESSAGE" =~ ^test(\(.+\))?: ]]; then
        TEST_COMMITS+=("${MESSAGE#test*: }")
    elif [[ "$MESSAGE" =~ ^chore(\(.+\))?: ]]; then
        CHORE_COMMITS+=("${MESSAGE#chore*: }")
    elif [[ "$MESSAGE" =~ ^build(\(.+\))?: ]]; then
        BUILD_COMMITS+=("${MESSAGE#build*: }")
    elif [[ "$MESSAGE" =~ ^ci(\(.+\))?: ]]; then
        CI_COMMITS+=("${MESSAGE#ci*: }")
    else
        OTHER_COMMITS+=("$MESSAGE")
    fi
done < <(git rev-list --reverse "${FROM_REF}..${TO_REF}")

# Get current date
CURRENT_DATE=$(date +%Y-%m-%d)

# Generate changelog
echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Generated Changelog Entry${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo
echo "## [VERSION] - $CURRENT_DATE"
echo

# Breaking Changes (most important)
if [ ${#BREAKING_COMMITS[@]} -gt 0 ]; then
    echo "### ⚠️ BREAKING CHANGES"
    echo
    for commit in "${BREAKING_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

# Added (new features)
if [ ${#FEAT_COMMITS[@]} -gt 0 ]; then
    echo "### Added"
    echo
    for commit in "${FEAT_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

# Fixed (bug fixes)
if [ ${#FIX_COMMITS[@]} -gt 0 ]; then
    echo "### Fixed"
    echo
    for commit in "${FIX_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

# Changed (refactoring, performance)
if [ ${#REFACTOR_COMMITS[@]} -gt 0 ] || [ ${#PERF_COMMITS[@]} -gt 0 ]; then
    echo "### Changed"
    echo
    for commit in "${REFACTOR_COMMITS[@]}"; do
        echo "- $commit"
    done
    for commit in "${PERF_COMMITS[@]}"; do
        echo "- $commit (performance)"
    done
    echo
fi

# Documentation
if [ ${#DOCS_COMMITS[@]} -gt 0 ]; then
    echo "### Documentation"
    echo
    for commit in "${DOCS_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

# Tests
if [ ${#TEST_COMMITS[@]} -gt 0 ]; then
    echo "### Tests"
    echo
    for commit in "${TEST_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

# Build/CI
if [ ${#BUILD_COMMITS[@]} -gt 0 ] || [ ${#CI_COMMITS[@]} -gt 0 ]; then
    echo "### Build/CI"
    echo
    for commit in "${BUILD_COMMITS[@]}"; do
        echo "- $commit"
    done
    for commit in "${CI_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

# Other changes
if [ ${#CHORE_COMMITS[@]} -gt 0 ] || [ ${#STYLE_COMMITS[@]} -gt 0 ] || [ ${#OTHER_COMMITS[@]} -gt 0 ]; then
    echo "### Other"
    echo
    for commit in "${CHORE_COMMITS[@]}"; do
        echo "- $commit"
    done
    for commit in "${STYLE_COMMITS[@]}"; do
        echo "- $commit"
    done
    for commit in "${OTHER_COMMITS[@]}"; do
        echo "- $commit"
    done
    echo
fi

echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo
echo -e "${YELLOW}Instructions:${NC}"
echo "1. Copy the above content"
echo "2. Replace [VERSION] with your actual version (e.g., 1.2.3)"
echo "3. Paste it at the top of CHANGELOG.md under the '## [Unreleased]' section"
echo "4. Review and edit as needed"
echo "5. Remove any irrelevant entries (e.g., WIP commits)"
echo
echo -e "${CYAN}Tip: You can also pipe to clipboard:${NC}"
echo "  ./scripts/changelog.sh | pbcopy    # macOS"
echo "  ./scripts/changelog.sh | xclip     # Linux"
echo
echo -e "${CYAN}Or append directly to CHANGELOG.md:${NC}"
echo "  ./scripts/changelog.sh >> CHANGELOG_TEMP.md"
echo

# Statistics
TOTAL_COMMITS=$(git rev-list --count "${FROM_REF}..${TO_REF}")
echo -e "${BLUE}Statistics:${NC}"
echo "  Total commits: $TOTAL_COMMITS"
echo "  Features: ${#FEAT_COMMITS[@]}"
echo "  Fixes: ${#FIX_COMMITS[@]}"
echo "  Breaking: ${#BREAKING_COMMITS[@]}"
echo

exit 0
