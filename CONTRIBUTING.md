# Contributing to Tiny School Hub API

Thank you for your interest in contributing! 🎉

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/tiny-school-hub-api-backend.git`
3. Create a branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit: `git commit -m "feat: add new feature"`
7. Push: `git push origin feature/my-feature`
8. Open a Pull Request

## Development Setup

See [README.md](README.md#development) for setup instructions.

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks
- `perf:` Performance improvements
- `ci:` CI/CD changes

Examples:
```
feat: add user profile endpoint
fix: resolve JWT token expiration issue
docs: update API documentation
test: add unit tests for auth middleware
```

## Code Style

- Run `gofmt` before committing
- Follow Go best practices
- Write tests for new features
- Keep functions small and focused
- Add comments for complex logic

## Testing

- Write unit tests for all new code
- Aim for >80% test coverage
- Use table-driven tests
- Test both success and error cases

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/core/auth
```

## Pull Request Process

1. **Update tests** - Add tests for new functionality
2. **Update documentation** - Update README.md, API docs, etc.
3. **Run pre-commit checks** - Ensure all checks pass
4. **Update CHANGELOG.md** - Add your changes under `[Unreleased]`
5. **Request review** - Tag maintainers for review
6. **Address feedback** - Make requested changes
7. **Squash commits** - Clean up commit history before merge

## Pre-commit Hooks

Install pre-commit hooks:

```bash
cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

The hooks will:
- Check Go formatting
- Run linter
- Run unit tests
- Check for sensitive data
- Validate YAML files

## Branching Strategy

- `main` - Production-ready code
- `develop` - Development branch (if used)
- `feature/*` - New features
- `fix/*` - Bug fixes
- `hotfix/*` - Urgent production fixes
- `release/*` - Release preparation

## Release Process

See [docs/VERSIONING.md](docs/VERSIONING.md) for release process.

## Code Review Guidelines

### For Authors
- Keep PRs small and focused
- Write clear descriptions
- Add screenshots/examples if applicable
- Respond to feedback promptly

### For Reviewers
- Be constructive and respectful
- Test the changes locally
- Check for security issues
- Verify test coverage

## Reporting Bugs

Use GitHub Issues with the bug template:

**Title**: Short, descriptive title

**Description**:
- Expected behavior
- Actual behavior
- Steps to reproduce
- Environment (OS, Go version, etc.)
- Error messages/logs

## Suggesting Features

Use GitHub Issues with the feature template:

**Title**: Short feature description

**Description**:
- Problem statement
- Proposed solution
- Alternative solutions considered
- Additional context

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone.

### Our Standards

- Be respectful and inclusive
- Accept constructive criticism
- Focus on what's best for the community
- Show empathy towards others

### Unacceptable Behavior

- Harassment or discrimination
- Trolling or insulting comments
- Public or private harassment
- Publishing others' private information

## Questions?

- Open a GitHub Discussion
- Join our community chat (if available)
- Email: [your-email@example.com]

## License

By contributing, you agree that your contributions will be licensed under the project's license.

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes
- Project README (for significant contributions)

Thank you for contributing! 🚀

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
