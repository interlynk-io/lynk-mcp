# lynk-mcp: AI-Powered SBOM & Vulnerability Management

MCP server for Interlynk API. This server enables AI assistants like Claude, Cursor, and VS Code Copilot to interact with your Lynk organization for SBOM management, vulnerability tracking, and compliance checking.

## Quick Start

```bash
# Install via Homebrew on macOS
brew install --cask interlynk-io/interlynk/lynk-mcp

# Configure your API token
lynk-mcp configure

# Verify connection
lynk-mcp verify
```

Then add to your AI assistant and start asking questions about your SBOMs!

## Why lynk-mcp?

Managing software supply chain security is complex. With lynk-mcp, you can use natural language to:

- **Query vulnerabilities instantly** - "Show me all critical CVEs affecting my products"
- **Track compliance** - "Which products are failing security policies?"
- **Analyze drift** - "What changed between these two versions?"
- **Generate reports** - "Create a security summary for the executive team"
- **Search across SBOMs** - "Find all instances of log4j in my organization"

## Key Features

- **Natural Language Queries**: Ask questions in plain English
- **Multi-Product Analysis**: Search vulnerabilities across your entire organization
- **Version Comparison**: Drift analysis between SBOM versions
- **Compliance Tracking**: Policy violations and license management
- **Works Everywhere**: Claude Desktop, Claude Code, VS Code, Cursor, Zed

## Example Queries

Once configured with your AI assistant, try these:

### Vulnerability Analysis

```
"Show me all critical vulnerabilities in my organization"
"List vulnerabilities with KEV (Known Exploited Vulnerabilities) status"
"What vulnerabilities in [product] have a fix available?"
"Which components have the most vulnerabilities?"
```

### Searching for Specific Attacks & CVEs

```
"Are any of my products affected by the XZ backdoor (CVE-2024-3094)?"
"Check if my organization is vulnerable to Log4Shell (CVE-2021-44228)"
"Search for any components affected by CVE-2023-44487 (HTTP/2 Rapid Reset)"
"Find all occurrences of OpenSSL vulnerabilities in my SBOMs"
```

### Security Reports

```
"Generate a security summary for [product] with all critical vulnerabilities"
"Create an executive summary of our vulnerability posture"
"List all components with known vulnerabilities grouped by severity"
"Summarize vulnerability trends between the last two versions"
```

### Drift Analysis

```
"Compare the last two versions of [product] and highlight security changes"
"What new vulnerabilities were introduced in the latest version?"
"Show me components that were added or removed between versions"
"Has our security posture improved since the last release?"
```

### Policy & Compliance

```
"What policies are currently failing for [environment]?"
"Show me all versions that violate security policies"
"List all components using GPL licenses"
"Which products have deprecated licenses?"
```

### Component Analysis

```
"Find all instances of log4j across my organization"
"List all components from [vendor]"
"Show me direct vs transitive dependencies in [version]"
"Which components are missing PURL identifiers?"
```

## Installation

### Homebrew (macOS)

```bash
brew install --cask interlynk-io/interlynk/lynk-mcp
```

Homebrew installs update with normal Homebrew workflows:

```bash
brew update
brew upgrade --cask lynk-mcp
```

The release workflow opens a PR against `interlynk-io/homebrew-interlynk` when a new tag is published.

### Linux Packages

Download the latest `.deb`, `.rpm`, or `.apk` package from the [GitHub releases](https://github.com/interlynk-io/lynk-mcp/releases/latest), then install the package for your distribution:

```bash
# Debian/Ubuntu
sudo dpkg -i lynk-mcp_*_linux_*.deb

# Fedora/RHEL
sudo rpm -Uvh lynk-mcp_*_linux_*.rpm

# Alpine
sudo apk add --allow-untrusted lynk-mcp_*_linux_*.apk
```

Linux packages are built automatically on every release. A hosted `apt`, `yum`, or `apk` repository is not currently published, so package-manager upgrades require installing the newer release package.

### Windows

Using Scoop:

```powershell
scoop bucket add interlynk https://github.com/interlynk-io/homebrew-interlynk
scoop install interlynk/lynk-mcp
```

Using winget:

```powershell
winget install Interlynk.lynk-mcp
```

The release workflow opens a PR for the Scoop bucket manifest and a winget package manifest PR when a new tag is published.

### Go Install

```bash
go install github.com/interlynk-io/lynk-mcp/cmd/lynk-mcp@latest
```

### Docker

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/interlynk-io/lynk-mcp:latest

# Run with API token
docker run -e LYNK_API_TOKEN=lynk_live_xxx ghcr.io/interlynk-io/lynk-mcp serve
```

### From Source

```bash
git clone https://github.com/interlynk-io/lynk-mcp.git
cd lynk-mcp
make build
```

The binary is placed in `./build/lynk-mcp`. You can run it directly from there, or run `make install` to install it to `$GOPATH/bin` (typically `~/go/bin`) and use it from anywhere.

### Release Automation

Tagged releases publish binaries, archives, checksums, Linux packages, Docker images, and package-manager manifests. The release workflow expects these repository secrets when package-manager publishing is enabled:

| Secret | Purpose |
|--------|---------|
| `INTERLYNK_RELEASE_GITHUB_TOKEN` | Opens Homebrew, Scoop, and winget manifest PRs |
| `INTERLYNK_RELEASE_SSH_KEY` | Pushes signed Homebrew/Scoop PR branches to `interlynk-io/homebrew-interlynk` |
| `INTERLYNK_RELEASE_GPG_PRIVATE_KEY` | Imports the release signing key used for tap commits |
| `INTERLYNK_RELEASE_GPG_PASSPHRASE` | Unlocks the release signing key |

The public key for `INTERLYNK_RELEASE_GPG_PRIVATE_KEY` must be uploaded to the GitHub account that owns the `interlynk-support-bot <support_eng@interlynk.io>` commit identity so GitHub marks tap PR commits as verified. See [Release Distribution](docs/release-distribution.md) for the shared release model used across Interlynk OSS tools.

## Configuration

### Initial Setup

```bash
lynk-mcp configure
```

This prompts for:
1. API Endpoint (defaults to https://api.interlynk.io/lynkapi)
2. API Token (your Lynk API key: `lynk_live_*`, `lynk_staging_*`, `lynk_test_*`, or `lynk_service_test_*`)

The token is stored securely in your system keychain.

### Verify Connection

```bash
lynk-mcp verify
```

### Configuration File

Stored in `~/.lynk-mcp/config.yaml`:

```yaml
api:
  endpoint: "https://api.interlynk.io/lynkapi"
  timeout: 30s
logging:
  level: "info"
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `LYNK_API_TOKEN` | API token (alternative to keychain) |
| `LYNK_MCP_API_ENDPOINT` | Override API endpoint |
| `LYNK_MCP_LOGGING_LEVEL` | Logging level (debug, info, warn, error) |

## AI Assistant Setup

### Claude Desktop

Add to your config file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "lynk": {
      "command": "lynk-mcp",
      "args": ["serve"]
    }
  }
}
```

### Claude Code (CLI)

```bash
claude mcp add lynk -- lynk-mcp serve
```

Or add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "lynk": {
      "command": "lynk-mcp",
      "args": ["serve"]
    }
  }
}
```

### VS Code (v1.99+)

Add to `settings.json` or `.vscode/mcp.json`:

```json
{
  "mcp": {
    "servers": {
      "lynk": {
        "command": "lynk-mcp",
        "args": ["serve"]
      }
    }
  }
}
```

### Cursor

Add to `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "lynk": {
      "command": "lynk-mcp",
      "args": ["serve"]
    }
  }
}
```

### Zed

Add to `~/.config/zed/settings.json`:

```json
{
  "context_servers": {
    "lynk": {
      "command": {
        "path": "lynk-mcp",
        "args": ["serve"]
      }
    }
  }
}
```

### Using Docker with AI Assistants

```json
{
  "mcpServers": {
    "lynk": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-e", "LYNK_API_TOKEN=lynk_live_xxx", "ghcr.io/interlynk-io/lynk-mcp", "serve"]
    }
  }
}
```

## Available Tools

### Organization & Products

| Tool | Description |
|------|-------------|
| `get_organization` | Get current organization information including metrics |
| `list_products` | List all products in the organization |
| `get_product` | Get details of a specific product including its environments |
| `list_environments` | List environments within a product |
| `get_environment` | Get details of a specific environment |

### Versions, SBOM Doctor & Components

| Tool | Description |
|------|-------------|
| `list_versions` | List versions in an environment |
| `get_version` | Get version details with statistics |
| `list_doctor_results` | List SBOM Doctor findings for a version |
| `compare_versions` | Compare two versions and show drift analysis |
| `list_components` | List components in a version |
| `get_component` | Get component details |
| `update_component` | Update component metadata; requires `confirm=true` |
| `update_component_supplier` | Update component supplier metadata; requires `confirm=true` |

### Vulnerabilities & VEX

| Tool | Description |
|------|-------------|
| `list_vulnerabilities` | List vulnerabilities in a version with optional filters |
| `get_vulnerability` | Get vulnerability details by CVE or UUID |
| `list_vex_statuses` | List VEX statuses with UUIDs for CVE triage |
| `list_vex_justifications` | List VEX justifications with UUIDs for CVE triage |
| `update_component_vex` | Update VEX data for a component vulnerability; requires `confirm=true` |
| `search_vulnerabilities` | Search across all products |

### Supply-Chain Security Incidents

| Tool | Description |
|------|-------------|
| `list_security_incidents` | List supply-chain security incidents visible to the current organization |
| `get_security_incident` | Get a supply-chain security incident, including markers and impact state |
| `create_security_incident` | Create a draft security incident; requires operator permissions and `confirm=true` |
| `update_security_incident` | Update editable security incident fields; requires operator permissions and `confirm=true` |
| `add_security_incident_markers` | Add markers to a security incident; requires `confirm=true` |
| `withdraw_security_incident_markers` | Withdraw active markers and resolve related active findings; requires `confirm=true` |
| `publish_security_incident` | Publish a draft incident and queue the initial impact scan; requires `confirm=true` |
| `resolve_security_incident` | Resolve an active security incident; requires `confirm=true` |
| `archive_security_incident` | Archive a security incident; requires `confirm=true` |
| `create_security_incident_update` | Add a timeline update to a security incident; requires operator permissions and `confirm=true` |
| `get_security_incident_findings` | Get customer-facing findings for a security incident in the current organization |
| `suppress_security_incident_finding` | Suppress a security incident finding; requires `confirm=true` and a reason |
| `rerun_security_incident_impact_scan` | Queue impact scanning for an active or resolved incident; requires `confirm=true` |
| `dry_run_security_incident_impact_scan` | Queue a dry-run impact scan for an incident; requires `confirm=true` |
| `get_security_incident_dry_run_result` | Get latest dry-run impact scan results |

### Policies & Compliance

| Tool | Description |
|------|-------------|
| `list_policies` | List security policies |
| `get_policy` | Get policy details with rules |
| `list_policy_violations` | List policy evaluation results |
| `get_ticketing_status` | Get ticketing provider connection and policy application status |
| `list_licenses` | List licenses with filtering |

## Available Resources

| Resource URI | Description |
|--------------|-------------|
| `version:///{version_id}` | Complete version information |
| `version:///{version_id}/components` | All components in a version |
| `version:///{version_id}/vulnerabilities` | All vulnerabilities in a version |
| `version:///{version_id}/doctor-results` | SBOM Doctor findings for a version |
| `environment:///{environment_id}/latest-version` | Most recent version |
| `organization:///summary` | Organization overview |
| `vulnerability:///{cve_id}` | Vulnerability details by CVE |

## Security

- API tokens stored in system keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Tokens never logged or exposed
- All API communication uses HTTPS
- Organization scoping enforced by Lynk API

## Development

### Prerequisites

- Go 1.24 or later

### Building

```bash
make build          # Build for current platform
make install        # Build and install to $GOPATH/bin
make build-all      # Build for all platforms
make test           # Run tests
make lint           # Run linter
```

### Project Structure

```
lynk-mcp/
├── cmd/lynk-mcp/          # CLI entry point
├── internal/
│   ├── api/               # High-level API client
│   ├── config/            # Configuration and keyring
│   ├── graphql/           # GraphQL client and queries
│   └── mcp/               # MCP server implementation
├── Dockerfile             # Multi-platform container build
├── go.mod
├── Makefile
└── README.md
```

## Other Interlynk Tools

- [**sbomqs**](https://github.com/interlynk-io/sbomqs) - SBOM quality scoring and compliance
- [**sbomasm**](https://github.com/interlynk-io/sbomasm) - SBOM assembler, merger, and editor
- [**sbomex**](https://github.com/interlynk-io/sbomex) - Search and download public SBOMs
- [**sbomgr**](https://github.com/interlynk-io/sbomgr) - Context-aware SBOM search

## License

Apache License 2.0

## Support

- [GitHub Issues](https://github.com/interlynk-io/lynk-mcp/issues)
- [Community Slack](https://join.slack.com/t/sbomqa/shared_invite/zt-2jzq1ttgy-4IGzOYBEtHwJdMyYj~BACA)
- [Email](mailto:hello@interlynk.io)

---

Made with care by [Interlynk.io](https://www.interlynk.io)
