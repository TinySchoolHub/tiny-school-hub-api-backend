# Tiny School Hub API Backend

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![AWS SDK](https://img.shields.io/badge/AWS_SDK-v2-FF9900?style=flat&logo=amazon-aws)](https://aws.github.io/aws-sdk-go-v2/)
[![License](https://img.shields.io/badge/License-Private-red.svg)](LICENSE)

> A production-ready, Kubernetes-native backend for school communication and management.

## 🎯 Overview

Tiny School Hub is a **vendor-neutral**, **cloud-agnostic** backend service built with Go, designed to run on any Kubernetes distribution. It provides secure communication between teachers and parents with class management, photo sharing, absence tracking, and messaging.

### Key Features

- 🔐 **Security-First:** Argon2id password hashing, JWT authentication, RBAC
- ☸️ **Kubernetes-Native:** Health probes, graceful shutdown, HPA, NetworkPolicies
- 🌐 **Vendor-Neutral:** No cloud-specific SDKs, works with any S3-compatible storage
- 📦 **12-Factor App:** Configuration via environment, stateless, portable
- 🚀 **Production-Ready:** AWS SDK v2, structured logging, rate limiting
- 🧪 **Testable:** Clean architecture, dependency injection, interfaces

## 📋 Quick Start

```bash
# 1. Clone and setup
git clone https://github.com/TinySchoolHub/tiny-school-hub-api-backend.git
cd tiny-school-hub-api-backend
cp .env.example .env

# 2. Install development tools and pre-commit hooks
make install-tools
make install-hooks

# 3. Start local infrastructure (PostgreSQL + MinIO)
make docker-up

# 4. Run database migrations
make migrate-up

# 5. Start the API server
make run
# Server starts on http://localhost:8080
```

## 🏗️ Architecture

```
tiny-school-hub-api-backend/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── core/
│   │   ├── auth/         # JWT & password hashing (Argon2id)
│   │   └── domain/       # Domain models & errors
│   ├── http/
│   │   ├── handlers/     # HTTP request handlers
│   │   └── middleware/   # Auth, CORS, rate limiting
│   ├── repository/       # Data persistence interfaces
│   │   └── postgres/     # PostgreSQL implementations
│   └── storage/          # S3-compatible storage (AWS SDK v2)
├── pkg/log/              # Structured logging
├── migrations/           # Database migrations
└── deploy/
    ├── helm/             # Helm charts
    └── kustomize/        # Kustomize overlays
```

### Technology Stack

- **Language:** Go 1.24+
- **Router:** Chi (lightweight, idiomatic)
- **Database:** PostgreSQL 16+
- **Storage:** S3-compatible (MinIO, AWS S3, etc.)
- **Auth:** JWT (golang-jwt/jwt), Argon2id
- **Logging:** Zerolog (structured JSON)
- **Migrations:** golang-migrate
- **Container:** Docker (multi-stage, scratch-based)
- **Orchestration:** Kubernetes + Helm

## 🔐 Security

### Authentication & Authorization
- **Password Hashing:** Argon2id (memory-hard, GPU-resistant)
- **JWT Tokens:** 
  - Short-lived access tokens (15 minutes)
  - Long-lived refresh tokens (7 days) with rotation
  - Token revocation support
- **RBAC:** Three roles (TEACHER, PARENT, ADMIN)
- **Class-Scoped Access:** Membership validation on all operations

### Infrastructure Security
- **Container:** Non-root user, read-only filesystem, dropped capabilities
- **Network:** Kubernetes NetworkPolicies restrict ingress/egress
- **Secrets:** Environment-based, never hardcoded
- **Rate Limiting:** Per-IP throttling to prevent abuse
- **CORS:** Configurable allowed origins

### Data Protection
- **PII Protection:** No sensitive data in logs
- **SQL Injection:** Parameterized queries throughout
- **File Uploads:** Content-type whitelist, size limits (5MB)
- **Error Handling:** Structured responses, no stack traces

## 📡 API Endpoints

### Authentication (Public)
```
POST   /v1/auth/register   - User registration
POST   /v1/auth/login      - User login
POST   /v1/auth/refresh    - Refresh access token
POST   /v1/auth/logout     - Logout & revoke token
```

### Classes (Protected)
```
POST   /v1/classes         - Create class (Teacher/Admin)
GET    /v1/classes         - List my classes
GET    /v1/classes/:id     - Get class details
GET    /v1/classes/:id/members - List members (Teacher)
```

### Photos (Protected)
```
POST   /v1/classes/:id/photos - Get presigned upload URL (Teacher)
GET    /v1/classes/:id/photos - List photos with view URLs
```

### Health Checks
```
GET    /healthz            - Liveness probe
GET    /readyz             - Readiness probe (checks DB + S3)
```

## 🗄️ Database Schema

- **users** - Authentication & roles
- **profiles** - User display information
- **classes** - Class definitions
- **class_members** - User-class associations
- **photos** - Photo metadata (S3 keys only)
- **absences** - Student absence tracking
- **messages** - Direct messaging
- **announcements** - Class/global announcements
- **refresh_tokens** - Token management

All tables include proper indexes, foreign keys, and timestamps.

## 🚀 Deployment

### Local Development
```bash
# Start infrastructure
docker-compose up -d

# Run migrations
make migrate-up

# Start server
make run
```

### Docker
```bash
# Build image
docker build -t tinyschoolhub/api:latest .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="..." \
  -e JWT_SECRET="..." \
  -e STORAGE_ENDPOINT="..." \
  tinyschoolhub/api:latest
```

### Kubernetes with Helm
```bash
# Install
helm install tiny-school-hub deploy/helm/tiny-school-hub \
  --set config.env=production \
  --set database.url="postgres://..." \
  --set storage.endpoint="s3.amazonaws.com" \
  --create-namespace \
  --namespace tiny-school-hub

# Upgrade
helm upgrade tiny-school-hub deploy/helm/tiny-school-hub

# Uninstall
helm uninstall tiny-school-hub
```

### Kubernetes with Kustomize
```bash
# Dev environment
kubectl apply -k deploy/kustomize/dev

# Staging environment
kubectl apply -k deploy/kustomize/staging

# Production environment
kubectl apply -k deploy/kustomize/prod
```

## 🔧 Configuration

All configuration via environment variables. See `.env.example` for complete list.

### Required Variables
```bash
# Server
PORT=8080
ENV=production

# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# S3-Compatible Storage
STORAGE_ENDPOINT=s3.amazonaws.com
STORAGE_REGION=us-east-1
STORAGE_BUCKET=your-bucket
STORAGE_ACCESS_KEY=your-access-key
STORAGE_SECRET_KEY=your-secret-key
STORAGE_USE_PATH_STYLE=false
STORAGE_INSECURE=false

# Application
RATE_LIMIT=100
CORS_ALLOWED_ORIGINS=https://app.example.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🧪 Testing & Code Quality

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Format code
make fmt

# Run all pre-commit checks
make pre-commit
```

### Pre-commit Hooks

This project uses Git pre-commit hooks to ensure code quality:

```bash
# Install hooks (runs automatically before each commit)
make install-hooks

# Run checks manually
make pre-commit

# Bypass hooks (emergency only)
git commit --no-verify -m "message"
```

**What gets checked:**
- ✅ Code formatting (gofmt)
- ✅ Go vet (suspicious constructs)
- ✅ Build compilation
- ✅ Unit tests
- ✅ Module dependencies (go mod tidy)
- ✅ Security scanning (gosec, if installed)
- ✅ Sensitive data detection
- ✅ SQL migration validation

See [Pre-commit Documentation](.github/PRE_COMMIT_CHECKS.md) for details.

## 🌊 Development Workflow

This project uses **GitFlow** for branch management and **Renovate** for dependency updates.

### 🎯 **NEW: Complete Workflow Guide**

👉 **[📖 Read the Complete Workflow Guide](docs/WORKFLOW_GUIDE.md)** 👈

**Everything you need to know:**
- ✅ Daily development workflow
- ✅ Creating releases step-by-step
- ✅ Hotfix procedures
- ✅ What happens automatically (CI/CD)
- ✅ Command cheat sheet
- ✅ Troubleshooting guide

### Quick Start

```bash
# Clone and setup
git clone https://github.com/TinySchoolHub/tiny-school-hub-api-backend.git
cd tiny-school-hub-api-backend

# Run GitFlow setup (first time only)
./scripts/setup-gitflow.sh

# Start working on a feature
git checkout develop
git pull origin develop
git checkout -b feature/my-awesome-feature

# Make changes, commit, and push
git add .
git commit -m "feat: add awesome feature"
git push -u origin feature/my-awesome-feature

# Create PR: feature/my-awesome-feature → develop
```

### Branches

- **`main`** - Production-ready code (protected)
- **`develop`** - Integration branch for next release (protected)
- **`feature/*`** - New features (branch from `develop`)
- **`bugfix/*`** - Bug fixes (branch from `develop`)
- **`hotfix/*`** - Urgent production fixes (branch from `main`)
- **`release/*`** - Release preparation (branch from `develop`)

### Automated Dependency Updates

**Renovate** automatically creates PRs for:
- ✅ Go module updates (grouped together)
- ✅ GitHub Actions updates
- ✅ Docker base image updates
- 🤖 Patch updates auto-merge after CI passes

Review and merge dependency PRs regularly to stay up to date.

### Documentation

- **[📖 Workflow Guide](docs/WORKFLOW_GUIDE.md)** - **⭐ START HERE** - Complete step-by-step guide
- **[GitFlow Guide](docs/GITFLOW.md)** - Detailed workflow documentation
- **[Quick Reference](docs/GITFLOW_QUICK_REFERENCE.md)** - Command cheat sheet
- **[Release Guide](docs/RELEASE_GUIDE.md)** - How to create releases
- **[CI/CD Analysis](docs/CICD_ANALYSIS.md)** - Pipeline optimization details
- **[Branch Protection](docs/BRANCH_PROTECTION_RULES.md)** - GitHub settings guide
- **[Contributing](CONTRIBUTING.md)** - Contribution guidelines

## 📚 Documentation

- **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Complete feature overview
- **[AWS_SDK_MIGRATION.md](AWS_SDK_MIGRATION.md)** - AWS SDK v2 migration details
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and changes

## 📄 License

Private. All rights reserved.

## 🙋 Support

For issues or questions, please contact the development team.

### Troubleshooting

**Server won't start:**
- Check `.env` file exists with all required variables
- Verify database is accessible: `psql $DATABASE_URL`
- Check storage endpoint is reachable

**Database connection errors:**
- Ensure PostgreSQL is running: `docker-compose ps`
- Run migrations: `make migrate-up`
- Check DATABASE_URL format

**Storage errors:**
- Verify S3 credentials are correct
- Check endpoint URL (include http:// for insecure)
- Test with MinIO locally first

**Build errors:**
- Update Go to 1.24+: `go version`
- Clean dependencies: `go clean -modcache && go mod tidy`
- Rebuild: `make build`

## 📊 Project Stats

- **Lines of Code:** ~3,500+
- **Binary Size:** 16MB (optimized)
- **Dependencies:** Minimal, all actively maintained
- **Go Version:** 1.24+
- **Database Tables:** 9
- **API Endpoints:** 13+
- **Security:** Argon2id + JWT + RBAC

---

**Built with ❤️ for TinySchoolHub**  
**Last Updated:** November 21, 2025  
**Status:** Production Ready ✅
