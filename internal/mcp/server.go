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
	"context"

	"github.com/interlynk-io/lynk-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type lynkClient interface {
	GetOrganization(ctx context.Context) (*api.Organization, error)
	ListProducts(ctx context.Context, input api.ListProductsInput) (*api.ProductsResult, error)
	GetProduct(ctx context.Context, id string) (*api.Product, error)
	ListVersions(ctx context.Context, input api.ListVersionsInput) (*api.VersionsResult, error)
	GetVersion(ctx context.Context, id string) (*api.Version, error)
	ListDoctorResults(ctx context.Context, input api.ListDoctorResultsInput) (*api.DoctorResultsResult, error)
	CompareVersions(ctx context.Context, sourceVersionID, targetVersionID string) ([]api.VersionDiff, error)
	ListComponents(ctx context.Context, input api.ListComponentsInput) (*api.ComponentsResult, error)
	GetComponent(ctx context.Context, id, versionID string) (*api.VersionComponent, error)
	UpdateComponent(ctx context.Context, input api.UpdateComponentInput) (*api.UpdateComponentResult, error)
	UpdateComponentSupplier(ctx context.Context, input api.UpdateComponentSupplierInput) (*api.UpdateComponentSupplierResult, error)
	ListVersionVulns(ctx context.Context, input api.ListVersionVulnsInput) (*api.ComponentVulnsResult, error)
	GetVuln(ctx context.Context, id, vulnID string) (*api.Vuln, error)
	GetVexStatuses(ctx context.Context) ([]api.VexStatus, error)
	GetVexJustifications(ctx context.Context) ([]api.VexJustification, error)
	UpdateComponentVex(ctx context.Context, input api.UpdateComponentVexInput) (*api.UpdateComponentVexResult, error)
	ListComponentVulns(ctx context.Context, input api.ListComponentVulnsInput) (*api.ComponentVulnsResult, error)
	ListSecurityIncidents(ctx context.Context, input api.ListSecurityIncidentsInput) ([]api.SecurityIncident, error)
	GetSecurityIncident(ctx context.Context, id string) (*api.SecurityIncident, error)
	CreateSecurityIncident(ctx context.Context, input api.CreateSecurityIncidentInput) (*api.CreateSecurityIncidentResult, error)
	UpdateSecurityIncident(ctx context.Context, input api.UpdateSecurityIncidentInput) (*api.SecurityIncidentMutationResult, error)
	GetSecurityIncidentFindings(ctx context.Context, input api.SecurityIncidentFindingsInput) (*api.SecurityIncidentFindingsResult, error)
	AddSecurityIncidentMarkers(ctx context.Context, incidentID string, markers []api.SecurityIncidentMarkerInput) (*api.AddSecurityIncidentMarkersResult, error)
	WithdrawSecurityIncidentMarkers(ctx context.Context, incidentID string, markerIDs []string) (*api.SecurityIncidentMarkersResult, error)
	PublishSecurityIncident(ctx context.Context, id string) (*api.SecurityIncidentMutationResult, error)
	ResolveSecurityIncident(ctx context.Context, id string) (*api.SecurityIncidentMutationResult, error)
	ArchiveSecurityIncident(ctx context.Context, id string) (*api.SecurityIncidentMutationResult, error)
	CreateSecurityIncidentUpdate(ctx context.Context, input api.CreateSecurityIncidentUpdateInput) (*api.CreateSecurityIncidentUpdateResult, error)
	SuppressSecurityIncidentFinding(ctx context.Context, input api.SuppressSecurityIncidentFindingInput) (*api.SuppressSecurityIncidentFindingResult, error)
	RerunSecurityIncidentImpactScan(ctx context.Context, id string) (*api.SecurityIncidentMutationResult, error)
	DryRunSecurityIncidentImpactScan(ctx context.Context, id string) (*api.DryRunSecurityIncidentImpactScanResult, error)
	GetSecurityIncidentDryRunResult(ctx context.Context, input api.SecurityIncidentDryRunResultInput) (*api.SecurityIncidentDryRunResult, error)
	ListPolicies(ctx context.Context, input api.ListPoliciesInput) (*api.PoliciesResult, error)
	GetPolicy(ctx context.Context, id string) (*api.Policy, error)
	ListPolicyResults(ctx context.Context, input api.ListPolicyResultsInput) (*api.PolicyResultsResult, error)
	GetTicketingStatus(ctx context.Context, input api.TicketingStatusInput) (*api.TicketingStatus, error)
	ListLicenses(ctx context.Context, input api.ListLicensesInput) (*api.LicensesResult, error)
	GetEnvironment(ctx context.Context, id string) (*api.Environment, error)
}

// Server is the MCP server for Lynk API
type Server struct {
	client lynkClient
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
		mcp.WithString("after", mcp.Description("Cursor for the next page")),
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

	s.mcp.AddTool(mcp.NewTool("list_doctor_results",
		mcp.WithDescription("List SBOM Doctor findings for a version"),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The UUID of the version")),
		mcp.WithString("search", mcp.Description("Case-insensitive substring search on component name")),
		mcp.WithString("component_id", mcp.Description("Filter to a single component UUID")),
		mcp.WithArray("severity", mcp.Description("Filter by Doctor severity"), mcp.WithStringItems()),
		mcp.WithArray("domain", mcp.Description("Filter by Doctor domain"), mcp.WithStringItems()),
		mcp.WithArray("check_code", mcp.Description("Filter by Doctor check code"), mcp.WithStringItems()),
		mcp.WithArray("component_name", mcp.Description("Filter by exact component name"), mcp.WithStringItems()),
		mcp.WithBoolean("force_refresh", mcp.Description("Bypass Doctor cache")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 25, max: 25)")),
		mcp.WithString("after", mcp.Description("Cursor for the next page")),
		mcp.WithNumber("last", mcp.Description("Maximum number of previous-page results to return (max: 25)")),
		mcp.WithString("before", mcp.Description("Cursor for the previous page")),
	), s.handleListDoctorResults)

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

	s.mcp.AddTool(mcp.NewTool("update_component",
		mcp.WithDescription("Destructively update component metadata. Requires confirm=true. Fetch the component first and only pass fields that should change."),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the component to update")),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The UUID of the version/SBOM containing the component")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to perform this destructive update")),
		mcp.WithString("kind", mcp.Description("Component kind")),
		mcp.WithString("name", mcp.Description("Component name")),
		mcp.WithString("description", mcp.Description("Component description")),
		mcp.WithString("copyright", mcp.Description("Component copyright")),
		mcp.WithString("version", mcp.Description("Component version")),
		mcp.WithString("group", mcp.Description("Component group")),
		mcp.WithObject("licenses", mcp.Description("License input object"), mcp.Properties(map[string]interface{}{
			"licensesExp": map[string]interface{}{"type": "string", "description": "SPDX license expression"},
		})),
		mcp.WithString("licenses_exp", mcp.Description("Convenience SPDX license expression; ignored if licenses is provided")),
		mcp.WithArray("cpes", mcp.Description("Component CPEs"), mcp.WithStringItems()),
		mcp.WithString("purl", mcp.Description("Component package URL")),
		mcp.WithBoolean("primary", mcp.Description("Whether this is the primary component")),
		mcp.WithBoolean("internal", mcp.Description("Whether this component is internal")),
		mcp.WithBoolean("generate_unique_id", mcp.Description("Generate a new component unique ID")),
		mcp.WithString("scope", mcp.Description("Component scope")),
		mcp.WithString("support_level", mcp.Description("Component support level enum: UNSPECIFIED, ACTIVELY_MAINTAINED, NO_LONGER_MAINTAINED, ABANDONED, NONE")),
		mcp.WithString("end_of_support", mcp.Description("End-of-support date or empty string")),
		mcp.WithString("notice", mcp.Description("Component notice")),
		mcp.WithArray("checksums", mcp.Description("Checksum objects with alg and content"), mcp.Items(map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"alg":     map[string]interface{}{"type": "string"},
				"content": map[string]interface{}{"type": "string"},
			},
			"required": []string{"alg", "content"},
		})),
		mcp.WithArray("external_urls", mcp.Description("External URL objects with name and url"), mcp.Items(map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{"type": "string"},
				"url":  map[string]interface{}{"type": "string"},
			},
		})),
	), s.handleUpdateComponent)

	s.mcp.AddTool(mcp.NewTool("update_component_supplier",
		mcp.WithDescription("Destructively update a component supplier. Requires confirm=true. Only pass fields that should change."),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the component supplier to update")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to perform this destructive update")),
		mcp.WithString("name", mcp.Description("Supplier name")),
		mcp.WithString("url", mcp.Description("Supplier URL")),
		mcp.WithString("contact_name", mcp.Description("Supplier contact name")),
		mcp.WithString("contact_email", mcp.Description("Supplier contact email")),
	), s.handleUpdateComponentSupplier)

	// Vulnerability tools
	s.mcp.AddTool(mcp.NewTool("list_vulnerabilities",
		mcp.WithDescription("List vulnerabilities in a version with optional filters"),
		mcp.WithString("version_id", mcp.Required(), mcp.Description("The UUID of the version")),
		mcp.WithString("severity", mcp.Description("Filter by severity (critical, high, medium, low)")),
		mcp.WithString("vex_status", mcp.Description("Filter by VEX status (e.g., affected, not_affected, fixed)")),
		mcp.WithBoolean("kev", mcp.Description("Filter to only KEV (Known Exploited Vulnerabilities)")),
		mcp.WithNumber("epss_min", mcp.Description("Filter by minimum EPSS score, e.g. 0.05 for 5%")),
		mcp.WithNumber("epss_max", mcp.Description("Filter by maximum EPSS score")),
		mcp.WithNumber("cvss_min", mcp.Description("Filter by minimum CVSS score")),
		mcp.WithNumber("cvss_max", mcp.Description("Filter by maximum CVSS score")),
		mcp.WithString("match_mode", mcp.Description("How to combine filters: all (default) or any")),
		mcp.WithBoolean("exceptional", mcp.Description("Shortcut for match_mode=any with cvss_min=9.0, epss_min=0.05, or kev=true")),
		mcp.WithString("search", mcp.Description("Search term to filter vulnerabilities")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleListVulnerabilities)

	s.mcp.AddTool(mcp.NewTool("get_vulnerability",
		mcp.WithDescription("Get details of a specific vulnerability"),
		mcp.WithString("vuln_id", mcp.Required(), mcp.Description("The CVE ID (e.g., CVE-2021-44228) or UUID")),
	), s.handleGetVulnerability)

	s.mcp.AddTool(mcp.NewTool("list_vex_statuses",
		mcp.WithDescription("List VEX statuses with UUIDs for CVE triage"),
	), s.handleListVexStatuses)

	s.mcp.AddTool(mcp.NewTool("list_vex_justifications",
		mcp.WithDescription("List VEX justifications with UUIDs for CVE triage"),
	), s.handleListVexJustifications)

	s.mcp.AddTool(mcp.NewTool("update_component_vex",
		mcp.WithDescription("Destructively update VEX data for a component vulnerability. Requires confirm=true. Only pass fields that should change. Status and justification may be supplied by UUID or by name."),
		mcp.WithString("component_vuln_id", mcp.Required(), mcp.Description("The UUID of the component vulnerability to update")),
		mcp.WithString("current_version_id", mcp.Required(), mcp.Description("The UUID of the current version/SBOM context")),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to perform this destructive update")),
		mcp.WithString("vex_status_id", mcp.Description("VEX status UUID")),
		mcp.WithString("vex_status", mcp.Description("VEX status name (e.g., affected, not_affected, fixed); resolved to UUID automatically")),
		mcp.WithString("vex_justification_id", mcp.Description("VEX justification UUID")),
		mcp.WithString("vex_justification", mcp.Description("VEX justification name (e.g., vulnerable_code_not_present); resolved to UUID automatically")),
		mcp.WithString("cdx_response_id", mcp.Description("CycloneDX response UUID")),
		mcp.WithString("note", mcp.Description("VEX note")),
		mcp.WithString("impact", mcp.Description("Impact statement")),
		mcp.WithString("detail", mcp.Description("Detail statement")),
		mcp.WithString("action", mcp.Description("Action statement")),
		mcp.WithString("fixed_in", mcp.Description("Fixed-in value")),
		mcp.WithBoolean("propagate_vex", mcp.Description("Propagate VEX to upstream")),
		mcp.WithString("resolution_date", mcp.Description("Resolution date in YYYY-MM-DD format")),
		mcp.WithArray("component_vuln_custom_field_attributes", mcp.Description("Custom field attribute objects"), mcp.Items(map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":                                   map[string]interface{}{"type": "string"},
				"componentVulnCustomFieldDefinitionId": map[string]interface{}{"type": "string"},
				"value":                                map[string]interface{}{"type": "string"},
				"_destroy":                             map[string]interface{}{"type": "boolean"},
			},
		})),
	), s.handleUpdateComponentVex)

	s.mcp.AddTool(mcp.NewTool("search_vulnerabilities",
		mcp.WithDescription("Search vulnerabilities across all products"),
		mcp.WithString("search", mcp.Description("Search term (CVE ID, component name, etc.)")),
		mcp.WithString("product_id", mcp.Description("Filter by product/repository UUID")),
		mcp.WithString("environment_id", mcp.Description("Filter by environment/project UUID")),
		mcp.WithString("severity", mcp.Description("Filter by severity")),
		mcp.WithBoolean("kev", mcp.Description("Filter to only KEV")),
		mcp.WithNumber("epss_min", mcp.Description("Filter by minimum EPSS score, e.g. 0.05 for 5%")),
		mcp.WithNumber("epss_max", mcp.Description("Filter by maximum EPSS score")),
		mcp.WithNumber("cvss_min", mcp.Description("Filter by minimum CVSS score")),
		mcp.WithNumber("cvss_max", mcp.Description("Filter by maximum CVSS score")),
		mcp.WithString("match_mode", mcp.Description("How to combine filters: all (default) or any")),
		mcp.WithBoolean("exceptional", mcp.Description("Shortcut for match_mode=any with cvss_min=9.0, epss_min=0.05, or kev=true")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results to return (default: 50)")),
	), s.handleSearchVulnerabilities)

	// Security incident tools
	s.mcp.AddTool(mcp.NewTool("list_security_incidents",
		mcp.WithDescription("List supply-chain security incidents visible to the current organization"),
		mcp.WithArray("status", mcp.Description("Filter by incident status: draft, active, resolved, archived"), mcp.WithStringItems()),
	), s.handleListSecurityIncidents)

	s.mcp.AddTool(mcp.NewTool("get_security_incident",
		mcp.WithDescription("Get a supply-chain security incident by ID, including markers and current organization impact state"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
	), s.handleGetSecurityIncident)

	s.mcp.AddTool(mcp.NewTool("create_security_incident",
		mcp.WithDescription("Create a draft supply-chain security incident. Requires operator organization permissions and confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to create the incident")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Incident title")),
		mcp.WithString("severity", mcp.Required(), mcp.Description("Incident severity: critical, high, medium, low, unknown")),
		mcp.WithString("confidence", mcp.Description("Incident confidence: confirmed, likely, investigating")),
		mcp.WithString("summary", mcp.Description("Incident summary")),
		mcp.WithString("recommended_actions", mcp.Description("Markdown recommended actions")),
		mcp.WithString("source_urls", mcp.Description("Markdown source URLs")),
	), s.handleCreateSecurityIncident)

	s.mcp.AddTool(mcp.NewTool("update_security_incident",
		mcp.WithDescription("Update editable fields on a supply-chain security incident. Requires operator organization permissions and confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to update the incident")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
		mcp.WithString("title", mcp.Description("Incident title")),
		mcp.WithString("severity", mcp.Description("Incident severity: critical, high, medium, low, unknown")),
		mcp.WithString("confidence", mcp.Description("Incident confidence: confirmed, likely, investigating")),
		mcp.WithString("summary", mcp.Description("Incident summary")),
		mcp.WithString("recommended_actions", mcp.Description("Markdown recommended actions")),
		mcp.WithString("source_urls", mcp.Description("Markdown source URLs")),
	), s.handleUpdateSecurityIncident)

	s.mcp.AddTool(mcp.NewTool("add_security_incident_markers",
		mcp.WithDescription("Add markers to a security incident. For active/resolved incidents, marker additions queue impact scanning. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to add markers")),
		mcp.WithString("security_incident_id", mcp.Required(), mcp.Description("The UUID of the security incident")),
		mcp.WithArray("markers", mcp.Required(), mcp.Description("Marker objects with marker_type and purl/component_name/component_version/github_url"), mcp.Items(map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"marker_type":       map[string]interface{}{"type": "string"},
				"purl":              map[string]interface{}{"type": "string"},
				"component_name":    map[string]interface{}{"type": "string"},
				"component_version": map[string]interface{}{"type": "string"},
				"github_url":        map[string]interface{}{"type": "string"},
			},
			"required": []string{"marker_type"},
		})),
	), s.handleAddSecurityIncidentMarkers)

	s.mcp.AddTool(mcp.NewTool("withdraw_security_incident_markers",
		mcp.WithDescription("Withdraw active markers from a security incident and resolve related active findings. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to withdraw markers")),
		mcp.WithString("security_incident_id", mcp.Required(), mcp.Description("The UUID of the security incident")),
		mcp.WithArray("marker_ids", mcp.Required(), mcp.Description("Marker UUIDs to withdraw"), mcp.WithStringItems()),
	), s.handleWithdrawSecurityIncidentMarkers)

	s.mcp.AddTool(mcp.NewTool("publish_security_incident",
		mcp.WithDescription("Publish a draft security incident and queue the initial impact scan. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to publish the incident")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
	), s.handlePublishSecurityIncident)

	s.mcp.AddTool(mcp.NewTool("resolve_security_incident",
		mcp.WithDescription("Resolve an active security incident. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to resolve the incident")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
	), s.handleResolveSecurityIncident)

	s.mcp.AddTool(mcp.NewTool("archive_security_incident",
		mcp.WithDescription("Archive a security incident. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to archive the incident")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
	), s.handleArchiveSecurityIncident)

	s.mcp.AddTool(mcp.NewTool("create_security_incident_update",
		mcp.WithDescription("Add a timeline update to a security incident. Requires operator organization permissions and confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to create the timeline update")),
		mcp.WithString("security_incident_id", mcp.Required(), mcp.Description("The UUID of the security incident")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Update title")),
		mcp.WithString("update_type", mcp.Required(), mcp.Description("Update type: indicator_added, indicator_withdrawn, guidance_changed, status_changed, source_added, correction")),
		mcp.WithString("occurred_at", mcp.Required(), mcp.Description("When the update occurred, as ISO8601 datetime")),
		mcp.WithString("body", mcp.Description("Update body")),
		mcp.WithBoolean("customer_visible", mcp.Description("Whether this update is visible to customers")),
	), s.handleCreateSecurityIncidentUpdate)

	s.mcp.AddTool(mcp.NewTool("get_security_incident_findings",
		mcp.WithDescription("Get customer-facing findings for a security incident in the current organization"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
		mcp.WithArray("statuses", mcp.Description("Filter finding statuses: active, resolved, suppressed"), mcp.WithStringItems()),
	), s.handleGetSecurityIncidentFindings)

	s.mcp.AddTool(mcp.NewTool("suppress_security_incident_finding",
		mcp.WithDescription("Suppress a specific security incident finding for the current organization. Requires confirm=true and a reason."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to suppress the finding")),
		mcp.WithString("finding_id", mcp.Required(), mcp.Description("The UUID of the finding")),
		mcp.WithString("reason", mcp.Required(), mcp.Description("Suppression reason")),
	), s.handleSuppressSecurityIncidentFinding)

	s.mcp.AddTool(mcp.NewTool("rerun_security_incident_impact_scan",
		mcp.WithDescription("Queue impact scanning for an active or resolved security incident. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to queue the scan")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
	), s.handleRerunSecurityIncidentImpactScan)

	s.mcp.AddTool(mcp.NewTool("dry_run_security_incident_impact_scan",
		mcp.WithDescription("Queue a dry-run impact scan for a security incident with active matchable markers. Requires confirm=true."),
		mcp.WithBoolean("confirm", mcp.Required(), mcp.Description("Must be true to queue the dry-run scan")),
		mcp.WithString("id", mcp.Required(), mcp.Description("The UUID of the security incident")),
	), s.handleDryRunSecurityIncidentImpactScan)

	s.mcp.AddTool(mcp.NewTool("get_security_incident_dry_run_result",
		mcp.WithDescription("Get the latest dry-run impact scan result. Operator-only. Poll no more than once every 2 seconds after queueing a dry run. Without org_id returns org summaries; with org_id returns paginated findings."),
		mcp.WithString("incident_id", mcp.Required(), mcp.Description("The UUID of the security incident")),
		mcp.WithString("org_id", mcp.Description("Organization UUID to fetch detailed dry-run findings for")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of findings to return when org_id is provided (default: 50, max: 100)")),
		mcp.WithString("after", mcp.Description("Cursor for the next findings page")),
	), s.handleGetSecurityIncidentDryRunResult)

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

	s.mcp.AddTool(mcp.NewTool("get_ticketing_status",
		mcp.WithDescription("Get ticketing provider connection/config status and policy application for products/repositories"),
		mcp.WithString("product_id", mcp.Description("Optional product UUID to inspect one product/repository")),
		mcp.WithNumber("products_limit", mcp.Description("Maximum number of products to inspect when product_id is omitted (default: 20)")),
		mcp.WithNumber("policies_limit", mcp.Description("Maximum number of policies to inspect (default: 50)")),
		mcp.WithNumber("ticket_links_limit", mcp.Description("Maximum number of component vulnerabilities to scan for created ticket links (default: 500)")),
	), s.handleGetTicketingStatus)

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
			"version:///{version_id}/doctor-results",
			"SBOM Doctor findings for a version",
			mcp.WithTemplateMIMEType("application/json"),
		),
		s.handleVersionDoctorResultsResource,
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
