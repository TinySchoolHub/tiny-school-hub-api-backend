# Pre-commit Checks

This document describes the pre-commit hooks configured for this project.

## Overview

Pre-commit hooks automatically run checks before each commit to ensure code quality, prevent common mistakes, and maintain consistent standards across the codebase.

## Installation

### Quick Install (Recommended)
```bash
make install-hooks
```

This installs a native Git hook that runs automatically on every commit.

## What Gets Checked

### 1. Code Formatting
- Ensures all Go files are formatted with `gofmt`
- Auto-fix: `make fmt`

### 2. Go Vet
- Checks for suspicious constructs and common mistakes
- Reports issues that may cause bugs

### 3. Module Dependencies
- Runs `go mod tidy` and checks for uncommitted changes
- Ensures `go.mod` and `go.sum` are in sync

### 4. Build Validation
- Verifies the code compiles successfully
- Catches syntax and type errors

### 5. Unit Tests
- Runs fast unit tests (with `-short` flag)
- 30-second timeout for quick feedback

### 6. Security Checks (Optional)
- Runs `gosec` if installed
- Scans for common security issues
- Install: `go install github.com/securego/gosec/v2/cmd/gosec@latest`

### 7. TODO/FIXME Detection
- Warns about pending TODO/FIXME comments
- Non-blocking, for awareness only

### 8. SQL Migration Validation
- Ensures every `.up.sql` has a corresponding `.down.sql`
- Maintains migration consistency

### 9. YAML Validation
- Validates YAML syntax in configuration files
- Prevents invalid Kubernetes/Helm configs

### 10. Sensitive Data Detection
- Scans for potential secrets, passwords, API keys
- Blocks commits containing sensitive patterns

## Manual Execution

Run checks without committing:
```bash
make pre-commit
```

## Bypassing Hooks

For emergency commits only:
```bash
git commit --no-verify -m "emergency fix"
```

**Note:** Use sparingly and always fix issues in follow-up commits.

## Advanced Setup: Python pre-commit Framework

For additional checks (SQL linting, Dockerfile linting, etc.):

```bash
# Install Python pre-commit
pip install pre-commit

# Setup hooks
make install-pre-commit-framework

# Run all checks
pre-commit run --all-files
```

## Configuration Files

- `scripts/pre-commit.sh` - Main pre-commit script
- `.pre-commit-config.yaml` - Python pre-commit framework config
- `.golangci.yml` - Go linter configuration

## Troubleshooting

### Hook not executing
```bash
# Verify installation
ls -l .git/hooks/pre-commit

# Reinstall
make uninstall-hooks
make install-hooks
```

### Permission denied
```bash
chmod +x scripts/pre-commit.sh
```

### Missing tools
```bash
# Install all development tools
make install-tools
```

## CI/CD Integration

These checks should also run in CI/CD pipelines to ensure consistency.

Example GitHub Actions workflow:
```yaml
- name: Run pre-commit checks
  run: |
    make install-tools
    make pre-commit
```

## Best Practices

1. **Fix locally first** - Don't bypass hooks to "fix in CI"
2. **Keep commits clean** - One logical change per commit
3. **Review warnings** - Even non-blocking warnings matter
4. **Update hooks** - Keep tooling up to date
5. **Document bypasses** - Explain why if you must use `--no-verify`

---

For detailed setup instructions, see `docs/PRE_COMMIT.md`
