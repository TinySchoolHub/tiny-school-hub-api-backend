# Pre-commit YAML Validation - Helm Template Exclusion

## ✅ Configuration Updated

Helm template files are now automatically excluded from YAML validation in both pre-commit configurations.

## 🎯 Why Exclude Helm Templates?

Helm templates contain **Go template syntax** (e.g., `{{ .Values.replicas }}`), which makes them invalid as pure YAML:

```yaml
# This looks like YAML but isn't valid YAML due to Go templates
replicas: {{ .Values.replicas }}
image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
```

These files are only valid after Helm processes them and replaces the template placeholders with actual values.

## 📝 What Was Changed

### 1. Python Pre-commit Framework (`.pre-commit-config.yaml`)

```yaml
- id: check-yaml
  exclude: ^deploy/helm/     # ✅ Already configured
```

This excludes **all files** under `deploy/helm/` from YAML validation.

### 2. Bash Git Hook (`scripts/pre-commit.sh`)

```bash
# Skip Helm template files (they contain Go templates, not pure YAML)
if [[ "$yaml_file" == deploy/helm/*/templates/* ]]; then
    echo -e "${YELLOW}  Skipping Helm template: $yaml_file${NC}"
    continue
fi
```

This specifically excludes files in `deploy/helm/*/templates/` directories.

## 🔍 What Gets Validated

### ✅ These YAML files ARE validated:
- `docker-compose.yml`
- `.github/workflows/*.yml`
- `.pre-commit-config.yaml`
- `deploy/kustomize/**/*.yaml`
- `deploy/helm/tiny-school-hub/Chart.yaml` ⚠️
- `deploy/helm/tiny-school-hub/values.yaml` ⚠️

### ⏭️ These files are SKIPPED:
- `deploy/helm/tiny-school-hub/templates/*.yaml` ✅ (Contains Go templates)
- `deploy/helm/tiny-school-hub/templates/_helpers.tpl` ✅ (Template partial)

## 🧪 Testing

To verify the exclusion works:

```bash
# Stage a Helm template file
git add deploy/helm/tiny-school-hub/templates/deployment.yaml

# Run pre-commit
make pre-commit

# You should see:
# 9. Validating YAML files...
#   Skipping Helm template: deploy/helm/tiny-school-hub/templates/deployment.yaml
# ✓ YAML validation
```

## 🔧 Manual YAML Validation

If you need to validate processed Helm templates:

```bash
# Render templates with Helm
helm template tiny-school-hub deploy/helm/tiny-school-hub > rendered.yaml

# Validate the rendered output
yamllint rendered.yaml

# Or use Helm's built-in validation
helm lint deploy/helm/tiny-school-hub
```

## 📋 Pattern Matching Details

### Python Pre-commit (Regex)
```regex
^deploy/helm/
```
- `^` = Start of path
- Matches: `deploy/helm/anything`

### Bash Script (Glob)
```bash
deploy/helm/*/templates/*
```
- `*` = Any directory name
- Matches: `deploy/helm/tiny-school-hub/templates/deployment.yaml`
- Matches: `deploy/helm/any-chart/templates/service.yaml`

## 🎓 Best Practices

1. **Keep `Chart.yaml` and `values.yaml` valid YAML** - They should NOT contain Go templates
2. **Use `templates/_helpers.tpl`** for reusable template functions
3. **Test Helm charts** with `helm lint` and `helm template`
4. **Validate rendered templates** in CI/CD pipeline

## 🚀 CI/CD Integration

Add Helm validation to your CI pipeline:

```yaml
# .github/workflows/ci.yml
- name: Validate Helm Charts
  run: |
    helm lint deploy/helm/tiny-school-hub
    helm template test deploy/helm/tiny-school-hub | kubeval --strict
```

## 📚 Related Commands

```bash
# Lint Helm chart
helm lint deploy/helm/tiny-school-hub

# Render templates (dry-run)
helm template tiny-school-hub deploy/helm/tiny-school-hub

# Validate against Kubernetes API
helm template tiny-school-hub deploy/helm/tiny-school-hub | kubectl apply --dry-run=client -f -

# Debug template rendering
helm template tiny-school-hub deploy/helm/tiny-school-hub --debug
```

## ✅ Summary

- ✅ Helm templates are excluded from YAML validation
- ✅ Regular YAML files are still validated
- ✅ Both pre-commit methods handle the exclusion
- ✅ Charts can be validated with `helm lint`
- ✅ No false positives from Go template syntax

---

**Issue Resolved:** Helm template YAML files are now properly excluded from pre-commit validation checks.
