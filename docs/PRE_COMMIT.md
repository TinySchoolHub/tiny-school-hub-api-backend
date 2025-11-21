# Pre-commit Setup Guide

This project uses pre-commit hooks to ensure code quality before committing.

## 🚀 Quick Setup

### Option 1: Native Git Hooks (Recommended)

```bash
# Install the git hook
make install-hooks

# The hook will now run automatically on every commit
git commit -m "your message"

# To bypass the hook if needed
git commit --no-verify -m "your message"
```

### Option 2: Python pre-commit Framework (Advanced)

```bash
# Install pre-commit framework
pip install pre-commit

# Install the hooks
make install-pre-commit-framework

# Run manually on all files
pre-commit run --all-files
```

## 🔍 What Gets Checked

### Go Checks
- ✅ **Code Formatting** - Ensures `gofmt` compliance
- ✅ **go vet** - Checks for suspicious constructs
- ✅ **go mod tidy** - Ensures dependencies are clean
- ✅ **Build** - Verifies code compiles
- ✅ **Unit Tests** - Runs fast unit tests
- ✅ **Security** - Runs gosec if installed
- ⚠️ **TODO/FIXME** - Warns about pending items

### File Checks
- ✅ **SQL Migrations** - Validates up/down pairs
- ✅ **YAML Validation** - Checks YAML syntax
- ✅ **Sensitive Data** - Prevents committing secrets
- ✅ **Large Files** - Blocks files >1MB
- ✅ **Merge Conflicts** - Detects conflict markers

### Code Quality (via golangci-lint)
- Cyclomatic complexity
- Unused code
- Error handling
- Common mistakes
- Performance issues
- Style violations

## 📋 Manual Commands

```bash
# Run all pre-commit checks manually
make pre-commit

# Run specific checks
make fmt          # Format code
make vet          # Run go vet
make lint         # Run golangci-lint
make test         # Run tests

# Install development tools
make install-tools  # Installs golangci-lint, migrate, etc.
```

## 🛠️ Installation Requirements

### Required
- Go 1.24+
- git

### Optional (for enhanced checks)
```bash
# Install linting tools
make install-tools

# Install gosec for security scanning
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Install pre-commit framework (Python)
pip install pre-commit
```

## ⚙️ Configuration Files

- **`.pre-commit-config.yaml`** - Python pre-commit framework config
- **`.golangci.yml`** - Go linter configuration
- **`scripts/pre-commit.sh`** - Native git hook script

## 🔧 Customization

### Skip Specific Checks

```bash
# Bypass all hooks for emergency commits
git commit --no-verify -m "emergency fix"

# Skip specific linters
SKIP=golangci-lint git commit -m "message"
```

### Adjust Linter Settings

Edit `.golangci.yml` to configure:
- Line length limits
- Cyclomatic complexity thresholds
- Enabled/disabled linters
- Per-file exclusions

### Modify Hook Behavior

Edit `scripts/pre-commit.sh` to:
- Add/remove checks
- Change error vs warning behavior
- Adjust timeouts

## 📊 CI Integration

Pre-commit checks are also recommended for CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run pre-commit
  run: |
    make install-tools
    make pre-commit
```

## 🐛 Troubleshooting

### Hook not running
```bash
# Check hook is installed
ls -l .git/hooks/pre-commit

# Reinstall
make uninstall-hooks
make install-hooks
```

### Hook fails with permission denied
```bash
# Make script executable
chmod +x scripts/pre-commit.sh
```

### golangci-lint not found
```bash
# Install linting tools
make install-tools
```

### Python pre-commit issues
```bash
# Clean and reinstall
pre-commit clean
pre-commit install
```

## 📚 Best Practices

1. **Fix issues locally** before committing
2. **Run `make pre-commit`** before pushing
3. **Keep hooks updated** with project needs
4. **Document bypass reasons** when using `--no-verify`
5. **Review warnings** even if commit succeeds

## 🔗 Resources

- [golangci-lint docs](https://golangci-lint.run/)
- [pre-commit framework](https://pre-commit.com/)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

---

**Need help?** Check the Makefile targets with `make help`
