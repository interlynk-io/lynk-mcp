# lynk-mcp: AI-Powered SBOM & Vulnerability Management

MCP server for Interlynk API. This server enables AI assistants like Claude, Cursor, and VS Code Copilot to interact with your Lynk organization for SBOM management, vulnerability tracking, and compliance checking.

## Quick Start

```bash
# Install via Homebrew
brew install interlynk-io/tap/lynk-mcp

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

### Using Homebrew (macOS/Linux)

```bash
brew install interlynk-io/tap/lynk-mcp
```

### Using Go Install

```bash
go install github.com/interlynk-io/lynk-mcp/cmd/lynk-mcp@latest
```

### Using Docker

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
| `get_organization` | Get organization info and metrics |
| `list_products` | List all products |
| `get_product` | Get product details with environments |
| `list_environments` | List environments in a product |
| `get_environment` | Get environment details |

### Versions & Components

| Tool | Description |
|------|-------------|
| `list_versions` | List versions in an environment |
| `get_version` | Get version details with statistics |
| `list_components` | List components in a version |
| `get_component` | Get component details |
| `compare_versions` | Compare two versions for drift |

### Vulnerabilities

| Tool | Description |
|------|-------------|
| `list_vulnerabilities` | List vulnerabilities with filters |
| `get_vulnerability` | Get vulnerability by CVE or UUID |
| `search_vulnerabilities` | Search across all products |

### Policies & Compliance

| Tool | Description |
|------|-------------|
| `list_policies` | List security policies |
| `get_policy` | Get policy details with rules |
| `list_policy_violations` | List policy evaluation results |
| `list_licenses` | List licenses with filtering |

## Available Resources

| Resource URI | Description |
|--------------|-------------|
| `version:///{version_id}` | Complete version information |
| `version:///{version_id}/components` | All components in a version |
| `version:///{version_id}/vulnerabilities` | All vulnerabilities in a version |
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
