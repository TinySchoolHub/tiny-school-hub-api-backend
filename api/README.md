# API Documentation

This directory contains the OpenAPI specification and documentation viewers for the Tiny School Hub API.

## 📁 Files

- **`openapi.yaml`** - OpenAPI 3.1 specification (the source of truth)
- **`swagger.html`** - Standalone Swagger UI viewer (open directly in browser)

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

Start Swagger UI and ReDoc alongside your services:

```bash
# Start all services including Swagger UI and ReDoc
make docker-up

# Open Swagger UI (interactive)
make swagger
# Or manually: http://localhost:8081

# Open ReDoc (clean, read-only)
make redoc
# Or manually: http://localhost:8082

# Open both
make docs
```

### Option 2: Standalone HTML (No Docker)

Simply open the HTML file in your browser:

```bash
# macOS
open api/swagger.html

# Linux
xdg-open api/swagger.html

# Windows
start api/swagger.html
```

This loads Swagger UI from CDN and reads `openapi.yaml` directly.

### Option 3: VS Code Extension

Install the [OpenAPI (Swagger) Editor](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi) extension:

1. Open `openapi.yaml` in VS Code
2. Press `Cmd+Shift+P` (or `Ctrl+Shift+P`)
3. Run "OpenAPI: Show Preview"

### Option 4: Online Viewer

Upload `openapi.yaml` to:
- **Swagger Editor**: https://editor.swagger.io/
- **Redocly**: https://redocly.github.io/redoc/

## 📚 Documentation Viewers

### Swagger UI (http://localhost:8081)

**Best for: Interactive testing**

- ✅ "Try it out" functionality
- ✅ Test API calls directly
- ✅ See request/response in real-time
- ✅ OAuth2 authentication support
- ✅ Download client SDKs
- 🎯 Great for developers building integrations

### ReDoc (http://localhost:8082)

**Best for: Reading documentation**

- ✅ Clean, three-panel layout
- ✅ Better for large APIs
- ✅ Search functionality
- ✅ Code samples in multiple languages
- ✅ Print-friendly
- 🎯 Great for API consumers and documentation

## 🔐 Testing Authentication

### 1. Login to Get Token

In Swagger UI:
1. Scroll to **POST /v1/auth/login**
2. Click "Try it out"
3. Enter credentials:
   ```json
   {
     "email": "teacher@example.com",
     "password": "password123"
   }
   ```
4. Click "Execute"
5. Copy the `access_token` from response

### 2. Authorize

1. Click the **🔒 Authorize** button at the top
2. Enter: `Bearer <your-access-token>`
3. Click "Authorize"
4. Now you can test protected endpoints!

## 🛠️ Development

### Validate the Spec

```bash
# Using npx (no install needed)
npx @apidevtools/swagger-cli validate api/openapi.yaml

# Or install globally
npm install -g @apidevtools/swagger-cli
swagger-cli validate api/openapi.yaml
```

### Generate Client SDKs

```bash
# TypeScript/Axios client for frontend
npx @openapitools/openapi-generator-cli generate \
  -i api/openapi.yaml \
  -g typescript-axios \
  -o clients/typescript

# Python client
npx @openapitools/openapi-generator-cli generate \
  -i api/openapi.yaml \
  -g python \
  -o clients/python

# Go client
npx @openapitools/openapi-generator-cli generate \
  -i api/openapi.yaml \
  -g go \
  -o clients/go
```

### Update Documentation

When you change the API:

1. Update `openapi.yaml`
2. Validate: `swagger-cli validate api/openapi.yaml`
3. Test in Swagger UI: `make swagger`
4. Commit changes

The documentation updates automatically - no rebuild needed!

## 🎨 Customization

### Swagger UI Themes

Edit `api/swagger.html` to customize:
- Colors
- Logo
- Layout
- Default expansions
- Authentication settings

### Docker Configuration

Edit `docker-compose.yml` to change:
- Ports (8081, 8082)
- Environment variables
- API URL

## 📖 Resources

- [OpenAPI Specification](https://spec.openapis.org/oas/v3.1.0)
- [Swagger UI Documentation](https://swagger.io/docs/open-source-tools/swagger-ui/)
- [ReDoc Documentation](https://redocly.com/docs/redoc/)
- [OpenAPI Generator](https://openapi-generator.tech/)

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
- name: Validate OpenAPI Spec
  run: |
    npm install -g @apidevtools/swagger-cli
    swagger-cli validate api/openapi.yaml

- name: Generate Documentation
  run: |
    docker run --rm -v ${PWD}:/local \
      swaggerapi/swagger-ui \
      /local/api/openapi.yaml
```

## 🐛 Troubleshooting

### Swagger UI shows "Failed to load API definition"

**Solution:** Make sure the file path is correct:
- Docker: `SWAGGER_JSON: /api/openapi.yaml`
- Standalone HTML: `url: "./openapi.yaml"`

### CORS errors in "Try it out"

**Solution:** 
1. Make sure API is running: `make run`
2. Check CORS_ALLOWED_ORIGINS in `.env`
3. Add `http://localhost:8081` to allowed origins

### Changes not reflecting

**Solution:**
1. Hard refresh browser: `Cmd+Shift+R` (macOS) or `Ctrl+Shift+F5` (Windows)
2. Restart Docker containers: `make docker-down && make docker-up`

### Ports already in use

**Solution:** Change ports in `docker-compose.yml`:
```yaml
swagger-ui:
  ports:
    - "8091:8080"  # Change 8081 to 8091
```

## 💡 Tips

- **Use Swagger UI** for development and testing
- **Use ReDoc** for client-facing documentation
- **Keep openapi.yaml** in sync with code
- **Validate** before committing changes
- **Generate clients** for type-safe API consumption
- **Version** your API (already at v1!)

---

**Happy documenting! 📚**

For questions, see the main [README.md](../README.md)
