# apidiff

A CLI tool that compares two OpenAPI/Swagger specifications and reports breaking changes, added/removed endpoints, and schema differences.

## Installation

Download the latest binary from [Releases](https://github.com/martyn/apidiff/releases) and add it to your PATH.

### Windows

```powershell
# Download and place in a directory on your PATH
Invoke-WebRequest -Uri "https://github.com/martyn/apidiff/releases/latest/download/apidiff-windows-amd64.exe" -OutFile "$env:LOCALAPPDATA\Microsoft\WindowsApps\apidiff.exe"
```

### macOS / Linux

```bash
# Download (adjust OS and arch as needed)
curl -Lo apidiff https://github.com/martyn/apidiff/releases/latest/download/apidiff-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
chmod +x apidiff
sudo mv apidiff /usr/local/bin/
```

### From source

```bash
go install github.com/martyn/apidiff@latest
```

## Usage

```
apidiff compare --old <old-spec> --new <new-spec> [flags]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--old` | Path to the old/original OpenAPI spec | (required) |
| `--new` | Path to the new/updated OpenAPI spec | (required) |
| `--format` | Output format: `console`, `json`, `markdown`, `sarif` | `console` |
| `--recommend-version` | Include a suggested semantic version bump | `false` |
| `--fail-on-breaking` | Exit with non-zero code if breaking changes are found | `false` |

### Examples

```bash
# Basic comparison with styled console output
apidiff compare --old swagger-v1.yaml --new swagger-v2.yaml

# JSON output for programmatic consumption
apidiff compare --old swagger-v1.yaml --new swagger-v2.yaml --format json

# Markdown report
apidiff compare --old swagger-v1.yaml --new swagger-v2.yaml --format markdown

# SARIF output for CI integration
apidiff compare --old swagger-v1.yaml --new swagger-v2.yaml --format sarif

# Include version bump recommendation
apidiff compare --old swagger-v1.yaml --new swagger-v2.yaml --recommend-version

# Fail in CI if breaking changes exist
apidiff compare --old swagger-v1.yaml --new swagger-v2.yaml --fail-on-breaking
```

### Example Output

```
Comparing OpenAPI specifications...

✓ 3 endpoints added
✓ 1 endpoint removed
✓ 3 schema changes

Breaking Changes

✗ DELETE /users/{id} removed

✗ Property lastName removed from User
  User.lastName

✓ Recommended version bump: MAJOR
```

## Breaking Change Detection

The following changes are classified as breaking:

| Change | Severity | Version Bump |
|--------|----------|--------------|
| Endpoint removed | Breaking | Major |
| Required field added | Breaking | Major |
| Property removed | Breaking | Major |
| Property type changed | Breaking | Major |
| New optional endpoint | Info | Minor |
| New optional property | Info | Minor |
| Documentation only | Info | Patch |

## Output Formats

- **console** — Styled terminal output with color-coded markers (default)
- **json** — Structured JSON with summary and change list
- **markdown** — Markdown tables suitable for PR comments
- **sarif** — [SARIF v2.1.0](https://sarifweb.azurewebsites.net/) for GitHub Code Scanning and CI tools

## Supported Specifications

- OpenAPI 3.x (YAML and JSON)
- Swagger 2.x (automatically converted to OpenAPI 3.x for comparison)

## License

[MIT](LICENSE)
