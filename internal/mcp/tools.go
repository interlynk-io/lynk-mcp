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

func (s *Server) handleListProducts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	input := api.ListProductsInput{
		First: getIntParam(args, "limit", 20),
	}
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}

	result, err := s.client.ListProducts(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list products: %v", err)), nil
	}

	products := make([]map[string]interface{}, len(result.Products))
	for i, p := range result.Products {
		products[i] = map[string]interface{}{
			"id":            p.ID,
			"name":          p.Name,
			"description":   p.Description,
			"enabled":       p.Enabled,
			"versionsCount": p.VersionsCount,
			"updatedAt":     p.UpdatedAt,
		}
	}

	return formatResult(map[string]interface{}{
		"products":   products,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
	})
}

func (s *Server) handleGetProduct(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	product, err := s.client.GetProduct(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get product: %v", err)), nil
	}

	environments := make([]map[string]interface{}, len(product.Environments))
	for i, e := range product.Environments {
		environments[i] = map[string]interface{}{
			"id":            e.ID,
			"name":          e.Name,
			"description":   e.Description,
			"enabled":       e.Enabled,
			"versionsCount": e.VersionsCount,
			"updatedAt":     e.UpdatedAt,
		}
	}

	return formatResult(map[string]interface{}{
		"id":            product.ID,
		"name":          product.Name,
		"description":   product.Description,
		"enabled":       product.Enabled,
		"versionsCount": product.VersionsCount,
		"updatedAt":     product.UpdatedAt,
		"environments":  environments,
	})
}

func (s *Server) handleListEnvironments(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	productID, ok := args["product_id"].(string)
	if !ok || productID == "" {
		return newToolResultError("Missing required parameter: product_id"), nil
	}

	product, err := s.client.GetProduct(ctx, productID)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list environments: %v", err)), nil
	}

	search, _ := args["search"].(string)

	environments := make([]map[string]interface{}, 0)
	for _, e := range product.Environments {
		if search != "" && !strings.Contains(strings.ToLower(e.Name), strings.ToLower(search)) {
			continue
		}
		environments = append(environments, map[string]interface{}{
			"id":            e.ID,
			"name":          e.Name,
			"description":   e.Description,
			"enabled":       e.Enabled,
			"versionsCount": e.VersionsCount,
			"updatedAt":     e.UpdatedAt,
		})
	}

	return formatResult(map[string]interface{}{
		"environments": environments,
		"totalCount":   len(environments),
	})
}

func (s *Server) handleGetEnvironment(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	environment, err := s.client.GetEnvironment(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get environment: %v", err)), nil
	}

	result := map[string]interface{}{
		"id":            environment.ID,
		"name":          environment.Name,
		"description":   environment.Description,
		"enabled":       environment.Enabled,
		"productId":     environment.ProductID,
		"versionsCount": environment.VersionsCount,
		"updatedAt":     environment.UpdatedAt,
	}

	if environment.Product != nil {
		result["product"] = map[string]interface{}{
			"id":   environment.Product.ID,
			"name": environment.Product.Name,
		}
	}

	return formatResult(result)
}

func (s *Server) handleListVersions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	environmentID, ok := args["environment_id"].(string)
	if !ok || environmentID == "" {
		return newToolResultError("Missing required parameter: environment_id"), nil
	}

	input := api.ListVersionsInput{
		EnvironmentID: environmentID,
		First:         getIntParam(args, "limit", 20),
	}
	if lifecycle, ok := args["lifecycle"].(string); ok && lifecycle != "" {
		input.Lifecycle = []string{lifecycle}
	}

	result, err := s.client.ListVersions(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list versions: %v", err)), nil
	}

	versions := make([]map[string]interface{}, len(result.Versions))
	for i, v := range result.Versions {
		versionData := map[string]interface{}{
			"id":          v.ID,
			"version":     v.Version,
			"spec":        v.Spec,
			"specVersion": v.SpecVersion,
			"format":      v.Format,
			"lifecycle":   v.Lifecycle,
			"createdAt":   v.CreatedAt,
			"updatedAt":   v.UpdatedAt,
		}
		if v.Stats != nil {
			versionData["stats"] = map[string]interface{}{
				"componentCount": v.Stats.CompCount,
				"vulnStats":      v.Stats.VulnStats,
			}
		}
		versions[i] = versionData
	}

	return formatResult(map[string]interface{}{
		"versions":   versions,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
	})
}

func (s *Server) handleGetVersion(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	version, err := s.client.GetVersion(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get version: %v", err)), nil
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

	return formatResult(result)
}

func (s *Server) handleCompareVersions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	sourceVersionID, ok := args["source_version_id"].(string)
	if !ok || sourceVersionID == "" {
		return newToolResultError("Missing required parameter: source_version_id"), nil
	}
	targetVersionID, ok := args["target_version_id"].(string)
	if !ok || targetVersionID == "" {
		return newToolResultError("Missing required parameter: target_version_id"), nil
	}

	diffs, err := s.client.CompareVersions(ctx, sourceVersionID, targetVersionID)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to compare versions: %v", err)), nil
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
		"sourceVersionId": sourceVersionID,
		"targetVersionId": targetVersionID,
		"diffs":           result,
		"totalChanges":    len(result),
	})
}

func (s *Server) handleListComponents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	versionID, ok := args["version_id"].(string)
	if !ok || versionID == "" {
		return newToolResultError("Missing required parameter: version_id"), nil
	}

	input := api.ListComponentsInput{
		VersionID: versionID,
		First:     getIntParam(args, "limit", 50),
	}
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}
	if kind, ok := args["kind"].(string); ok && kind != "" {
		input.Kind = []string{kind}
	}
	if direct, ok := args["direct"].(bool); ok {
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
	args := toolArguments(request)
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	versionID, ok := args["version_id"].(string)
	if !ok || versionID == "" {
		return newToolResultError("Missing required parameter: version_id"), nil
	}

	component, err := s.client.GetComponent(ctx, id, versionID)
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
		"versionId":   component.VersionID,
		"updatedAt":   component.UpdatedAt,
	}

	if component.VersionInfo != nil {
		result["versionInfo"] = map[string]interface{}{
			"id":      component.VersionInfo.ID,
			"version": component.VersionInfo.Version,
		}
		if component.VersionInfo.Environment != nil {
			result["environment"] = map[string]interface{}{
				"id":   component.VersionInfo.Environment.ID,
				"name": component.VersionInfo.Environment.Name,
			}
		}
	}

	return formatResult(result)
}

func (s *Server) handleListVulnerabilities(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	versionID, ok := args["version_id"].(string)
	if !ok || versionID == "" {
		return newToolResultError("Missing required parameter: version_id"), nil
	}

	input := api.ListVersionVulnsInput{
		VersionID: versionID,
		First:     getIntParam(args, "limit", 50),
	}
	if severity, ok := args["severity"].(string); ok && severity != "" {
		input.Severity = []string{severity}
	}
	if status, ok := args["vex_status"].(string); ok && status != "" {
		input.Status = []string{status}
	}
	if kev, ok := args["kev"].(bool); ok {
		input.Kev = &kev
	}
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}

	result, err := s.client.ListVersionVulns(ctx, input)
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
	args := toolArguments(request)
	vulnID, ok := args["vuln_id"].(string)
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
	args := toolArguments(request)
	input := api.ListComponentVulnsInput{
		First: getIntParam(args, "limit", 50),
	}
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}
	if severity, ok := args["severity"].(string); ok && severity != "" {
		input.Severity = []string{severity}
	}
	if kev, ok := args["kev"].(bool); ok {
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
			"versionId": cv.VersionID,
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
	args := toolArguments(request)
	input := api.ListPoliciesInput{
		First: getIntParam(args, "limit", 20),
	}
	if search, ok := args["search"].(string); ok {
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
	args := toolArguments(request)
	id, ok := args["id"].(string)
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
			"id":       r.ID,
			"name":     r.Name,
			"subject":  r.Subject,
			"operator": r.Operator,
			"value":    r.Value,
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
	args := toolArguments(request)
	input := api.ListPolicyResultsInput{
		First: getIntParam(args, "limit", 50),
	}
	if policyID, ok := args["policy_id"].(string); ok {
		input.PolicyID = policyID
	}
	if versionID, ok := args["version_id"].(string); ok {
		input.VersionID = versionID
	}
	if resultType, ok := args["result_type"].(string); ok {
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
			"versionId":  pr.VersionID,
			"resultType": pr.ResultType,
			"result":     pr.Result,
			"createdAt":  pr.CreatedAt,
		}
		if pr.Policy != nil {
			violation["policyName"] = pr.Policy.Name
		}
		if pr.Version != nil {
			violation["version"] = pr.Version.Version
			if pr.Version.Environment != nil {
				violation["environmentName"] = pr.Version.Environment.Name
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
	args := toolArguments(request)
	input := api.ListLicensesInput{
		First: getIntParam(args, "limit", 50),
	}
	if status, ok := args["status"].(string); ok {
		input.Status = status
	}
	if search, ok := args["search"].(string); ok {
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

func toolArguments(request mcp.CallToolRequest) map[string]interface{} {
	if args, ok := request.Params.Arguments.(map[string]interface{}); ok {
		return args
	}
	return map[string]interface{}{}
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
