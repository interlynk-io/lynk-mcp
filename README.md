# lynk-mcp

MCP (Model Context Protocol) server for Lynk SBOM management API. This server enables AI assistants like Claude to interact with your Lynk organization for SBOM management, vulnerability tracking, and compliance checking.

## Features

- **Organization Management**: View organization info and metrics
- **Products & Projects**: List and explore project groups (products) and projects (streams)
- **SBOM Operations**: List, view, and compare SBOMs with drift analysis
- **Component Analysis**: Search and explore components across SBOMs
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

## Usage with Claude Desktop

Add the following to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "lynk-sbom": {
      "command": "lynk-mcp",
      "args": ["serve"]
    }
  }
}
```

Restart Claude Desktop after making this change.

## Available Tools

### Organization & Projects

| Tool | Description |
|------|-------------|
| `get_organization` | Get current organization info and metrics |
| `list_project_groups` | List all products/project groups |
| `get_project_group` | Get project group details with its projects |
| `list_projects` | List projects within a project group |
| `get_project` | Get project details |

### SBOMs & Components

| Tool | Description |
|------|-------------|
| `list_sboms` | List SBOMs in a project |
| `get_sbom` | Get SBOM details with statistics |
| `list_components` | List components in an SBOM |
| `get_component` | Get component details |
| `compare_sboms` | Compare two SBOMs for drift analysis |

### Vulnerabilities

| Tool | Description |
|------|-------------|
| `list_vulnerabilities` | List vulnerabilities in an SBOM with filters |
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
| `sbom:///{sbom_id}` | Complete SBOM information |
| `sbom:///{sbom_id}/components` | All components in an SBOM |
| `sbom:///{sbom_id}/vulnerabilities` | All vulnerabilities in an SBOM |
| `project:///{project_id}/latest-sbom` | Most recent SBOM for a project |
| `organization:///summary` | Organization overview |
| `vulnerability:///{cve_id}` | Vulnerability details by CVE ID |

## Example Queries

Once configured with Claude Desktop, you can ask questions like:

- "List all my products"
- "Show me the vulnerabilities in the latest SBOM for [project name]"
- "What critical vulnerabilities are in my organization?"
- "Compare the last two versions of [product name]"
- "What policies are failing for [project name]?"
- "Show me all components using Apache-2.0 license"

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
