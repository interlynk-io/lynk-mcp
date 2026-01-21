// Copyright 2025 Interlynk.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"github.com/interlynk-io/lynk-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// Server is the MCP server for Lynk API
type Server struct {
	client *api.Client
	logger *zap.Logger
	mcp    *server.MCPServer
}

// NewServer creates a new MCP server
func NewServer(client *api.Client, logger *zap.Logger) *Server {
	s := &Server{
		client: client,
		logger: logger,
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"lynk-version",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, false),
	)

	s.mcp = mcpServer
	s.registerTools()
	s.registerResources()

	return s
}

// Serve starts the MCP server in stdio mode
func (s *Server) Serve() error {
	return server.ServeStdio(s.mcp)
}

// registerTools registers all MCP tools
func (s *Server) registerTools() {
	// Organization tools
	s.mcp.AddTool(mcp.NewTool("get_organization",
		mcp.WithDescription("Get current organization information including metrics"),
	), s.handleGetOrganization)

	// Product tools
	s.mcp.AddTool(mcp.NewTool("list_products",
		mcp.WithDescription("List all products in the organization"),
		mcp.WithString("search", mcp.Description("Search term to filter by name")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 20)")),
	), s.handleListProducts)

	s.mcp.AddTool(mcp.NewTool("get_product",
		mcp.WithDescription("Get details of a specific product including its environments"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the product")),
	), s.handleGetProduct)

	// Environment tools
	s.mcp.AddTool(mcp.NewTool("list_environments",
		mcp.WithDescription("List environments within a product"),
		mcp.WithString("product_id", mcp.Required(), mcp.Description("The UUID of the product")),
		mcp.WithString("search", mcp.Description("Search term to filter by name")),
	), s.handleListEnvironments)

	s.mcp.AddTool(mcp.NewTool("get_environment",
		mcp.WithDescription("Get details of a specific environment"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the environment")),
	), s.handleGetEnvironment)

	// Version tools
	s.mcp.AddTool(mcp.NewTool("list_versions",
		mcp.WithDescription("List versions in an environment"),
		mcp.WithString("environment_id", mcp.Required(), mcp.Description("The UUID of the environment")),
		mcp.WithString("lifecycle", mcp.Description("Filter by lifecycle stage (e.g., released, development)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 20)")),
	), s.handleListVersions)

	s.mcp.AddTool(mcp.NewTool("get_version",
		mcp.WithDescription("Get details of a specific version including statistics"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the version")),
	), s.handleGetVersion)

	s.mcp.AddTool(mcp.NewTool("compare_versions",
		mcp.WithDescription("Compare two versions and show the differences (drift analysis)"),
		mcp.WithString("source_version_id", mcp.Required(), mcp.Description("The UUID of the source version")),
		mcp.WithString("target_version_id", mcp.Required(), mcp.Description("The UUID of the target version to compare against")),
	), s.handleCompareVersions)

	// Component tools
	s.mcp.AddTool(mcp.NewTool("list_components",
		mcp.WithDescription("List components in a version"),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The UUID of the version")),
		mcp.WithString("search", mcp.Description("Search term to filter components")),
		mcp.WithString("kind", mcp.Description("Filter by component kind (e.g., library, application)")),
		mcp.WithBoolean("direct", mcp.Description("Filter to direct dependencies only")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleListComponents)

	s.mcp.AddTool(mcp.NewTool("get_component",
		mcp.WithDescription("Get details of a specific component"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the component")),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The UUID of the version containing the component")),
	), s.handleGetComponent)

	// Vulnerability tools
	s.mcp.AddTool(mcp.NewTool("list_vulnerabilities",
		mcp.WithDescription("List vulnerabilities in a version with optional filters"),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The UUID of the version")),
		mcp.WithString("severity", mcp.Description("Filter by severity (critical, high, medium, low)")),
		mcp.WithString("vex_status", mcp.Description("Filter by VEX status (e.g., affected, not_affected, fixed)")),
		mcp.WithBoolean("kev", mcp.Description("Filter to only KEV (Known Exploited Vulnerabilities)")),
		mcp.WithString("search", mcp.Description("Search term to filter vulnerabilities")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleListVulnerabilities)

	s.mcp.AddTool(mcp.NewTool("get_vulnerability",
		mcp.WithDescription("Get details of a specific vulnerability"),
		mcp.WithString("vuln_id", mcp.Required(), mcp.Description("The CVE ID (e.g., CVE-2021-44228) or UUID")),
	), s.handleGetVulnerability)

	s.mcp.AddTool(mcp.NewTool("search_vulnerabilities",
		mcp.WithDescription("Search vulnerabilities across all products"),
		mcp.WithString("search", mcp.Description("Search term (CVE ID, component name, etc.)")),
		mcp.WithString("severity", mcp.Description("Filter by severity")),
		mcp.WithBoolean("kev", mcp.Description("Filter to only KEV")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleSearchVulnerabilities)

	// Policy tools
	s.mcp.AddTool(mcp.NewTool("list_policies",
		mcp.WithDescription("List security policies in the organization"),
		mcp.WithString("search", mcp.Description("Search term to filter policies")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 20)")),
	), s.handleListPolicies)

	s.mcp.AddTool(mcp.NewTool("get_policy",
		mcp.WithDescription("Get details of a specific policy including its rules"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the policy")),
	), s.handleGetPolicy)

	s.mcp.AddTool(mcp.NewTool("list_policy_violations",
		mcp.WithDescription("List policy evaluation results/violations"),
		mcp.WithString("policy_id", mcp.Description("Filter by policy UUID")),
		mcp.WithString("version_id", mcp.Description("Filter by version UUID")),
		mcp.WithString("result_type", mcp.Description("Filter by result type (pass, fail, warn)")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleListPolicyViolations)

	// License tools
	s.mcp.AddTool(mcp.NewTool("list_licenses",
		mcp.WithDescription("List licenses used in the organization's versions"),
		mcp.WithString("status", mcp.Description("Filter by license status (approved, rejected, unspecified)")),
		mcp.WithString("search", mcp.Description("Search term to filter licenses")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleListLicenses)
}

// registerResources registers MCP resources
func (s *Server) registerResources() {
	// Register resource templates
	s.mcp.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"version:///{version_id}",
			"Complete version information",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleVersionResource,
	)

	s.mcp.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"version:///{version_id}/components",
			"All components in a version",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleVersionComponentsResource,
	)

	s.mcp.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"version:///{version_id}/vulnerabilities",
			"All vulnerabilities in a version",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleVersionVulnerabilitiesResource,
	)

	s.mcp.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"environment:///{environment_id}/latest-version",
			"Most recent version for an environment",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleEnvironmentLatestVersionResource,
	)

	s.mcp.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"organization:///summary",
			"Organization overview and summary",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleOrganizationSummaryResource,
	)

	s.mcp.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"vulnerability:///{cve_id}",
			"Vulnerability details by CVE ID",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleVulnerabilityResource,
	)
}
