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

// Tool handler implementations

func (s *Server) handleGetOrganization(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	org, err := s.client.GetOrganization(ctx)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get organization: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":        org.ID,
		"name":      org.Name,
		"email":     org.Email,
		"url":       org.URL,
		"status":    org.Status,
		"tier":      org.Tier,
		"updatedAt": org.UpdatedAt,
	}

	if org.Metrics != nil {
		result["metrics"] = map[string]interface{}{
			"projectCount":   org.Metrics.ProjectCount,
			"versionCount":   org.Metrics.VersionCount,
			"componentCount": org.Metrics.ComponentCount,
			"vulnsMetric":    org.Metrics.VulnsMetric,
		}
	}

	return formatResult(result)
}

func (s *Server) handleListProjectGroups(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := api.ListProjectGroupsInput{
		First: getIntParam(request.Params.Arguments, "limit", 20),
	}
	if search, ok := request.Params.Arguments["search"].(string); ok {
		input.Search = search
	}

	result, err := s.client.ListProjectGroups(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list project groups: %v", err)), nil
	}

	groups := make([]map[string]interface{}, len(result.ProjectGroups))
	for i, g := range result.ProjectGroups {
		groups[i] = map[string]interface{}{
			"id":          g.ID,
			"name":        g.Name,
			"description": g.Description,
			"enabled":     g.Enabled,
			"sbomsCount":  g.SbomsCount,
			"updatedAt":   g.UpdatedAt,
		}
	}

	return formatResult(map[string]interface{}{
		"projectGroups": groups,
		"totalCount":    result.TotalCount,
		"hasMore":       result.HasNextPage,
	})
}

func (s *Server) handleGetProjectGroup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	group, err := s.client.GetProjectGroup(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get project group: %v", err)), nil
	}

	projects := make([]map[string]interface{}, len(group.Projects))
	for i, p := range group.Projects {
		projects[i] = map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"enabled":     p.Enabled,
			"sbomsCount":  p.SbomsCount,
			"updatedAt":   p.UpdatedAt,
		}
	}

	return formatResult(map[string]interface{}{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"enabled":     group.Enabled,
		"sbomsCount":  group.SbomsCount,
		"updatedAt":   group.UpdatedAt,
		"projects":    projects,
	})
}

func (s *Server) handleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectGroupID, ok := request.Params.Arguments["project_group_id"].(string)
	if !ok || projectGroupID == "" {
		return newToolResultError("Missing required parameter: project_group_id"), nil
	}

	group, err := s.client.GetProjectGroup(ctx, projectGroupID)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list projects: %v", err)), nil
	}

	search, _ := request.Params.Arguments["search"].(string)

	projects := make([]map[string]interface{}, 0)
	for _, p := range group.Projects {
		if search != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
			continue
		}
		projects = append(projects, map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"enabled":     p.Enabled,
			"sbomsCount":  p.SbomsCount,
			"updatedAt":   p.UpdatedAt,
		})
	}

	return formatResult(map[string]interface{}{
		"projects":   projects,
		"totalCount": len(projects),
	})
}

func (s *Server) handleGetProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	project, err := s.client.GetProject(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get project: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":             project.ID,
		"name":           project.Name,
		"description":    project.Description,
		"enabled":        project.Enabled,
		"projectGroupId": project.ProjectGroupID,
		"sbomsCount":     project.SbomsCount,
		"updatedAt":      project.UpdatedAt,
	}

	if project.ProjectGroup != nil {
		result["projectGroup"] = map[string]interface{}{
			"id":   project.ProjectGroup.ID,
			"name": project.ProjectGroup.Name,
		}
	}

	return formatResult(result)
}

func (s *Server) handleListSboms(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, ok := request.Params.Arguments["project_id"].(string)
	if !ok || projectID == "" {
		return newToolResultError("Missing required parameter: project_id"), nil
	}

	input := api.ListSbomsInput{
		ProjectID: projectID,
		First:     getIntParam(request.Params.Arguments, "limit", 20),
	}
	if lifecycle, ok := request.Params.Arguments["lifecycle"].(string); ok && lifecycle != "" {
		input.Lifecycle = []string{lifecycle}
	}

	result, err := s.client.ListSboms(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list SBOMs: %v", err)), nil
	}

	sboms := make([]map[string]interface{}, len(result.Sboms))
	for i, sb := range result.Sboms {
		sbomData := map[string]interface{}{
			"id":             sb.ID,
			"projectVersion": sb.ProjectVersion,
			"spec":           sb.Spec,
			"specVersion":    sb.SpecVersion,
			"format":         sb.Format,
			"lifecycle":      sb.Lifecycle,
			"createdAt":      sb.CreatedAt,
			"updatedAt":      sb.UpdatedAt,
		}
		if sb.Stats != nil {
			sbomData["stats"] = map[string]interface{}{
				"componentCount":  sb.Stats.CompCount,
				"vulnStats":       sb.Stats.VulnStats,
			}
		}
		sboms[i] = sbomData
	}

	return formatResult(map[string]interface{}{
		"sboms":      sboms,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
	})
}

func (s *Server) handleGetSbom(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	sbom, err := s.client.GetSbom(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get SBOM: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":             sbom.ID,
		"projectVersion": sbom.ProjectVersion,
		"spec":           sbom.Spec,
		"specVersion":    sbom.SpecVersion,
		"format":         sbom.Format,
		"lifecycle":      sbom.Lifecycle,
		"projectId":      sbom.ProjectID,
		"createdAt":      sbom.CreatedAt,
		"updatedAt":      sbom.UpdatedAt,
	}

	if sbom.Stats != nil {
		result["stats"] = map[string]interface{}{
			"componentCount":         sbom.Stats.CompCount,
			"componentWithPurl":      sbom.Stats.CompPurlCount,
			"componentWithCpe":       sbom.Stats.CompCpeCount,
			"componentWithLicense":   sbom.Stats.CompLicenseCount,
			"componentWithSupplier":  sbom.Stats.CompSupplierCount,
			"vulnerabilities":        sbom.Stats.VulnStats,
		}
	}

	if sbom.Project != nil {
		result["project"] = map[string]interface{}{
			"id":             sbom.Project.ID,
			"name":           sbom.Project.Name,
			"projectGroupId": sbom.Project.ProjectGroupID,
		}
	}

	return formatResult(result)
}

func (s *Server) handleCompareSboms(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sourceSbomID, ok := request.Params.Arguments["source_sbom_id"].(string)
	if !ok || sourceSbomID == "" {
		return newToolResultError("Missing required parameter: source_sbom_id"), nil
	}
	targetSbomID, ok := request.Params.Arguments["target_sbom_id"].(string)
	if !ok || targetSbomID == "" {
		return newToolResultError("Missing required parameter: target_sbom_id"), nil
	}

	diffs, err := s.client.CompareSboms(ctx, sourceSbomID, targetSbomID)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to compare SBOMs: %v", err)), nil
	}

	result := make([]map[string]interface{}, len(diffs))
	for i, d := range diffs {
		diff := map[string]interface{}{
			"diffType": d.DiffType,
			"diffTags": d.DiffTags,
		}
		if d.SubjectComponent != nil {
			diff["subjectComponent"] = map[string]interface{}{
				"id":      d.SubjectComponent.ID,
				"name":    d.SubjectComponent.Name,
				"version": d.SubjectComponent.Version,
				"purl":    d.SubjectComponent.Purl,
			}
		}
		if d.TargetComponent != nil {
			diff["targetComponent"] = map[string]interface{}{
				"id":      d.TargetComponent.ID,
				"name":    d.TargetComponent.Name,
				"version": d.TargetComponent.Version,
				"purl":    d.TargetComponent.Purl,
			}
		}
		result[i] = diff
	}

	return formatResult(map[string]interface{}{
		"sourceSbomId": sourceSbomID,
		"targetSbomId": targetSbomID,
		"diffs":        result,
		"totalChanges": len(result),
	})
}

func (s *Server) handleListComponents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sbomID, ok := request.Params.Arguments["sbom_id"].(string)
	if !ok || sbomID == "" {
		return newToolResultError("Missing required parameter: sbom_id"), nil
	}

	input := api.ListComponentsInput{
		SbomID: sbomID,
		First:  getIntParam(request.Params.Arguments, "limit", 50),
	}
	if search, ok := request.Params.Arguments["search"].(string); ok {
		input.Search = search
	}
	if kind, ok := request.Params.Arguments["kind"].(string); ok && kind != "" {
		input.Kind = []string{kind}
	}
	if direct, ok := request.Params.Arguments["direct"].(bool); ok {
		input.Direct = &direct
	}

	result, err := s.client.ListComponents(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list components: %v", err)), nil
	}

	components := make([]map[string]interface{}, len(result.Components))
	for i, c := range result.Components {
		components[i] = map[string]interface{}{
			"id":          c.ID,
			"name":        c.Name,
			"version":     c.Version,
			"kind":        c.Kind,
			"purl":        c.Purl,
			"licensesExp": c.LicensesExp,
			"primary":     c.Primary,
			"internal":    c.Internal,
		}
	}

	return formatResult(map[string]interface{}{
		"components": components,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
	})
}

func (s *Server) handleGetComponent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	sbomID, ok := request.Params.Arguments["sbom_id"].(string)
	if !ok || sbomID == "" {
		return newToolResultError("Missing required parameter: sbom_id"), nil
	}

	component, err := s.client.GetComponent(ctx, id, sbomID)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get component: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":          component.ID,
		"name":        component.Name,
		"version":     component.Version,
		"kind":        component.Kind,
		"purl":        component.Purl,
		"cpes":        component.Cpes,
		"licensesExp": component.LicensesExp,
		"group":       component.Group,
		"description": component.Description,
		"primary":     component.Primary,
		"internal":    component.Internal,
		"sbomId":      component.SbomID,
		"updatedAt":   component.UpdatedAt,
	}

	if component.Sbom != nil {
		result["sbom"] = map[string]interface{}{
			"id":             component.Sbom.ID,
			"projectVersion": component.Sbom.ProjectVersion,
		}
		if component.Sbom.Project != nil {
			result["project"] = map[string]interface{}{
				"id":   component.Sbom.Project.ID,
				"name": component.Sbom.Project.Name,
			}
		}
	}

	return formatResult(result)
}

func (s *Server) handleListVulnerabilities(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sbomID, ok := request.Params.Arguments["sbom_id"].(string)
	if !ok || sbomID == "" {
		return newToolResultError("Missing required parameter: sbom_id"), nil
	}

	input := api.ListSbomVulnsInput{
		SbomID: sbomID,
		First:  getIntParam(request.Params.Arguments, "limit", 50),
	}
	if severity, ok := request.Params.Arguments["severity"].(string); ok && severity != "" {
		input.Severity = []string{severity}
	}
	if status, ok := request.Params.Arguments["vex_status"].(string); ok && status != "" {
		input.Status = []string{status}
	}
	if kev, ok := request.Params.Arguments["kev"].(bool); ok {
		input.Kev = &kev
	}
	if search, ok := request.Params.Arguments["search"].(string); ok {
		input.Search = search
	}

	result, err := s.client.ListSbomVulns(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list vulnerabilities: %v", err)), nil
	}

	vulns := make([]map[string]interface{}, len(result.ComponentVulns))
	for i, cv := range result.ComponentVulns {
		vuln := map[string]interface{}{
			"id":        cv.ID,
			"fixedIn":   cv.FixedIn,
			"detail":    cv.Detail,
			"updatedAt": cv.UpdatedAt,
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
				"id":          cv.Vuln.ID,
				"vulnId":      cv.Vuln.VulnID,
				"description": cv.Vuln.Description,
				"severity":    cv.Vuln.Severity,
				"cvssScore":   cv.Vuln.CvssScore,
				"source":      cv.Vuln.Source,
			}
			if cv.Vuln.VulnInfo != nil {
				vulnData["epssScore"] = cv.Vuln.VulnInfo.EpssScore
				vulnData["epssPercentile"] = cv.Vuln.VulnInfo.EpssPercentile
				vulnData["kev"] = cv.Vuln.VulnInfo.Kev
				vulnData["cwes"] = cv.Vuln.VulnInfo.Cwes
			}
			vuln["vulnerability"] = vulnData
		}
		if cv.VexStatus != nil {
			vuln["vexStatus"] = cv.VexStatus.Name
		}
		if cv.VexJustification != nil {
			vuln["vexJustification"] = cv.VexJustification.Name
		}
		vulns[i] = vuln
	}

	return formatResult(map[string]interface{}{
		"vulnerabilities": vulns,
		"totalCount":      result.TotalCount,
		"hasMore":         result.HasNextPage,
	})
}

func (s *Server) handleGetVulnerability(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	vulnID, ok := request.Params.Arguments["vuln_id"].(string)
	if !ok || vulnID == "" {
		return newToolResultError("Missing required parameter: vuln_id"), nil
	}

	// Determine if it's a CVE ID or UUID
	var id, cveID string
	if strings.HasPrefix(strings.ToUpper(vulnID), "CVE-") {
		cveID = vulnID
	} else {
		id = vulnID
	}

	vuln, err := s.client.GetVuln(ctx, id, cveID)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get vulnerability: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":             vuln.ID,
		"vulnId":         vuln.VulnID,
		"description":    vuln.Description,
		"severity":       vuln.Severity,
		"cvssScore":      vuln.CvssScore,
		"cvssVector":     vuln.CvssVector,
		"source":         vuln.Source,
		"publishedAt":    vuln.PublishedAt,
		"lastModifiedAt": vuln.LastModifiedAt,
		"updatedAt":      vuln.UpdatedAt,
	}

	if vuln.VulnInfo != nil {
		result["epssScore"] = vuln.VulnInfo.EpssScore
		result["epssPercentile"] = vuln.VulnInfo.EpssPercentile
		result["kev"] = vuln.VulnInfo.Kev
		result["cwes"] = vuln.VulnInfo.Cwes
		result["advisories"] = vuln.VulnInfo.Advisories
	}

	return formatResult(result)
}

func (s *Server) handleSearchVulnerabilities(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := api.ListComponentVulnsInput{
		First: getIntParam(request.Params.Arguments, "limit", 50),
	}
	if search, ok := request.Params.Arguments["search"].(string); ok {
		input.Search = search
	}
	if severity, ok := request.Params.Arguments["severity"].(string); ok && severity != "" {
		input.Severity = []string{severity}
	}
	if kev, ok := request.Params.Arguments["kev"].(bool); ok {
		input.Kev = &kev
	}

	result, err := s.client.ListComponentVulns(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to search vulnerabilities: %v", err)), nil
	}

	vulns := make([]map[string]interface{}, len(result.ComponentVulns))
	for i, cv := range result.ComponentVulns {
		vuln := map[string]interface{}{
			"id":        cv.ID,
			"sbomId":    cv.SbomID,
			"fixedIn":   cv.FixedIn,
			"updatedAt": cv.UpdatedAt,
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
				"vulnId":    cv.Vuln.VulnID,
				"severity":  cv.Vuln.Severity,
				"cvssScore": cv.Vuln.CvssScore,
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

	return formatResult(map[string]interface{}{
		"vulnerabilities": vulns,
		"totalCount":      result.TotalCount,
		"hasMore":         result.HasNextPage,
	})
}

func (s *Server) handleListPolicies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := api.ListPoliciesInput{
		First: getIntParam(request.Params.Arguments, "limit", 20),
	}
	if search, ok := request.Params.Arguments["search"].(string); ok {
		input.Search = search
	}

	result, err := s.client.ListPolicies(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list policies: %v", err)), nil
	}

	policies := make([]map[string]interface{}, len(result.Policies))
	for i, p := range result.Policies {
		policies[i] = map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"enabled":     p.Enabled,
			"resultType":  p.ResultType,
			"updatedAt":   p.UpdatedAt,
		}
	}

	return formatResult(map[string]interface{}{
		"policies":   policies,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
	})
}

func (s *Server) handleGetPolicy(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	policy, err := s.client.GetPolicy(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get policy: %v", err)), nil
	}

	rules := make([]map[string]interface{}, len(policy.PolicyRules))
	for i, r := range policy.PolicyRules {
		rules[i] = map[string]interface{}{
			"id":          r.ID,
			"name":        r.Name,
			"subject":     r.Subject,
			"operator":    r.Operator,
			"value":       r.Value,
			"enabled":     r.Enabled,
			"failMessage": r.FailMessage,
		}
	}

	return formatResult(map[string]interface{}{
		"id":          policy.ID,
		"name":        policy.Name,
		"description": policy.Description,
		"enabled":     policy.Enabled,
		"resultType":  policy.ResultType,
		"updatedAt":   policy.UpdatedAt,
		"rules":       rules,
	})
}

func (s *Server) handleListPolicyViolations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := api.ListPolicyResultsInput{
		First: getIntParam(request.Params.Arguments, "limit", 50),
	}
	if policyID, ok := request.Params.Arguments["policy_id"].(string); ok {
		input.PolicyID = policyID
	}
	if sbomID, ok := request.Params.Arguments["sbom_id"].(string); ok {
		input.SbomID = sbomID
	}
	if resultType, ok := request.Params.Arguments["result_type"].(string); ok {
		input.ResultType = resultType
	}

	result, err := s.client.ListPolicyResults(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list policy violations: %v", err)), nil
	}

	violations := make([]map[string]interface{}, len(result.PolicyResults))
	for i, pr := range result.PolicyResults {
		violation := map[string]interface{}{
			"id":         pr.ID,
			"policyId":   pr.PolicyID,
			"sbomId":     pr.SbomID,
			"resultType": pr.ResultType,
			"result":     pr.Result,
			"createdAt":  pr.CreatedAt,
		}
		if pr.Policy != nil {
			violation["policyName"] = pr.Policy.Name
		}
		if pr.Sbom != nil {
			violation["sbomVersion"] = pr.Sbom.ProjectVersion
			if pr.Sbom.Project != nil {
				violation["projectName"] = pr.Sbom.Project.Name
			}
		}
		violations[i] = violation
	}

	return formatResult(map[string]interface{}{
		"policyResults": violations,
		"totalCount":    result.TotalCount,
		"hasMore":       result.HasNextPage,
	})
}

func (s *Server) handleListLicenses(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := api.ListLicensesInput{
		First: getIntParam(request.Params.Arguments, "limit", 50),
	}
	if status, ok := request.Params.Arguments["status"].(string); ok {
		input.Status = status
	}
	if search, ok := request.Params.Arguments["search"].(string); ok {
		input.Search = search
	}

	result, err := s.client.ListLicenses(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list licenses: %v", err)), nil
	}

	licenses := make([]map[string]interface{}, len(result.Licenses))
	for i, l := range result.Licenses {
		licenses[i] = map[string]interface{}{
			"shortId":            l.ShortID,
			"name":               l.Name,
			"state":              l.DerivedState,
			"copyLeft":           l.CopyLeft,
			"osiApproved":        l.OsiApproved,
			"fsfLibre":           l.FsfLibre,
			"deprecated":         l.Deprecated,
			"attribution":        l.Attribution,
			"sourceDistribution": l.SourceDistribution,
			"modifications":      l.Modifications,
		}
	}

	return formatResult(map[string]interface{}{
		"licenses":   licenses,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
	})
}

// Helper functions

// newToolResultError creates a CallToolResult with IsError set to true
func newToolResultError(message string) *mcp.CallToolResult {
	result := mcp.NewToolResultText(message)
	result.IsError = true
	return result
}

func formatResult(data interface{}) (*mcp.CallToolResult, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

func getIntParam(args map[string]interface{}, key string, defaultVal int) int {
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
	}
	return defaultVal
}
