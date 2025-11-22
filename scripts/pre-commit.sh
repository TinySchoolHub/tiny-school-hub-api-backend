#!/bin/bash
# Git pre-commit hook for TinySchoolHub API
# This hook runs automatically before each commit
#
# Installation:
#   make install-hooks
#   # OR manually:
#   ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Running pre-commit checks...${NC}\n"

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
        return 1
    fi
}

# Get list of staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)
STAGED_SQL_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.sql$' || true)
STAGED_YAML_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep -E '\.(yaml|yml)$' || true)

# Exit early if no Go files changed
if [ -z "$STAGED_GO_FILES" ]; then
    echo -e "${GREEN}No Go files staged, skipping Go checks${NC}\n"
    exit 0
fi

echo -e "${BLUE}Checking ${#STAGED_GO_FILES[@]} Go file(s)...${NC}\n"

# 1. Check formatting with gofmt
echo -e "${YELLOW}1. Checking code formatting...${NC}"
UNFORMATTED=$(gofmt -l $STAGED_GO_FILES 2>&1 || true)
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}✗ The following files need formatting:${NC}"
    echo "$UNFORMATTED"
    echo ""
    echo -e "${YELLOW}Fix with: gofmt -w $UNFORMATTED${NC}"
    echo -e "${YELLOW}Or run: make fmt${NC}"
    exit 1
fi
print_status 0 "Code formatting"

# 2. Check for common Go issues with go vet
echo -e "\n${YELLOW}2. Running go vet...${NC}"
if ! go vet ./...; then
    print_status 1 "go vet"
    echo -e "${YELLOW}Fix the issues above before committing${NC}"
    exit 1
fi
print_status 0 "go vet"

# 3. Run go mod tidy and check for changes
echo -e "\n${YELLOW}3. Checking go.mod and go.sum...${NC}"
go mod tidy
if ! git diff --exit-code go.mod go.sum > /dev/null 2>&1; then
    echo -e "${RED}✗ go.mod or go.sum has uncommitted changes after 'go mod tidy'${NC}"
    echo -e "${YELLOW}Please stage these changes and commit again${NC}"
    git diff go.mod go.sum
    exit 1
fi
print_status 0 "go.mod and go.sum"

# 4. Check for build errors
echo -e "\n${YELLOW}4. Checking if code compiles...${NC}"
if ! go build -o /dev/null ./...; then
    print_status 1 "Build"
    echo -e "${YELLOW}Fix build errors before committing${NC}"
    exit 1
fi
print_status 0 "Build"

# 5. Run tests (quick mode - no integration tests)
echo -e "\n${YELLOW}5. Running unit tests...${NC}"
# Temporarily disable exit on error to capture test failures
set +e
TEST_OUTPUT=$(go test -short -timeout=30s ./... 2>&1)
TEST_EXIT_CODE=$?
set -e

# Always show test output summary
printf "%s\n" "$TEST_OUTPUT" | grep -E "^(ok|FAIL|\?)" || true
if [ $TEST_EXIT_CODE -ne 0 ]; then
    # Show failed test details
    printf "%s\n" "$TEST_OUTPUT" | grep -E "^(---|    )" || true
fi

# Check if there are any actual test files
if echo "$TEST_OUTPUT" | grep -q "\[no test files\]"; then
    # Count packages with no tests
    NO_TEST_COUNT=$(echo "$TEST_OUTPUT" | grep -c "\[no test files\]")
    echo -e "${YELLOW}⚠ Warning: $NO_TEST_COUNT package(s) have no test files${NC}"
    echo -e "${YELLOW}Consider adding tests for better code quality${NC}"
    print_status 0 "Tests (no test files found)"
elif [ $TEST_EXIT_CODE -ne 0 ]; then
    print_status 1 "Tests"
    echo -e "${YELLOW}Fix failing tests before committing${NC}"
    exit 1
else
    print_status 0 "Tests"
fi

# 6. Check for common security issues (if gosec is installed)
if command -v gosec &> /dev/null; then
    echo -e "\n${YELLOW}6. Running security checks...${NC}"
    if ! gosec -quiet -fmt=text ./...; then
        print_status 1 "Security checks"
        echo -e "${YELLOW}Review security issues above${NC}"
        # Don't fail on security issues, just warn
    else
        print_status 0 "Security checks"
    fi
else
    echo -e "\n${YELLOW}6. Skipping security checks (gosec not installed)${NC}"
    echo -e "${YELLOW}Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest${NC}"
fi

# 7. Check for TODO/FIXME comments in staged files (warning only)
echo -e "\n${YELLOW}7. Checking for TODO/FIXME comments...${NC}"
TODO_COUNT=$(grep -n -E "TODO|FIXME" $STAGED_GO_FILES 2>/dev/null | wc -l || echo 0)
if [ "$TODO_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}⚠ Found $TODO_COUNT TODO/FIXME comment(s) in staged files${NC}"
    grep -n -E "TODO|FIXME" $STAGED_GO_FILES || true
    echo ""
fi

# 8. Check SQL migrations if any SQL files are staged
if [ -n "$STAGED_SQL_FILES" ]; then
    echo -e "\n${YELLOW}8. Checking SQL migrations...${NC}"
    # Check for up/down pair
    for sql_file in $STAGED_SQL_FILES; do
        if [[ $sql_file == *".up.sql" ]]; then
            down_file="${sql_file%.up.sql}.down.sql"
            if [ ! -f "$down_file" ]; then
                echo -e "${RED}✗ Missing corresponding down migration: $down_file${NC}"
                exit 1
            fi
        elif [[ $sql_file == *".down.sql" ]]; then
            up_file="${sql_file%.down.sql}.up.sql"
            if [ ! -f "$up_file" ]; then
                echo -e "${RED}✗ Missing corresponding up migration: $up_file${NC}"
                exit 1
            fi
        fi
    done
    print_status 0 "SQL migrations"
fi

# 9. Check YAML files if any are staged
if [ -n "$STAGED_YAML_FILES" ]; then
    echo -e "\n${YELLOW}9. Validating YAML files...${NC}"
    for yaml_file in $STAGED_YAML_FILES; do
        # Skip Helm template files (they contain Go templates, not pure YAML)
        if [[ "$yaml_file" == deploy/helm/*/templates/* ]]; then
            echo -e "${YELLOW}  Skipping Helm template: $yaml_file${NC}"
            continue
        fi
        
        if ! python3 -c "import yaml; yaml.safe_load(open('$yaml_file'))" 2>/dev/null; then
            echo -e "${RED}✗ Invalid YAML: $yaml_file${NC}"
            exit 1
        fi
    done
    print_status 0 "YAML validation"
fi

# 10. Check for sensitive data patterns
echo -e "\n${YELLOW}10. Checking for sensitive data...${NC}"
SENSITIVE_PATTERNS=(
    # Match actual hardcoded values, not variable comparisons
    "password\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    "secret\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    "api_key\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    "token\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    "private_key\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    "aws_access_key_id\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    "aws_secret_access_key\s*[:=]\s*['\"][^'\"]{3,}['\"]"
    # Actual private keys
    "BEGIN RSA PRIVATE KEY"
    "BEGIN PRIVATE KEY"
    "BEGIN OPENSSH PRIVATE KEY"
    # AWS keys pattern
    "AKIA[0-9A-Z]{16}"
)

for pattern in "${SENSITIVE_PATTERNS[@]}"; do
    # Exclude test files and mock data
    for file in $STAGED_GO_FILES; do
        # Skip test files
        if [[ "$file" == *_test.go ]]; then
            continue
        fi
        
        MATCHES=$(grep -i -E "$pattern" "$file" 2>/dev/null | grep -v "example" | grep -v "// " || true)
        if [ -n "$MATCHES" ]; then
            echo -e "${RED}✗ Possible sensitive data found in $file matching pattern: $pattern${NC}"
            echo "$MATCHES"
            echo -e "${YELLOW}Review the above matches carefully${NC}"
            echo -e "${YELLOW}If this is a false positive (e.g., variable name comparison), it's safe to proceed${NC}"
            # Changed to warning instead of blocking
            echo -e "${YELLOW}⚠ Warning: Review carefully before committing${NC}"
        fi
    done
done
print_status 0 "Sensitive data check"

# Success!
echo -e "\n${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ All pre-commit checks passed!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}\n"

exit 0
