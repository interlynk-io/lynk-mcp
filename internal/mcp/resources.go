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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/interlynk-io/lynk-mcp/internal/api"
	"github.com/mark3labs/mcp-go/mcp"
)

// Resource handler implementations

func (s *Server) handleVersionResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	versionID := extractPathParam(request.Params.URI, "version_id")
	if versionID == "" {
		return nil, fmt.Errorf("missing version_id in URI")
	}

	version, err := s.client.GetVersion(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	result := map[string]interface{}{
		"id":            version.ID,
		"version":       version.Version,
		"spec":          version.Spec,
		"specVersion":   version.SpecVersion,
		"format":        version.Format,
		"lifecycle":     version.Lifecycle,
		"environmentId": version.EnvironmentID,
		"createdAt":     version.CreatedAt,
		"updatedAt":     version.UpdatedAt,
	}

	if version.Stats != nil {
		result["stats"] = map[string]interface{}{
			"componentCount":        version.Stats.CompCount,
			"componentWithPurl":     version.Stats.CompPurlCount,
			"componentWithCpe":      version.Stats.CompCpeCount,
			"componentWithLicense":  version.Stats.CompLicenseCount,
			"componentWithSupplier": version.Stats.CompSupplierCount,
			"vulnerabilities":       version.Stats.VulnStats,
		}
	}

	if version.Environment != nil {
		result["environment"] = map[string]interface{}{
			"id":        version.Environment.ID,
			"name":      version.Environment.Name,
			"productId": version.Environment.ProductID,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal version: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (s *Server) handleVersionComponentsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	versionID := extractPathParam(request.Params.URI, "version_id")
	if versionID == "" {
		return nil, fmt.Errorf("missing version_id in URI")
	}

	// Fetch all components (using a large limit)
	result, err := s.client.ListComponents(ctx, api.ListComponentsInput{
		VersionID: versionID,
		First:     1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list components: %w", err)
	}

	components := make([]map[string]interface{}, len(result.Components))
	for i, c := range result.Components {
		components[i] = map[string]interface{}{
			"id":          c.ID,
			"name":        c.Name,
			"version":     c.Version,
			"kind":        c.Kind,
			"purl":        c.Purl,
			"cpes":        c.Cpes,
			"licensesExp": c.LicensesExp,
			"group":       c.Group,
			"primary":     c.Primary,
			"internal":    c.Internal,
		}
	}

	output := map[string]interface{}{
		"versionId":  versionID,
		"components": components,
		"totalCount": result.TotalCount,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal components: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (s *Server) handleVersionVulnerabilitiesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	versionID := extractPathParam(request.Params.URI, "version_id")
	if versionID == "" {
		return nil, fmt.Errorf("missing version_id in URI")
	}

	// Fetch all vulnerabilities
	result, err := s.client.ListVersionVulns(ctx, api.ListVersionVulnsInput{
		VersionID: versionID,
		First:     1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list vulnerabilities: %w", err)
	}

	vulns := make([]map[string]interface{}, len(result.ComponentVulns))
	for i, cv := range result.ComponentVulns {
		vuln := map[string]interface{}{
			"id":      cv.ID,
			"fixedIn": cv.FixedIn,
		}
		if cv.Component != nil {
			vuln["component"] = map[string]interface{}{
				"id":      cv.Component.ID,
				"name":    cv.Component.Name,
				"version": cv.Component.Version,
				"purl":    cv.Component.Purl,
			}
		}
		if cv.Vuln != nil {
			vulnData := map[string]interface{}{
				"vulnId":      cv.Vuln.VulnID,
				"description": cv.Vuln.Description,
				"severity":    cv.Vuln.Severity,
				"cvssScore":   cv.Vuln.CvssScore,
			}
			if cv.Vuln.VulnInfo != nil {
				vulnData["kev"] = cv.Vuln.VulnInfo.Kev
				vulnData["epssScore"] = cv.Vuln.VulnInfo.EpssScore
			}
			vuln["vulnerability"] = vulnData
		}
		if cv.VexStatus != nil {
			vuln["vexStatus"] = cv.VexStatus.Name
		}
		vulns[i] = vuln
	}

	output := map[string]interface{}{
		"versionId":       versionID,
		"vulnerabilities": vulns,
		"totalCount":      result.TotalCount,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vulnerabilities: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (s *Server) handleEnvironmentLatestVersionResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	environmentID := extractPathParam(request.Params.URI, "environment_id")
	if environmentID == "" {
		return nil, fmt.Errorf("missing environment_id in URI")
	}

	// Get the most recent version for the environment
	result, err := s.client.ListVersions(ctx, api.ListVersionsInput{
		EnvironmentID: environmentID,
		First:         1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	if len(result.Versions) == 0 {
		return nil, fmt.Errorf("no versions found for environment")
	}

	version := result.Versions[0]
	output := map[string]interface{}{
		"id":          version.ID,
		"version":     version.Version,
		"spec":        version.Spec,
		"specVersion": version.SpecVersion,
		"format":      version.Format,
		"lifecycle":   version.Lifecycle,
		"createdAt":   version.CreatedAt,
		"updatedAt":   version.UpdatedAt,
	}

	if version.Stats != nil {
		output["stats"] = map[string]interface{}{
			"componentCount":  version.Stats.CompCount,
			"vulnerabilities": version.Stats.VulnStats,
		}
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal version: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (s *Server) handleOrganizationSummaryResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	org, err := s.client.GetOrganization(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	output := map[string]interface{}{
		"id":        org.ID,
		"name":      org.Name,
		"email":     org.Email,
		"url":       org.URL,
		"status":    org.Status,
		"tier":      org.Tier,
		"updatedAt": org.UpdatedAt,
	}

	if org.Metrics != nil {
		output["metrics"] = map[string]interface{}{
			"projectCount":   org.Metrics.ProjectCount,
			"versionCount":   org.Metrics.VersionCount,
			"componentCount": org.Metrics.ComponentCount,
			"vulnsMetric":    org.Metrics.VulnsMetric,
		}
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal organization: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

func (s *Server) handleVulnerabilityResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	cveID := extractPathParam(request.Params.URI, "cve_id")
	if cveID == "" {
		return nil, fmt.Errorf("missing cve_id in URI")
	}

	vuln, err := s.client.GetVuln(ctx, "", cveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerability: %w", err)
	}

	output := map[string]interface{}{
		"id":             vuln.ID,
		"vulnId":         vuln.VulnID,
		"description":    vuln.Description,
		"severity":       vuln.Severity,
		"cvssScore":      vuln.CvssScore,
		"cvssVector":     vuln.CvssVector,
		"source":         vuln.Source,
		"publishedAt":    vuln.PublishedAt,
		"lastModifiedAt": vuln.LastModifiedAt,
	}

	if vuln.VulnInfo != nil {
		output["epssScore"] = vuln.VulnInfo.EpssScore
		output["epssPercentile"] = vuln.VulnInfo.EpssPercentile
		output["kev"] = vuln.VulnInfo.Kev
		output["cwes"] = vuln.VulnInfo.Cwes
		output["advisories"] = vuln.VulnInfo.Advisories
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vulnerability: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil
}

// Helper function to extract path parameters from URI
func extractPathParam(uri, paramName string) string {
	// Parse URIs like version:///abc-123/components
	// or vulnerability:///CVE-2021-44228

	// Remove the scheme prefix
	path := uri
	if idx := strings.Index(uri, ":///"); idx != -1 {
		path = uri[idx+4:]
	}

	// Split by /
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}

	// For simple URIs like version:///{version_id}, the ID is the first part
	// For URIs like version:///{version_id}/components, the ID is still the first part
	// For URIs like vulnerability:///{cve_id}, the ID is the first part
	switch paramName {
	case "version_id", "environment_id", "cve_id":
		return parts[0]
	}

	return ""
}
