# lynk-mcp

MCP (Model Context Protocol) server for Lynk version management API. This server enables AI assistants like Claude to interact with your Lynk organization for version management, vulnerability tracking, and compliance checking.

## Features

- **Organization Management**: View organization info and metrics
- **Products & Environments**: List and explore products and environments
- **Version Operations**: List, view, and compare versions with drift analysis
- **Component Analysis**: Search and explore components across versions
- **Vulnerability Management**: Query vulnerabilities with filtering by severity, KEV status, and VEX status
- **Policy Compliance**: View policies and their evaluation results
- **License Management**: Track and filter licenses across your organization

## Installation

### Using Go Install

```bash
go install github.com/interlynk-io/lynk-mcp/cmd/lynk-mcp@latest
```

### Using Homebrew (macOS/Linux)

```bash
brew install interlynk-io/tap/lynk-mcp
```

### From Source

```bash
git clone https://github.com/interlynk-io/lynk-mcp.git
cd lynk-mcp
make build
```

## Configuration

### Initial Setup

Run the configuration command to set up your API token:

```bash
lynk-mcp configure
```

This will prompt you for:
1. API Endpoint (defaults to https://api.interlynk.io/lynkapi)
2. API Token (your Lynk API key starting with `lynk_live_`, `lynk_staging_`, or `lynk_test_`)

The token is stored securely in your system keychain.

### Verify Connection

Test your configuration:

```bash
lynk-mcp verify
```

This will display your organization information if the connection is successful.

## Usage with AI Assistants

### Claude Desktop

Add the following to your Claude Desktop configuration file:

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

Restart Claude Desktop after making this change.

### Claude Code (CLI)

Add the MCP server to Claude Code using the CLI:

```bash
claude mcp add lynk -- lynk-mcp serve
```

Or manually add to your Claude Code settings file:

**macOS**: `~/.claude/settings.json`
**Linux**: `~/.claude/settings.json`
**Windows**: `%USERPROFILE%\.claude\settings.json`

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

For project-specific configuration, create a `.mcp.json` file in your project root:

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

### VS Code

VS Code supports MCP servers through the built-in agent mode or extensions like Continue.

#### Using VS Code Agent Mode (v1.99+)

Add to your VS Code settings (`settings.json`):

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

Or create a `.vscode/mcp.json` file in your workspace:

```json
{
  "servers": {
    "lynk": {
      "command": "lynk-mcp",
      "args": ["serve"]
    }
  }
}
```

#### Using Continue Extension

If using the [Continue](https://continue.dev) extension, add to your `~/.continue/config.json`:

```json
{
  "experimental": {
    "modelContextProtocolServers": [
      {
        "transport": {
          "type": "stdio",
          "command": "lynk-mcp",
          "args": ["serve"]
        }
      }
    ]
  }
}
```

### Cursor

Add the MCP server to your Cursor configuration:

**macOS**: `~/.cursor/mcp.json`
**Windows**: `%USERPROFILE%\.cursor\mcp.json`
**Linux**: `~/.cursor/mcp.json`

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

Alternatively, create a `.cursor/mcp.json` file in your project root for project-specific configuration.

### Zed

Add the MCP server to your Zed settings:

**macOS**: `~/.config/zed/settings.json`
**Linux**: `~/.config/zed/settings.json`

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

Or for project-specific configuration, add to `.zed/settings.json` in your project root.

## Available Tools

### Organization & Products

| Tool | Description |
|------|-------------|
| `get_organization` | Get current organization info and metrics |
| `list_products` | List all products |
| `get_product` | Get product details with its environments |
| `list_environments` | List environments within a product |
| `get_environment` | Get environment details |

### Versions & Components

| Tool | Description |
|------|-------------|
| `list_versions` | List versions in an environment |
| `get_version` | Get version details with statistics |
| `list_components` | List components in a version |
| `get_component` | Get component details |
| `compare_versions` | Compare two versions for drift analysis |

### Vulnerabilities

| Tool | Description |
|------|-------------|
| `list_vulnerabilities` | List vulnerabilities in a version with filters |
| `get_vulnerability` | Get vulnerability details by CVE or UUID |
| `search_vulnerabilities` | Search vulnerabilities across all products |

### Policies & Compliance

| Tool | Description |
|------|-------------|
| `list_policies` | List security policies |
| `get_policy` | Get policy details with rules |
| `list_policy_violations` | List policy evaluation results |
| `list_licenses` | List licenses with state filtering |

## Available Resources

| Resource URI | Description |
|--------------|-------------|
| `version:///{version_id}` | Complete version information |
| `version:///{version_id}/components` | All components in a version |
| `version:///{version_id}/vulnerabilities` | All vulnerabilities in a version |
| `environment:///{environment_id}/latest-version` | Most recent version for an environment |
| `organization:///summary` | Organization overview |
| `vulnerability:///{cve_id}` | Vulnerability details by CVE ID |

## Example Queries

Once configured, you can ask your AI assistant questions like these:

### Getting Started

- "List all my products"
- "Show me the environments in [product name]"
- "What's the latest SBOM version for [environment name]?"

### Vulnerability Analysis

- "Show me all critical vulnerabilities in my organization"
- "List vulnerabilities with KEV (Known Exploited Vulnerabilities) status"
- "What vulnerabilities in [product name] have a fix available?"
- "Show me all high and critical severity vulnerabilities across all products"
- "Which components have the most vulnerabilities?"

### Searching for Specific Attacks & CVEs

- "Are any of my products affected by the Shai Hulud attack (CVE-2024-3094)?"
- "Check if my organization is vulnerable to Log4Shell (CVE-2021-44228)"
- "Search for any components affected by CVE-2023-44487 (HTTP/2 Rapid Reset)"
- "Do I have any xz-utils components that might be affected by the backdoor?"
- "Find all occurrences of OpenSSL vulnerabilities in my SBOMs"

### Generating Cybersecurity Reports

- "Generate a security summary report for [product name] including all critical vulnerabilities, their CVSS scores, and remediation status"
- "Create an executive summary of our organization's vulnerability posture"
- "List all components with known vulnerabilities grouped by severity for compliance reporting"
- "Generate a report of all policy violations across my products"
- "Summarize the vulnerability trends between the last two versions of [product name]"
- "Create a license compliance report showing all copyleft licenses in use"

### Drift Analysis & Comparisons

- "Compare the last two versions of [product name] and highlight security changes"
- "What new vulnerabilities were introduced in the latest version?"
- "Show me components that were added or removed between versions"
- "Has our security posture improved or degraded since the last release?"

### Policy & Compliance

- "What policies are currently failing for [environment name]?"
- "Show me all SBOM versions that violate our security policies"
- "List all components using GPL licenses"
- "Which products have components with deprecated licenses?"

### Component Analysis

- "Find all instances of log4j across my organization"
- "List all components from [vendor name]"
- "Show me all direct dependencies vs transitive dependencies in [version]"
- "Which components are missing PURL identifiers?"

## Configuration File

Configuration is stored in `~/.lynk-mcp/config.yaml`:

```yaml
api:
  endpoint: "https://api.interlynk.io/lynkapi"
  timeout: 30s
logging:
  level: "info"
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `LYNK_MCP_API_ENDPOINT` | Override API endpoint |
| `LYNK_MCP_LOGGING_LEVEL` | Set logging level (debug, info, warn, error) |

## Security

- API tokens are stored in your system's native keychain
  - macOS: Keychain
  - Windows: Credential Manager
  - Linux: Secret Service or encrypted file
- Tokens are never logged or exposed
- All API communication uses HTTPS
- Organization scoping is enforced by the Lynk API

## Development

### Prerequisites

- Go 1.22 or later
- Make (optional)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint
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
├── go.mod
├── Makefile
└── README.md
```

## License

Apache License 2.0

## Support

For issues and feature requests, please visit: https://github.com/interlynk-io/lynk-mcp/issues
