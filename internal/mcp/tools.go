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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
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
	if after, ok := args["after"].(string); ok {
		input.After = after
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
		if p.Repository != nil {
			products[i]["repository"] = formatImportedRepository(p.Repository)
		}
		if p.TicketingDefaultsSummary != nil {
			products[i]["ticketingDefaultsSummary"] = formatTicketingDefaultsSummary(p.TicketingDefaultsSummary)
		}
	}

	return formatResult(map[string]interface{}{
		"products":   products,
		"totalCount": result.TotalCount,
		"hasMore":    result.HasNextPage,
		"endCursor":  result.EndCursor,
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
		if e.JiraDefaults != nil {
			environments[i]["jiraDefaults"] = formatJiraDefaults(e.JiraDefaults)
		}
	}

	result := map[string]interface{}{
		"id":            product.ID,
		"name":          product.Name,
		"description":   product.Description,
		"enabled":       product.Enabled,
		"versionsCount": product.VersionsCount,
		"updatedAt":     product.UpdatedAt,
		"environments":  environments,
	}
	if product.Repository != nil {
		result["repository"] = formatImportedRepository(product.Repository)
	}

	return formatResult(result)
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
		if e.JiraDefaults != nil {
			environments[len(environments)-1]["jiraDefaults"] = formatJiraDefaults(e.JiraDefaults)
		}
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
	if environment.JiraDefaults != nil {
		result["jiraDefaults"] = formatJiraDefaults(environment.JiraDefaults)
	}

	if environment.Product != nil {
		product := map[string]interface{}{
			"id":   environment.Product.ID,
			"name": environment.Product.Name,
		}
		if environment.Product.Repository != nil {
			product["repository"] = formatImportedRepository(environment.Product.Repository)
		}
		result["product"] = product
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
		versions[i] = formatVersionSummary(&v)
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
	if includeSummary, ok := args["include_component_vuln_summary"].(bool); ok && includeSummary {
		limit := getIntParam(args, "component_summary_limit", 100)
		components, err := s.client.ListComponents(ctx, api.ListComponentsInput{
			VersionID: id,
			First:     limit,
		})
		if err != nil {
			return newToolResultError(fmt.Sprintf("Failed to list component vulnerability summaries: %v", err)), nil
		}
		result["componentVulnerabilitySummary"] = formatComponentVulnerabilitySummaries(components)
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

func (s *Server) handleFindVersion(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	versionString, ok := args["version"].(string)
	if !ok || versionString == "" {
		return newToolResultError("Missing required parameter: version"), nil
	}

	searchResult, err := s.client.SearchVersions(ctx, api.VersionSearchInput{
		Search: versionString,
		First:  getIntParam(args, "limit", 50),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to find version: %v", err)), nil
	}

	productName := strings.TrimSpace(stringParam(args, "product_name"))
	environmentName := strings.TrimSpace(stringParam(args, "environment_name"))
	matches := make([]map[string]interface{}, 0, len(searchResult.Versions))
	for _, version := range searchResult.Versions {
		if version.Version != versionString {
			continue
		}
		if productName != "" {
			if version.Environment == nil || version.Environment.Product == nil || !sameName(version.Environment.Product.Name, productName) {
				continue
			}
		}
		if environmentName != "" {
			if version.Environment == nil || !sameName(version.Environment.Name, environmentName) {
				continue
			}
		}
		matches = append(matches, formatVersionSummary(&version))
	}

	output := map[string]interface{}{
		"matches":       matches,
		"matchCount":    len(matches),
		"searchedCount": len(searchResult.Versions),
		"totalCount":    searchResult.TotalCount,
		"hasMore":       searchResult.HasNextPage,
		"endCursor":     searchResult.EndCursor,
	}
	if len(matches) == 1 {
		output["version"] = matches[0]
	}
	return formatResult(output)
}

func (s *Server) handleDownloadSBOM(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	versionID := strings.TrimSpace(stringParam(args, "version_id"))
	var resolvedVersion map[string]interface{}
	if versionID == "" {
		version, err := s.resolveDownloadVersion(ctx, args)
		if err != nil {
			return newToolResultError(err.Error()), nil
		}
		versionID = version.ID
		resolvedVersion = formatVersionSummary(version)
	}

	input := api.DownloadSBOMInput{
		VersionID:                 versionID,
		Spec:                      strings.TrimSpace(stringParam(args, "spec")),
		SpecVersion:               strings.TrimSpace(stringParam(args, "spec_version")),
		IncludeVulns:              getBoolPtrParam(args, "include_vulns"),
		IncludeFiles:              getBoolPtrParam(args, "include_files"),
		Lite:                      getBoolPtrParam(args, "lite"),
		DontPackageSBOM:           getBoolPtrParam(args, "dont_package_sbom"),
		Original:                  getBoolPtrParam(args, "original"),
		ExcludeParts:              getBoolPtrParam(args, "exclude_parts"),
		IncludeSupportStatus:      getBoolPtrParam(args, "include_support_status"),
		SupportLevelOnly:          getBoolPtrParam(args, "support_level_only"),
		RedactInternalComponents:  getBoolPtrParam(args, "redact_internal_components"),
		TLPClassificationOverride: strings.TrimSpace(stringParam(args, "tlp_classification_override")),
		RequireCompleted:          compactStrings(getStringSliceParam(args, "require_completed")),
	}
	download, err := s.client.DownloadSBOM(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to download SBOM: %v", err)), nil
	}

	includeContent := true
	if val, ok := args["include_content"].(bool); ok {
		includeContent = val
	}
	output := formatDownloadResult(download, includeContent)
	output["versionId"] = versionID
	if input.Spec != "" {
		output["spec"] = input.Spec
	}
	if input.SpecVersion != "" {
		output["specVersion"] = input.SpecVersion
	}
	if input.IncludeVulns != nil {
		output["includeVulns"] = *input.IncludeVulns
	}
	if resolvedVersion != nil {
		output["resolvedVersion"] = resolvedVersion
	}
	return formatResult(output)
}

func (s *Server) resolveDownloadVersion(ctx context.Context, args map[string]interface{}) (*api.Version, error) {
	environmentID := strings.TrimSpace(stringParam(args, "environment_id"))
	if environmentID == "" {
		productID := strings.TrimSpace(stringParam(args, "product_id"))
		if productID == "" {
			productName := strings.TrimSpace(stringParam(args, "product_name"))
			if productName == "" {
				return nil, fmt.Errorf("missing required parameter: version_id or product_id/product_name")
			}
			product, err := s.findProductByName(ctx, productName)
			if err != nil {
				return nil, err
			}
			productID = product.ID
		}

		product, err := s.client.GetProduct(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product for latest version lookup: %w", err)
		}
		environmentID, err = selectEnvironmentID(product.Environments, strings.TrimSpace(stringParam(args, "environment_name")))
		if err != nil {
			return nil, err
		}
	}

	versions, err := s.client.ListVersions(ctx, api.ListVersionsInput{
		EnvironmentID: environmentID,
		First:         getIntParam(args, "latest_version_limit", 100),
		OrderByField:  "SBOMS_UPDATED_AT",
		OrderByDir:    "DESC",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list versions for latest version lookup: %w", err)
	}
	versionString := strings.TrimSpace(stringParam(args, "version"))
	var selected *api.Version
	for i := range versions.Versions {
		version := &versions.Versions[i]
		if versionString != "" && version.Version != versionString {
			continue
		}
		if selected == nil || versionNewer(version, selected) {
			selected = version
		}
	}
	if selected == nil {
		if versionString != "" {
			return nil, fmt.Errorf("no version %q found in environment %s", versionString, environmentID)
		}
		return nil, fmt.Errorf("no versions found in environment %s", environmentID)
	}
	return selected, nil
}

func (s *Server) findProductByName(ctx context.Context, name string) (*api.Product, error) {
	products, err := s.client.ListProducts(ctx, api.ListProductsInput{Search: name, First: 50})
	if err != nil {
		return nil, fmt.Errorf("failed to find product %q: %w", name, err)
	}
	var matches []api.Product
	for _, product := range products.Products {
		if sameName(product.Name, name) {
			matches = append(matches, product)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no product found with exact name %q", name)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple products found with exact name %q; use product_id", name)
	}
	return &matches[0], nil
}

func selectEnvironmentID(environments []api.Environment, environmentName string) (string, error) {
	if environmentName != "" {
		for _, environment := range environments {
			if sameName(environment.Name, environmentName) {
				return environment.ID, nil
			}
		}
		return "", fmt.Errorf("no environment found with exact name %q", environmentName)
	}
	for _, environment := range environments {
		if sameName(environment.Name, "production") {
			return environment.ID, nil
		}
	}
	if len(environments) == 1 {
		return environments[0].ID, nil
	}
	return "", fmt.Errorf("environment_name or environment_id is required when the product has %d environments and no production environment was found", len(environments))
}

func versionNewer(candidate, current *api.Version) bool {
	if !candidate.UpdatedAt.IsZero() || !current.UpdatedAt.IsZero() {
		return candidate.UpdatedAt.After(current.UpdatedAt)
	}
	if !candidate.CreatedAt.IsZero() || !current.CreatedAt.IsZero() {
		return candidate.CreatedAt.After(current.CreatedAt)
	}
	return candidate.Version > current.Version
}

func formatDownloadResult(download *api.DownloadSBOMResult, includeContent bool) map[string]interface{} {
	output := map[string]interface{}{
		"ready":            download.Ready,
		"filename":         download.Filename,
		"contentType":      download.ContentType,
		"contentLength":    len(download.Content),
		"processingStatus": download.ProcessingStatus,
	}
	if includeContent {
		output["content"] = download.Content
	}
	return output
}

func (s *Server) handleListDoctorResults(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	versionID, ok := args["version_id"].(string)
	if !ok || versionID == "" {
		return newToolResultError("Missing required parameter: version_id"), nil
	}

	input := api.ListDoctorResultsInput{
		VersionID:     versionID,
		Severity:      getStringSliceParam(args, "severity"),
		Domain:        getStringSliceParam(args, "domain"),
		CheckCode:     getStringSliceParam(args, "check_code"),
		ComponentName: getStringSliceParam(args, "component_name"),
	}
	if limit := getIntParam(args, "limit", 0); limit > 0 {
		input.First = limit
	}
	if last := getIntParam(args, "last", 0); last > 0 {
		input.Last = last
	}
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}
	if componentID, ok := args["component_id"].(string); ok {
		input.ComponentID = componentID
	}
	if forceRefresh, ok := args["force_refresh"].(bool); ok {
		input.ForceRefresh = &forceRefresh
	}
	if after, ok := args["after"].(string); ok {
		input.After = after
	}
	if before, ok := args["before"].(string); ok {
		input.Before = before
	}

	result, err := s.client.ListDoctorResults(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list Doctor results: %v", err)), nil
	}

	return formatResult(formatDoctorResults(versionID, result))
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

func (s *Server) handleUpdateComponent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for update_component"), nil
	}

	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}
	versionID, ok := args["version_id"].(string)
	if !ok || versionID == "" {
		return newToolResultError("Missing required parameter: version_id"), nil
	}
	if !hasAnyParam(args, "kind", "name", "description", "copyright", "version", "group", "licenses", "licenses_exp", "cpes", "purl", "primary", "internal", "generate_unique_id", "scope", "support_level", "end_of_support", "notice", "checksums", "external_urls") {
		return newToolResultError("No update fields provided"), nil
	}

	input := api.UpdateComponentInput{
		ID:               id,
		VersionID:        versionID,
		Kind:             getStringPtrParam(args, "kind"),
		Name:             getStringPtrParam(args, "name"),
		Description:      getStringPtrParam(args, "description"),
		Copyright:        getStringPtrParam(args, "copyright"),
		Version:          getStringPtrParam(args, "version"),
		Group:            getStringPtrParam(args, "group"),
		Cpes:             getStringSlicePtrParam(args, "cpes"),
		Purl:             getStringPtrParam(args, "purl"),
		Primary:          getBoolPtrParam(args, "primary"),
		Internal:         getBoolPtrParam(args, "internal"),
		GenerateUniqueID: getBoolPtrParam(args, "generate_unique_id"),
		Scope:            getStringPtrParam(args, "scope"),
		SupportLevel:     getStringPtrParam(args, "support_level"),
		EndOfSupport:     getStringPtrParam(args, "end_of_support"),
		Notice:           getStringPtrParam(args, "notice"),
	}

	licenses, err := getLicenseInputParam(args)
	if err != nil {
		return newToolResultError(err.Error()), nil
	}
	input.Licenses = licenses
	checksums, err := getChecksumInputsParam(args, "checksums")
	if err != nil {
		return newToolResultError(err.Error()), nil
	}
	input.Checksums = checksums
	externalURLs, err := getExternalURLInputsParam(args, "external_urls")
	if err != nil {
		return newToolResultError(err.Error()), nil
	}
	input.ExternalURLs = externalURLs

	result, err := s.client.UpdateComponent(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to update component: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to update component: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.Component == nil {
		return newToolResultError("Failed to update component: API returned no component"), nil
	}

	return formatResult(formatComponent(result.Component))
}

func (s *Server) handleUpdateComponentSupplier(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for update_component_supplier"), nil
	}

	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}
	if !hasAnyParam(args, "name", "url", "contact_name", "contact_email") {
		return newToolResultError("No update fields provided"), nil
	}

	result, err := s.client.UpdateComponentSupplier(ctx, api.UpdateComponentSupplierInput{
		ID:           id,
		Name:         getStringPtrParam(args, "name"),
		URL:          getStringPtrParam(args, "url"),
		ContactName:  getStringPtrParam(args, "contact_name"),
		ContactEmail: getStringPtrParam(args, "contact_email"),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to update component supplier: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to update component supplier: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.Supplier == nil {
		return newToolResultError("Failed to update component supplier: API returned no supplier"), nil
	}

	return formatResult(map[string]interface{}{
		"id":           result.Supplier.ID,
		"name":         result.Supplier.Name,
		"url":          result.Supplier.URL,
		"contactName":  result.Supplier.ContactName,
		"contactEmail": result.Supplier.ContactEmail,
	})
}

func (s *Server) handleListVulnerabilities(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	versionID, ok := args["version_id"].(string)
	if !ok || versionID == "" {
		return newToolResultError("Missing required parameter: version_id"), nil
	}

	filters := getVulnerabilityMetadataFilters(args)
	input := api.ListVersionVulnsInput{
		VersionID: versionID,
		First:     getIntParam(args, "limit", 50),
	}
	if componentID, ok := args["component_id"].(string); ok && componentID != "" {
		input.ComponentIDs = []string{componentID}
	}
	if purl, ok := args["purl"].(string); ok {
		if purl == "" {
			return newToolResultError("purl must not be empty"), nil
		}
		input.Purl = purl
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
	input.EpssMin = filters.EpssMin
	input.EpssMax = filters.EpssMax
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}
	if after, ok := args["after"].(string); ok {
		input.After = after
	}

	var result *api.ComponentVulnsResult
	var matchReasons map[string][]string
	var err error
	if len(input.ComponentIDs) > 0 || input.Purl != "" {
		result, matchReasons, err = s.listVersionVulnsByComponentIdentity(ctx, input, filters)
	} else if filters.MatchAny {
		result, matchReasons, err = s.listVersionVulnsAny(ctx, input, filters)
	} else {
		query := input
		if filters.HasClientSideThresholds() {
			query.First = vulnerabilityAnyQueryLimit(input.First)
		}
		result, err = s.client.ListVersionVulns(ctx, query)
		if err == nil {
			matchReasons = matchReasonsForComponentVulns(result.ComponentVulns, filters)
			result.ComponentVulns = filterComponentVulnsByClientThresholds(result.ComponentVulns, filters)
			result.ComponentVulns = filterComponentVulnsByComponentIdentity(result.ComponentVulns, input.ComponentIDs, input.Purl)
			if filters.HasClientSideThresholds() {
				result.TotalCount = len(result.ComponentVulns)
				result.ComponentVulns = limitComponentVulns(result.ComponentVulns, input.First)
				result.HasNextPage = result.HasNextPage || result.TotalCount > len(result.ComponentVulns)
			}
		}
	}
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list vulnerabilities: %v", err)), nil
	}

	return formatResult(map[string]interface{}{
		"vulnerabilities": formatComponentVulns(result.ComponentVulns, matchReasons, true),
		"totalCount":      result.TotalCount,
		"hasMore":         result.HasNextPage,
		"endCursor":       result.EndCursor,
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

func (s *Server) handleListVexStatuses(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	statuses, err := s.client.GetVexStatuses(ctx)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list VEX statuses: %v", err)), nil
	}

	result := make([]map[string]interface{}, len(statuses))
	for i, status := range statuses {
		result[i] = map[string]interface{}{
			"id":   status.ID,
			"name": status.Name,
		}
	}

	return formatResult(map[string]interface{}{
		"vexStatuses": result,
	})
}

func (s *Server) handleListVexJustifications(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	justifications, err := s.client.GetVexJustifications(ctx)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list VEX justifications: %v", err)), nil
	}

	result := make([]map[string]interface{}, len(justifications))
	for i, justification := range justifications {
		result[i] = map[string]interface{}{
			"id":   justification.ID,
			"name": justification.Name,
		}
	}

	return formatResult(map[string]interface{}{
		"vexJustifications": result,
	})
}

func (s *Server) handleUpdateComponentVex(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for update_component_vex"), nil
	}

	componentVulnID, ok := args["component_vuln_id"].(string)
	if !ok || componentVulnID == "" {
		return newToolResultError("Missing required parameter: component_vuln_id"), nil
	}
	currentVersionID, ok := args["current_version_id"].(string)
	if !ok || currentVersionID == "" {
		return newToolResultError("Missing required parameter: current_version_id"), nil
	}
	if !hasAnyParam(args, "vex_status_id", "vex_status", "vex_justification_id", "vex_justification", "cdx_response_id", "note", "impact", "detail", "action", "fixed_in", "propagate_vex", "resolution_date", "component_vuln_custom_field_attributes") {
		return newToolResultError("No update fields provided"), nil
	}

	customFields, err := getComponentVulnCustomFieldInputsParam(args, "component_vuln_custom_field_attributes")
	if err != nil {
		return newToolResultError(err.Error()), nil
	}

	vexStatusID, err := s.resolveVexStatusID(ctx, args)
	if err != nil {
		return newToolResultError(err.Error()), nil
	}
	vexJustificationID, err := s.resolveVexJustificationID(ctx, args)
	if err != nil {
		return newToolResultError(err.Error()), nil
	}

	result, err := s.client.UpdateComponentVex(ctx, api.UpdateComponentVexInput{
		ComponentVulnID:                    componentVulnID,
		CurrentVersionID:                   currentVersionID,
		VexStatusID:                        vexStatusID,
		VexJustificationID:                 vexJustificationID,
		CDXResponseID:                      getStringPtrParam(args, "cdx_response_id"),
		Note:                               getStringPtrParam(args, "note"),
		Impact:                             getStringPtrParam(args, "impact"),
		Detail:                             getStringPtrParam(args, "detail"),
		Action:                             getStringPtrParam(args, "action"),
		FixedIn:                            getStringPtrParam(args, "fixed_in"),
		PropagateVex:                       getBoolPtrParam(args, "propagate_vex"),
		ResolutionDate:                     getStringPtrParam(args, "resolution_date"),
		ComponentVulnCustomFieldAttributes: customFields,
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to update component VEX: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to update component VEX: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.ComponentVuln == nil {
		return newToolResultError("Failed to update component VEX: API returned no component vulnerability"), nil
	}

	return formatResult(formatComponentVuln(result.ComponentVuln))
}

func (s *Server) handleBulkUpdateComponentVex(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for bulk_update_component_vex"), nil
	}

	componentVulnIDs := getStringSliceParam(args, "component_vuln_ids")
	componentVulnIDs = compactStrings(componentVulnIDs)
	if len(componentVulnIDs) == 0 {
		return newToolResultError("Missing required parameter: component_vuln_ids"), nil
	}
	if !hasAnyParam(args, "vex_status_id", "vex_status", "vex_justification_id", "vex_justification", "cdx_response_id", "note", "impact", "detail", "action", "fixed_in", "propagate_vex", "resolution_date", "component_vuln_custom_field_attributes") {
		return newToolResultError("No update fields provided"), nil
	}

	customFields, err := getComponentVulnCustomFieldInputsParam(args, "component_vuln_custom_field_attributes")
	if err != nil {
		return newToolResultError(err.Error()), nil
	}

	vexStatusID, err := s.resolveVexStatusID(ctx, args)
	if err != nil {
		return newToolResultError(err.Error()), nil
	}
	vexJustificationID, err := s.resolveVexJustificationID(ctx, args)
	if err != nil {
		return newToolResultError(err.Error()), nil
	}

	result, err := s.client.BulkUpdateComponentVex(ctx, api.BulkUpdateComponentVexInput{
		ComponentVulnIDs:                   componentVulnIDs,
		CurrentVersionID:                   getStringPtrParam(args, "current_version_id"),
		VexStatusID:                        vexStatusID,
		VexJustificationID:                 vexJustificationID,
		CDXResponseID:                      getStringPtrParam(args, "cdx_response_id"),
		Note:                               getStringPtrParam(args, "note"),
		Impact:                             getStringPtrParam(args, "impact"),
		Detail:                             getStringPtrParam(args, "detail"),
		Action:                             getStringPtrParam(args, "action"),
		FixedIn:                            getStringPtrParam(args, "fixed_in"),
		PropagateVex:                       getBoolPtrParam(args, "propagate_vex"),
		ResolutionDate:                     getStringPtrParam(args, "resolution_date"),
		ComponentVulnCustomFieldAttributes: customFields,
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to bulk update component VEX: %v", err)), nil
	}

	return formatResult(formatBulkComponentVexResult(componentVulnIDs, result))
}

func (s *Server) handleSearchVulnerabilities(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	filters := getVulnerabilityMetadataFilters(args)
	input := api.ListComponentVulnsInput{
		First: getIntParam(args, "limit", 50),
	}
	if after, ok := args["after"].(string); ok {
		input.After = after
	}
	if search, ok := args["search"].(string); ok {
		input.Search = search
	}
	input.ComponentIDs = getStringSliceParam(args, "component_ids")
	if componentID, ok := args["component_id"].(string); ok && componentID != "" {
		input.ComponentIDs = append(input.ComponentIDs, componentID)
	}
	if purl, ok := args["purl"].(string); ok {
		if purl == "" {
			return newToolResultError("purl must not be empty"), nil
		}
		input.Purl = purl
	}
	if severity, ok := args["severity"].(string); ok && severity != "" {
		input.Severity = []string{severity}
	}
	if kev, ok := args["kev"].(bool); ok {
		input.Kev = &kev
	}
	input.EpssMin = filters.EpssMin
	input.EpssMax = filters.EpssMax
	if productID, ok := args["product_id"].(string); ok && productID != "" {
		input.ProductIDs = []string{productID}
	}
	if environmentID, ok := args["environment_id"].(string); ok && environmentID != "" {
		input.EnvironmentIDs = []string{environmentID}
	}

	var result *api.ComponentVulnsResult
	var matchReasons map[string][]string
	var err error
	if filters.MatchAny {
		result, matchReasons, err = s.listComponentVulnsAny(ctx, input, filters)
	} else {
		query := input
		if filters.HasClientSideThresholds() {
			query.First = vulnerabilityAnyQueryLimit(input.First)
		}
		result, err = s.client.ListComponentVulns(ctx, query)
		if err == nil {
			matchReasons = matchReasonsForComponentVulns(result.ComponentVulns, filters)
			result.ComponentVulns = filterComponentVulnsByClientThresholds(result.ComponentVulns, filters)
			result.ComponentVulns = filterComponentVulnsByComponentIdentity(result.ComponentVulns, nil, input.Purl)
			if filters.HasClientSideThresholds() {
				result.TotalCount = len(result.ComponentVulns)
				result.ComponentVulns = limitComponentVulns(result.ComponentVulns, input.First)
				result.HasNextPage = result.HasNextPage || result.TotalCount > len(result.ComponentVulns)
			}
			if input.Purl != "" {
				result.TotalCount = len(result.ComponentVulns)
			}
		}
	}
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to search vulnerabilities: %v", err)), nil
	}

	return formatResult(map[string]interface{}{
		"vulnerabilities": formatComponentVulns(result.ComponentVulns, matchReasons, false),
		"totalCount":      result.TotalCount,
		"hasMore":         result.HasNextPage,
		"endCursor":       result.EndCursor,
	})
}

func (s *Server) handleListSecurityIncidents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	incidents, err := s.client.ListSecurityIncidents(ctx, api.ListSecurityIncidentsInput{
		Status: getStringSliceParam(args, "status"),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to list security incidents: %v", err)), nil
	}

	formatted := make([]map[string]interface{}, len(incidents))
	for i := range incidents {
		formatted[i] = formatSecurityIncident(&incidents[i])
	}

	return formatResult(map[string]interface{}{
		"securityIncidents": formatted,
		"totalCount":        len(formatted),
	})
}

func (s *Server) handleGetSecurityIncident(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	incident, err := s.client.GetSecurityIncident(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get security incident: %v", err)), nil
	}
	if incident == nil {
		return newToolResultError("Security incident not found"), nil
	}

	return formatResult(formatSecurityIncident(incident))
}

func (s *Server) handleCreateSecurityIncident(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for create_security_incident"), nil
	}

	title, ok := args["title"].(string)
	if !ok || title == "" {
		return newToolResultError("Missing required parameter: title"), nil
	}
	severity, ok := args["severity"].(string)
	if !ok || severity == "" {
		return newToolResultError("Missing required parameter: severity"), nil
	}

	result, err := s.client.CreateSecurityIncident(ctx, api.CreateSecurityIncidentInput{
		Title:              title,
		Severity:           severity,
		Confidence:         getStringPtrParam(args, "confidence"),
		Summary:            getStringPtrParam(args, "summary"),
		RecommendedActions: getStringPtrParam(args, "recommended_actions"),
		SourceURLs:         getStringPtrParam(args, "source_urls"),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to create security incident: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to create security incident: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.SecurityIncident == nil {
		return newToolResultError("Failed to create security incident: API returned no incident"), nil
	}

	return formatResult(formatSecurityIncident(result.SecurityIncident))
}

func (s *Server) handleUpdateSecurityIncident(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for update_security_incident"), nil
	}

	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}
	if !hasAnyParam(args, "title", "severity", "confidence", "summary", "recommended_actions", "source_urls") {
		return newToolResultError("At least one editable field must be provided"), nil
	}

	result, err := s.client.UpdateSecurityIncident(ctx, api.UpdateSecurityIncidentInput{
		ID:                 id,
		Title:              getStringPtrParam(args, "title"),
		Severity:           getStringPtrParam(args, "severity"),
		Confidence:         getStringPtrParam(args, "confidence"),
		Summary:            getStringPtrParam(args, "summary"),
		RecommendedActions: getStringPtrParam(args, "recommended_actions"),
		SourceURLs:         getStringPtrParam(args, "source_urls"),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to update security incident: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to update security incident: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.SecurityIncident == nil {
		return newToolResultError("Failed to update security incident: API returned no incident"), nil
	}

	return formatResult(formatSecurityIncident(result.SecurityIncident))
}

func (s *Server) handleAddSecurityIncidentMarkers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for add_security_incident_markers"), nil
	}

	incidentID, ok := args["security_incident_id"].(string)
	if !ok || incidentID == "" {
		return newToolResultError("Missing required parameter: security_incident_id"), nil
	}
	markers, err := getSecurityIncidentMarkerInputsParam(args, "markers")
	if err != nil {
		return newToolResultError(err.Error()), nil
	}
	if len(markers) == 0 {
		return newToolResultError("Missing required parameter: markers"), nil
	}

	result, err := s.client.AddSecurityIncidentMarkers(ctx, incidentID, markers)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to add security incident markers: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to add security incident markers: %s", strings.Join(result.Errors, "; "))), nil
	}

	formatted := make([]map[string]interface{}, len(result.Markers))
	for i := range result.Markers {
		formatted[i] = formatSecurityIncidentMarker(&result.Markers[i])
	}

	return formatResult(map[string]interface{}{
		"markers":      formatted,
		"totalCount":   len(formatted),
		"scanBehavior": "If the incident is active or resolved, the API queues impact scanning for the added markers.",
	})
}

func (s *Server) handleWithdrawSecurityIncidentMarkers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for withdraw_security_incident_markers"), nil
	}

	incidentID, ok := args["security_incident_id"].(string)
	if !ok || incidentID == "" {
		return newToolResultError("Missing required parameter: security_incident_id"), nil
	}
	markerIDs := getStringSliceParam(args, "marker_ids")
	if len(markerIDs) == 0 {
		return newToolResultError("Missing required parameter: marker_ids"), nil
	}

	result, err := s.client.WithdrawSecurityIncidentMarkers(ctx, incidentID, markerIDs)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to withdraw security incident markers: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to withdraw security incident markers: %s", strings.Join(result.Errors, "; "))), nil
	}

	formatted := make([]map[string]interface{}, len(result.Markers))
	for i := range result.Markers {
		formatted[i] = formatSecurityIncidentMarker(&result.Markers[i])
	}

	return formatResult(map[string]interface{}{
		"markers":      formatted,
		"totalCount":   len(formatted),
		"scanBehavior": "Active findings for withdrawn markers were resolved and organization impact state was recalculated by the API.",
	})
}

func (s *Server) handlePublishSecurityIncident(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return s.handleSecurityIncidentMutation(
		ctx,
		request,
		"publish_security_incident",
		"publish security incident",
		s.client.PublishSecurityIncident,
		"Publishing queues the initial impact scan through the incident_published event.",
	)
}

func (s *Server) handleResolveSecurityIncident(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return s.handleSecurityIncidentMutation(
		ctx,
		request,
		"resolve_security_incident",
		"resolve security incident",
		s.client.ResolveSecurityIncident,
		"Resolving emits an incident_resolved event for downstream state updates.",
	)
}

func (s *Server) handleArchiveSecurityIncident(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return s.handleSecurityIncidentMutation(
		ctx,
		request,
		"archive_security_incident",
		"archive security incident",
		s.client.ArchiveSecurityIncident,
		"Archiving emits an incident_archived event.",
	)
}

func (s *Server) handleCreateSecurityIncidentUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for create_security_incident_update"), nil
	}

	incidentID, ok := args["security_incident_id"].(string)
	if !ok || incidentID == "" {
		return newToolResultError("Missing required parameter: security_incident_id"), nil
	}
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return newToolResultError("Missing required parameter: title"), nil
	}
	updateType, ok := args["update_type"].(string)
	if !ok || updateType == "" {
		return newToolResultError("Missing required parameter: update_type"), nil
	}
	occurredAt, ok := args["occurred_at"].(string)
	if !ok || occurredAt == "" {
		return newToolResultError("Missing required parameter: occurred_at"), nil
	}

	result, err := s.client.CreateSecurityIncidentUpdate(ctx, api.CreateSecurityIncidentUpdateInput{
		SecurityIncidentID: incidentID,
		Title:              title,
		UpdateType:         updateType,
		OccurredAt:         occurredAt,
		Body:               getStringPtrParam(args, "body"),
		CustomerVisible:    getBoolPtrParam(args, "customer_visible"),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to create security incident update: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to create security incident update: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.Update == nil {
		return newToolResultError("Failed to create security incident update: API returned no update"), nil
	}

	return formatResult(formatSecurityIncidentUpdate(result.Update))
}

func (s *Server) handleGetSecurityIncidentFindings(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	result, err := s.client.GetSecurityIncidentFindings(ctx, api.SecurityIncidentFindingsInput{
		IncidentID: id,
		Statuses:   getStringSliceParam(args, "statuses"),
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get security incident findings: %v", err)), nil
	}
	if result == nil {
		return newToolResultError("Security incident not found"), nil
	}

	return formatResult(formatSecurityIncidentFindingsResult(result))
}

func (s *Server) handleSuppressSecurityIncidentFinding(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for suppress_security_incident_finding"), nil
	}

	findingID, ok := args["finding_id"].(string)
	if !ok || findingID == "" {
		return newToolResultError("Missing required parameter: finding_id"), nil
	}
	reason, ok := args["reason"].(string)
	if !ok || strings.TrimSpace(reason) == "" {
		return newToolResultError("Missing required parameter: reason"), nil
	}

	result, err := s.client.SuppressSecurityIncidentFinding(ctx, api.SuppressSecurityIncidentFindingInput{
		FindingID: findingID,
		Reason:    reason,
	})
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to suppress security incident finding: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to suppress security incident finding: %s", strings.Join(result.Errors, "; "))), nil
	}
	if result.Finding == nil {
		return newToolResultError("Failed to suppress security incident finding: API returned no finding"), nil
	}

	return formatResult(formatSecurityIncidentFinding(result.Finding))
}

func (s *Server) handleRerunSecurityIncidentImpactScan(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return s.handleSecurityIncidentMutation(
		ctx,
		request,
		"rerun_security_incident_impact_scan",
		"rerun security incident impact scan",
		s.client.RerunSecurityIncidentImpactScan,
		"Impact scanning was queued for this active or resolved incident.",
	)
}

func (s *Server) handleDryRunSecurityIncidentImpactScan(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError("Missing required confirmation: confirm must be true for dry_run_security_incident_impact_scan"), nil
	}

	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	result, err := s.client.DryRunSecurityIncidentImpactScan(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to queue dry-run impact scan: %v", err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to queue dry-run impact scan: %s", strings.Join(result.Errors, "; "))), nil
	}

	return formatResult(map[string]interface{}{
		"status":       result.Status,
		"scanBehavior": "Dry-run impact scanning was queued; poll get_security_incident_dry_run_result for status and findings no more than once every 2 seconds.",
	})
}

func (s *Server) handleGetSecurityIncidentDryRunResult(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	incidentID, ok := args["incident_id"].(string)
	if !ok || incidentID == "" {
		incidentID, _ = args["id"].(string)
	}
	if incidentID == "" {
		return newToolResultError("Missing required parameter: incident_id"), nil
	}

	input := api.SecurityIncidentDryRunResultInput{
		IncidentID: incidentID,
		First:      getIntParam(args, "limit", 50),
	}
	if input.First > 100 {
		input.First = 100
	}
	if orgID, ok := args["org_id"].(string); ok {
		input.OrgID = orgID
	}
	if after, ok := args["after"].(string); ok {
		input.After = after
	}

	result, err := s.client.GetSecurityIncidentDryRunResult(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get dry-run result: %v", err)), nil
	}

	return formatResult(formatSecurityIncidentDryRunResult(result))
}

func (s *Server) handleSecurityIncidentMutation(
	ctx context.Context,
	request mcp.CallToolRequest,
	toolName string,
	actionName string,
	mutation func(context.Context, string) (*api.SecurityIncidentMutationResult, error),
	scanBehavior string,
) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	if !confirmed(args) {
		return newToolResultError(fmt.Sprintf("Missing required confirmation: confirm must be true for %s", toolName)), nil
	}

	id, ok := args["id"].(string)
	if !ok || id == "" {
		return newToolResultError("Missing required parameter: id"), nil
	}

	result, err := mutation(ctx, id)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to %s: %v", actionName, err)), nil
	}
	if len(result.Errors) > 0 {
		return newToolResultError(fmt.Sprintf("Failed to %s: %s", actionName, strings.Join(result.Errors, "; "))), nil
	}
	if result.SecurityIncident == nil {
		return newToolResultError(fmt.Sprintf("Failed to %s: API returned no incident", actionName)), nil
	}

	formatted := formatSecurityIncident(result.SecurityIncident)
	formatted["scanBehavior"] = scanBehavior
	return formatResult(formatted)
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
			"id":           p.ID,
			"name":         p.Name,
			"description":  p.Description,
			"enabled":      p.Enabled,
			"resultType":   p.ResultType,
			"createTicket": p.CreateTicket,
			"updatedAt":    p.UpdatedAt,
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
		"id":           policy.ID,
		"name":         policy.Name,
		"description":  policy.Description,
		"enabled":      policy.Enabled,
		"resultType":   policy.ResultType,
		"createTicket": policy.CreateTicket,
		"updatedAt":    policy.UpdatedAt,
		"rules":        rules,
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

func (s *Server) handleGetTicketingStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := toolArguments(request)
	input := api.TicketingStatusInput{
		ProductsFirst: getIntParam(args, "products_limit", 20),
		PoliciesFirst: getIntParam(args, "policies_limit", 50),
		TicketsFirst:  getIntParam(args, "ticket_links_limit", 500),
	}
	if productID, ok := args["product_id"].(string); ok {
		input.ProductID = productID
	}
	if after, ok := args["products_after"].(string); ok {
		input.ProductsAfter = after
	}
	if after, ok := args["policies_after"].(string); ok {
		input.PoliciesAfter = after
	}
	if after, ok := args["ticket_links_after"].(string); ok {
		input.TicketsAfter = after
	}
	if includeCreatedTickets, ok := args["include_created_tickets"].(bool); ok {
		input.IncludeCreatedTickets = &includeCreatedTickets
	}

	status, err := s.client.GetTicketingStatus(ctx, input)
	if err != nil {
		return newToolResultError(fmt.Sprintf("Failed to get ticketing status: %v", err)), nil
	}

	return formatResult(formatTicketingStatus(status))
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

type vulnerabilityMetadataFilters struct {
	EpssMin     *float64
	EpssMax     *float64
	CvssMin     *float64
	CvssMax     *float64
	Kev         *bool
	MatchAny    bool
	Exceptional bool
}

func getVulnerabilityMetadataFilters(args map[string]interface{}) vulnerabilityMetadataFilters {
	filters := vulnerabilityMetadataFilters{
		EpssMin: getFloatPtrParam(args, "epss_min"),
		EpssMax: getFloatPtrParam(args, "epss_max"),
		CvssMin: getFloatPtrParam(args, "cvss_min"),
		CvssMax: getFloatPtrParam(args, "cvss_max"),
	}
	if kev, ok := args["kev"].(bool); ok {
		filters.Kev = &kev
	}
	if exceptional, ok := args["exceptional"].(bool); ok && exceptional {
		epssMin := 0.05
		cvssMin := 9.0
		kev := true
		filters.EpssMin = &epssMin
		filters.CvssMin = &cvssMin
		filters.Kev = &kev
		filters.MatchAny = true
		filters.Exceptional = true
	}
	if matchMode, ok := args["match_mode"].(string); ok && strings.EqualFold(matchMode, "any") {
		filters.MatchAny = true
	}
	return filters
}

func (f vulnerabilityMetadataFilters) HasClientSideThresholds() bool {
	return f.CvssMin != nil || f.CvssMax != nil
}

func (s *Server) listVersionVulnsByComponentIdentity(ctx context.Context, input api.ListVersionVulnsInput, filters vulnerabilityMetadataFilters) (*api.ComponentVulnsResult, map[string][]string, error) {
	offset, err := decodeFilteredCursor(input.After)
	if err != nil {
		return nil, nil, err
	}
	limit := input.First
	if limit <= 0 {
		limit = 50
	}

	query := input
	query.After = ""
	query.First = 100
	query.ComponentIDs = nil
	query.Purl = ""
	if filters.MatchAny {
		query.Kev = nil
		query.EpssMin = nil
		query.EpssMax = nil
	}

	var page []api.ComponentVuln
	total := 0
	for {
		result, err := s.client.ListVersionVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		if query.After != "" && result.EndCursor == query.After {
			break
		}

		filtered := result.ComponentVulns
		if filters.MatchAny {
			filtered = filterComponentVulnsByAnyMatch(filtered, filters)
		} else {
			filtered = filterComponentVulnsByClientThresholds(filtered, filters)
		}
		filtered = filterComponentVulnsByComponentIdentity(filtered, input.ComponentIDs, input.Purl)
		for _, vuln := range filtered {
			total++
			if total <= offset {
				continue
			}
			if len(page) < limit {
				page = append(page, vuln)
			}
		}

		if !result.HasNextPage || result.EndCursor == "" {
			break
		}
		query.After = result.EndCursor
	}

	consumed := offset + len(page)
	endCursor := ""
	if len(page) > 0 {
		endCursor = encodeFilteredCursor(consumed)
	}
	return &api.ComponentVulnsResult{
		ComponentVulns: page,
		TotalCount:     total,
		HasNextPage:    total > consumed,
		EndCursor:      endCursor,
	}, matchReasonsForComponentVulns(page, filters), nil
}

func (s *Server) listVersionVulnsAny(ctx context.Context, input api.ListVersionVulnsInput, filters vulnerabilityMetadataFilters) (*api.ComponentVulnsResult, map[string][]string, error) {
	merged := make(map[string]api.ComponentVuln)
	reasons := make(map[string][]string)
	hasMore := false
	queryLimit := vulnerabilityAnyQueryLimit(input.First)

	if filters.Kev != nil {
		query := input
		query.First = queryLimit
		query.Kev = filters.Kev
		query.EpssMin = nil
		query.EpssMax = nil
		result, err := s.client.ListVersionVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		mergeComponentVulns(merged, reasons, result.ComponentVulns, "kev")
		hasMore = hasMore || result.HasNextPage
	}
	if filters.EpssMin != nil || filters.EpssMax != nil {
		query := input
		query.First = queryLimit
		query.Kev = nil
		query.EpssMin = filters.EpssMin
		query.EpssMax = filters.EpssMax
		result, err := s.client.ListVersionVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		mergeComponentVulns(merged, reasons, result.ComponentVulns, "epss")
		hasMore = hasMore || result.HasNextPage
	}
	if filters.CvssMin != nil || filters.CvssMax != nil {
		query := input
		query.First = queryLimit
		query.Kev = nil
		query.EpssMin = nil
		query.EpssMax = nil
		result, err := s.client.ListVersionVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		mergeComponentVulns(merged, reasons, filterComponentVulnsByCvss(result.ComponentVulns, filters), "cvss")
		hasMore = hasMore || result.HasNextPage
	}

	filtered := filterComponentVulnsByComponentIdentity(componentVulnMapValues(merged, 0), input.ComponentIDs, input.Purl)
	vulns := limitComponentVulns(filtered, input.First)
	return &api.ComponentVulnsResult{
		ComponentVulns: vulns,
		TotalCount:     len(filtered),
		HasNextPage:    hasMore || len(filtered) > len(vulns),
	}, reasons, nil
}

func (s *Server) listComponentVulnsAny(ctx context.Context, input api.ListComponentVulnsInput, filters vulnerabilityMetadataFilters) (*api.ComponentVulnsResult, map[string][]string, error) {
	merged := make(map[string]api.ComponentVuln)
	reasons := make(map[string][]string)
	hasMore := false
	queryLimit := vulnerabilityAnyQueryLimit(input.First)

	if filters.Kev != nil {
		query := input
		query.First = queryLimit
		query.Kev = filters.Kev
		query.EpssMin = nil
		query.EpssMax = nil
		result, err := s.client.ListComponentVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		mergeComponentVulns(merged, reasons, result.ComponentVulns, "kev")
		hasMore = hasMore || result.HasNextPage
	}
	if filters.EpssMin != nil || filters.EpssMax != nil {
		query := input
		query.First = queryLimit
		query.Kev = nil
		query.EpssMin = filters.EpssMin
		query.EpssMax = filters.EpssMax
		result, err := s.client.ListComponentVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		mergeComponentVulns(merged, reasons, result.ComponentVulns, "epss")
		hasMore = hasMore || result.HasNextPage
	}
	if filters.CvssMin != nil || filters.CvssMax != nil {
		query := input
		query.First = queryLimit
		query.Kev = nil
		query.EpssMin = nil
		query.EpssMax = nil
		result, err := s.client.ListComponentVulns(ctx, query)
		if err != nil {
			return nil, nil, err
		}
		mergeComponentVulns(merged, reasons, filterComponentVulnsByCvss(result.ComponentVulns, filters), "cvss")
		hasMore = hasMore || result.HasNextPage
	}

	filtered := filterComponentVulnsByComponentIdentity(componentVulnMapValues(merged, 0), nil, input.Purl)
	vulns := limitComponentVulns(filtered, input.First)
	return &api.ComponentVulnsResult{
		ComponentVulns: vulns,
		TotalCount:     len(filtered),
		HasNextPage:    hasMore || len(filtered) > len(vulns),
	}, reasons, nil
}

func mergeComponentVulns(merged map[string]api.ComponentVuln, reasons map[string][]string, vulns []api.ComponentVuln, reason string) {
	for _, vuln := range vulns {
		merged[vuln.ID] = vuln
		if !stringSliceContains(reasons[vuln.ID], reason) {
			reasons[vuln.ID] = append(reasons[vuln.ID], reason)
		}
	}
}

func componentVulnMapValues(merged map[string]api.ComponentVuln, limit int) []api.ComponentVuln {
	if limit <= 0 {
		limit = len(merged)
	}
	vulns := make([]api.ComponentVuln, 0, minInt(len(merged), limit))
	for _, vuln := range merged {
		if len(vulns) >= limit {
			break
		}
		vulns = append(vulns, vuln)
	}
	return vulns
}

func limitComponentVulns(vulns []api.ComponentVuln, limit int) []api.ComponentVuln {
	if limit <= 0 {
		limit = 50
	}
	if len(vulns) <= limit {
		return vulns
	}
	return vulns[:limit]
}

func decodeFilteredCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	if offset, err := strconv.Atoi(cursor); err == nil && offset >= 0 {
		return offset, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid filtered cursor %q", cursor)
	}
	offset, err := strconv.Atoi(string(decoded))
	if err != nil || offset < 0 {
		return 0, fmt.Errorf("invalid filtered cursor %q", cursor)
	}
	return offset, nil
}

func encodeFilteredCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}

func vulnerabilityAnyQueryLimit(limit int) int {
	if limit <= 0 {
		limit = 50
	}
	return maxInt(limit, 500)
}

func filterComponentVulnsByClientThresholds(vulns []api.ComponentVuln, filters vulnerabilityMetadataFilters) []api.ComponentVuln {
	vulns = filterComponentVulnsByCvss(vulns, filters)
	return vulns
}

func filterComponentVulnsByAnyMatch(vulns []api.ComponentVuln, filters vulnerabilityMetadataFilters) []api.ComponentVuln {
	if !filters.MatchAny {
		return vulns
	}
	filtered := make([]api.ComponentVuln, 0, len(vulns))
	for _, vuln := range vulns {
		if len(componentVulnMatchReasons(vuln, filters)) > 0 {
			filtered = append(filtered, vuln)
		}
	}
	return filtered
}

func filterComponentVulnsByComponentIdentity(vulns []api.ComponentVuln, componentIDs []string, purl string) []api.ComponentVuln {
	if len(componentIDs) == 0 && purl == "" {
		return vulns
	}
	componentIDSet := make(map[string]bool, len(componentIDs))
	for _, id := range componentIDs {
		if id != "" {
			componentIDSet[id] = true
		}
	}

	filtered := make([]api.ComponentVuln, 0, len(vulns))
	for _, vuln := range vulns {
		if len(componentIDSet) > 0 && !componentIDSet[vuln.ComponentID] {
			continue
		}
		if purl != "" {
			if vuln.Component == nil || vuln.Component.Purl != purl {
				continue
			}
		}
		filtered = append(filtered, vuln)
	}
	return filtered
}

func filterComponentVulnsByCvss(vulns []api.ComponentVuln, filters vulnerabilityMetadataFilters) []api.ComponentVuln {
	if filters.CvssMin == nil && filters.CvssMax == nil {
		return vulns
	}
	filtered := make([]api.ComponentVuln, 0, len(vulns))
	for _, vuln := range vulns {
		if vuln.Vuln == nil {
			continue
		}
		score := vuln.Vuln.CvssScore
		if filters.CvssMin != nil && score < *filters.CvssMin {
			continue
		}
		if filters.CvssMax != nil && score > *filters.CvssMax {
			continue
		}
		filtered = append(filtered, vuln)
	}
	return filtered
}

func matchReasonsForComponentVulns(vulns []api.ComponentVuln, filters vulnerabilityMetadataFilters) map[string][]string {
	reasons := make(map[string][]string)
	for _, vuln := range vulns {
		reasons[vuln.ID] = componentVulnMatchReasons(vuln, filters)
	}
	return reasons
}

func componentVulnMatchReasons(vuln api.ComponentVuln, filters vulnerabilityMetadataFilters) []string {
	reasons := []string{}
	if filters.Kev != nil && *filters.Kev && vuln.Vuln != nil && vuln.Vuln.VulnInfo != nil && vuln.Vuln.VulnInfo.Kev {
		reasons = append(reasons, "kev")
	}
	if (filters.EpssMin != nil || filters.EpssMax != nil) && vuln.Vuln != nil && vuln.Vuln.VulnInfo != nil {
		score := vuln.Vuln.VulnInfo.EpssScore
		if (filters.EpssMin == nil || score >= *filters.EpssMin) && (filters.EpssMax == nil || score <= *filters.EpssMax) {
			reasons = append(reasons, "epss")
		}
	}
	if (filters.CvssMin != nil || filters.CvssMax != nil) && vuln.Vuln != nil {
		score := vuln.Vuln.CvssScore
		if (filters.CvssMin == nil || score >= *filters.CvssMin) && (filters.CvssMax == nil || score <= *filters.CvssMax) {
			reasons = append(reasons, "cvss")
		}
	}
	return reasons
}

func formatComponentVulns(componentVulns []api.ComponentVuln, matchReasons map[string][]string, includeDetail bool) []map[string]interface{} {
	vulns := make([]map[string]interface{}, len(componentVulns))
	for i, cv := range componentVulns {
		vuln := map[string]interface{}{
			"id":            cv.ID,
			"versionId":     cv.VersionID,
			"fixedIn":       cv.FixedIn,
			"fixedVersions": cv.FixedVersions,
			"updatedAt":     cv.UpdatedAt,
		}
		if includeDetail {
			vuln["detail"] = cv.Detail
		}
		if reasons := matchReasons[cv.ID]; len(reasons) > 0 {
			vuln["matchReasons"] = reasons
		}
		if cv.Component != nil {
			component := map[string]interface{}{
				"id":      cv.Component.ID,
				"name":    cv.Component.Name,
				"version": cv.Component.Version,
				"purl":    cv.Component.Purl,
			}
			if cv.Component.VersionID != "" {
				component["sbomId"] = cv.Component.VersionID
				component["versionId"] = cv.Component.VersionID
			}
			vuln["component"] = component
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
			vuln["vexStatusId"] = cv.VexStatus.ID
		}
		if cv.VexJustification != nil {
			vuln["vexJustification"] = cv.VexJustification.Name
			vuln["vexJustificationId"] = cv.VexJustification.ID
		}
		vulns[i] = vuln
	}
	return vulns
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func formatTicketingStatus(status *api.TicketingStatus) map[string]interface{} {
	result := map[string]interface{}{
		"productsTotalCount":  status.ProductsTotalCount,
		"productsHasMore":     status.ProductsHasNextPage,
		"productsEndCursor":   status.ProductsEndCursor,
		"policiesTotalCount":  status.PoliciesTotalCount,
		"policiesHasMore":     status.PoliciesHasNextPage,
		"policiesEndCursor":   status.PoliciesEndCursor,
		"ticketsScannedCount": status.TicketsScannedCount,
		"ticketsHasMore":      status.TicketsHasNextPage,
		"ticketsEndCursor":    status.TicketsEndCursor,
	}

	connections := make([]map[string]interface{}, len(status.Connections))
	for i, connection := range status.Connections {
		connectionData := map[string]interface{}{
			"provider":     connection.Provider,
			"connectionId": connection.ConnectionID,
			"providerId":   connection.ProviderID,
			"enabled":      connection.Enabled,
			"url":          connection.URL,
			"updatedAt":    connection.UpdatedAt,
		}
		if connection.UserName != "" {
			connectionData["userName"] = connection.UserName
		}
		if connection.HealthCheckStatus != "" {
			connectionData["healthCheckStatus"] = connection.HealthCheckStatus
		}
		if !connection.LastHealthCheckAt.IsZero() {
			connectionData["lastHealthCheckAt"] = connection.LastHealthCheckAt
		}
		connections[i] = connectionData
	}
	result["ticketingConnections"] = connections

	if status.JiraVulnManagementConfig != nil {
		result["providerConfigs"] = map[string]interface{}{
			"jiraVulnManagement": map[string]interface{}{
				"id":                 status.JiraVulnManagementConfig.ID,
				"enabled":            status.JiraVulnManagementConfig.Enabled,
				"provisioningStatus": status.JiraVulnManagementConfig.ProvisioningStatus,
				"provisioningStep":   status.JiraVulnManagementConfig.ProvisioningStep,
				"provisioningErrors": status.JiraVulnManagementConfig.ProvisioningErrors,
				"issueTypeId":        status.JiraVulnManagementConfig.IssueTypeID,
				"workflowId":         status.JiraVulnManagementConfig.WorkflowID,
				"screenId":           status.JiraVulnManagementConfig.ScreenID,
				"updatedAt":          status.JiraVulnManagementConfig.UpdatedAt,
			},
		}
	} else {
		result["providerConfigs"] = map[string]interface{}{}
	}

	products := make([]map[string]interface{}, len(status.Products))
	for i, product := range status.Products {
		productData := map[string]interface{}{
			"id":           product.ID,
			"name":         product.Name,
			"enabled":      product.Enabled,
			"repository":   formatImportedRepository(product.Repository),
			"environments": formatTicketingEnvironments(product.Environments),
		}
		products[i] = productData
	}
	result["products"] = products

	policies := make([]map[string]interface{}, len(status.Policies))
	for i, policy := range status.Policies {
		inclusions := make([]map[string]interface{}, len(policy.Inclusions))
		for j, inclusion := range policy.Inclusions {
			inclusions[j] = map[string]interface{}{
				"environmentId":   inclusion.EnvironmentID,
				"environmentName": inclusion.EnvironmentName,
				"productId":       inclusion.ProductID,
				"productName":     inclusion.ProductName,
			}
		}
		policies[i] = map[string]interface{}{
			"id":           policy.ID,
			"name":         policy.Name,
			"enabled":      policy.Enabled,
			"resultType":   policy.ResultType,
			"createTicket": policy.CreateTicket,
			"inclusions":   inclusions,
		}
	}
	result["policies"] = policies
	result["createdTickets"] = formatCreatedTickets(status.CreatedTickets)

	return result
}

func formatTicketingDefaultsSummary(summary *api.TicketingDefaultsSummary) map[string]interface{} {
	return map[string]interface{}{
		"jiraConfigured":       summary.JiraConfigured,
		"jiraProjectKeys":      summary.JiraProjectKeys,
		"environmentsWithJira": summary.EnvironmentsWithJira,
	}
}

func formatJiraDefaults(defaults *api.JiraDefaults) map[string]interface{} {
	result := map[string]interface{}{
		"id":         defaults.ID,
		"projectKey": defaults.ProjectKey,
		"issueType":  defaults.IssueType,
		"assignee":   defaults.Assignee,
		"reporter":   defaults.Reporter,
		"epic":       defaults.Epic,
		"enableSync": defaults.EnableSync,
		"updatedAt":  defaults.UpdatedAt,
	}
	if defaults.Components != nil {
		result["components"] = defaults.Components
	}
	return result
}

func formatCreatedTickets(tickets []api.CreatedTicket) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tickets))
	for i, ticket := range tickets {
		result[i] = map[string]interface{}{
			"id":               ticket.ID,
			"provider":         ticket.Provider,
			"issueKey":         ticket.IssueKey,
			"issueUrl":         ticket.IssueURL,
			"createdAt":        ticket.CreatedAt,
			"updatedAt":        ticket.UpdatedAt,
			"componentVulnId":  ticket.ComponentVulnID,
			"componentName":    ticket.ComponentName,
			"componentVersion": ticket.ComponentVersion,
			"vulnId":           ticket.VulnID,
			"vulnerabilityId":  ticket.VulnerabilityID,
			"severity":         ticket.Severity,
			"versionId":        ticket.VersionID,
			"version":          ticket.Version,
			"environmentId":    ticket.EnvironmentID,
			"environmentName":  ticket.EnvironmentName,
			"productId":        ticket.ProductID,
			"productName":      ticket.ProductName,
		}
	}
	return result
}

func formatImportedRepository(repo *api.ImportedRepository) interface{} {
	if repo == nil {
		return nil
	}

	data := map[string]interface{}{
		"type":           repo.Type,
		"id":             repo.ID,
		"name":           repo.Name,
		"importStatus":   repo.ImportStatus,
		"webhookEnabled": repo.WebhookEnabled,
	}
	if repo.FullName != "" {
		data["fullName"] = repo.FullName
	}
	if repo.Owner != "" {
		data["owner"] = repo.Owner
	}
	if repo.DefaultBranch != "" {
		data["defaultBranch"] = repo.DefaultBranch
	}
	if repo.Slug != "" {
		data["slug"] = repo.Slug
	}
	if repo.Workspace != "" {
		data["workspace"] = repo.Workspace
	}
	if repo.FullPath != "" {
		data["fullPath"] = repo.FullPath
	}
	if repo.GitlabID != "" {
		data["gitlabId"] = repo.GitlabID
	}
	return data
}

func formatTicketingEnvironments(environments []api.TicketingEnvironment) []map[string]interface{} {
	result := make([]map[string]interface{}, len(environments))
	for i, environment := range environments {
		settings := make([]map[string]interface{}, len(environment.IssueTrackerSettings))
		for j, setting := range environment.IssueTrackerSettings {
			settings[j] = map[string]interface{}{
				"id":             setting.ID,
				"provider":       setting.Provider,
				"projectKey":     setting.ProjectKey,
				"issueType":      setting.IssueType,
				"assignee":       setting.Assignee,
				"reporter":       setting.Reporter,
				"epic":           setting.Epic,
				"components":     setting.Components,
				"teamId":         setting.TeamID,
				"stateId":        setting.StateID,
				"enableSync":     setting.EnableSync,
				"lastSyncedAt":   setting.LastSyncedAt,
				"lastSyncStatus": setting.LastSyncStatus,
				"updatedAt":      setting.UpdatedAt,
			}
		}

		policies := make([]map[string]interface{}, len(environment.AppliedTicketPolicies))
		for j, policy := range environment.AppliedTicketPolicies {
			policies[j] = map[string]interface{}{
				"id":         policy.ID,
				"name":       policy.Name,
				"enabled":    policy.Enabled,
				"resultType": policy.ResultType,
			}
		}

		result[i] = map[string]interface{}{
			"id":                    environment.ID,
			"name":                  environment.Name,
			"enabled":               environment.Enabled,
			"issueTrackerSettings":  settings,
			"appliedTicketPolicies": policies,
		}
	}
	return result
}

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

func getFloatPtrParam(args map[string]interface{}, key string) *float64 {
	val, ok := args[key]
	if !ok {
		return nil
	}
	switch v := val.(type) {
	case float64:
		return &v
	case float32:
		value := float64(v)
		return &value
	case int:
		value := float64(v)
		return &value
	case int64:
		value := float64(v)
		return &value
	case json.Number:
		value, err := v.Float64()
		if err == nil {
			return &value
		}
	}
	return nil
}

func getStringSliceParam(args map[string]interface{}, key string) []string {
	val, ok := args[key]
	if !ok {
		return nil
	}
	switch v := val.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
		return result
	case string:
		if v != "" {
			return []string{v}
		}
	}
	return nil
}

func compactStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func confirmed(args map[string]interface{}) bool {
	confirm, ok := args["confirm"].(bool)
	return ok && confirm
}

func hasAnyParam(args map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := args[key]; ok {
			return true
		}
	}
	return false
}

func getStringPtrParam(args map[string]interface{}, key string) *string {
	val, ok := args[key]
	if !ok {
		return nil
	}
	str, ok := val.(string)
	if !ok {
		return nil
	}
	return &str
}

func getBoolPtrParam(args map[string]interface{}, key string) *bool {
	val, ok := args[key]
	if !ok {
		return nil
	}
	boolean, ok := val.(bool)
	if !ok {
		return nil
	}
	return &boolean
}

func (s *Server) resolveVexStatusID(ctx context.Context, args map[string]interface{}) (*string, error) {
	if id := getStringPtrParam(args, "vex_status_id"); id != nil && *id != "" {
		return id, nil
	}

	name := strings.TrimSpace(stringParam(args, "vex_status"))
	if name == "" {
		return nil, nil
	}

	statuses, err := s.client.GetVexStatuses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve VEX status %q: %w", name, err)
	}
	for _, status := range statuses {
		if sameVexName(status.Name, name) {
			id := status.ID
			return &id, nil
		}
	}

	return nil, fmt.Errorf("unknown VEX status %q; use list_vex_statuses to see supported values", name)
}

func (s *Server) resolveVexJustificationID(ctx context.Context, args map[string]interface{}) (*string, error) {
	if id := getStringPtrParam(args, "vex_justification_id"); id != nil && *id != "" {
		return id, nil
	}

	name := strings.TrimSpace(stringParam(args, "vex_justification"))
	if name == "" {
		return nil, nil
	}

	justifications, err := s.client.GetVexJustifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve VEX justification %q: %w", name, err)
	}
	for _, justification := range justifications {
		if sameVexName(justification.Name, name) {
			id := justification.ID
			return &id, nil
		}
	}

	return nil, fmt.Errorf("unknown VEX justification %q; use list_vex_justifications to see supported values", name)
}

func stringParam(args map[string]interface{}, key string) string {
	val, ok := args[key].(string)
	if !ok {
		return ""
	}
	return val
}

func sameVexName(left, right string) bool {
	return normalizeVexName(left) == normalizeVexName(right)
}

func sameName(left, right string) bool {
	return strings.EqualFold(strings.TrimSpace(left), strings.TrimSpace(right))
}

func normalizeVexName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	return name
}

func getStringSlicePtrParam(args map[string]interface{}, key string) *[]string {
	if _, ok := args[key]; !ok {
		return nil
	}
	values := getStringSliceParam(args, key)
	return &values
}

func getLicenseInputParam(args map[string]interface{}) (*api.LicenseInput, error) {
	if val, ok := args["licenses"]; ok {
		obj, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("licenses must be an object")
		}
		licensesExp, _ := obj["licensesExp"].(string)
		if licensesExp == "" {
			licensesExp, _ = obj["licenses_exp"].(string)
		}
		return &api.LicenseInput{LicensesExp: licensesExp}, nil
	}
	if licensesExp := getStringPtrParam(args, "licenses_exp"); licensesExp != nil {
		return &api.LicenseInput{LicensesExp: *licensesExp}, nil
	}
	return nil, nil
}

func getChecksumInputsParam(args map[string]interface{}, key string) (*[]api.ChecksumInput, error) {
	val, ok := args[key]
	if !ok {
		return nil, nil
	}
	items, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be an array", key)
	}
	result := make([]api.ChecksumInput, 0, len(items))
	for i, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be an object", key, i)
		}
		alg, _ := obj["alg"].(string)
		content, _ := obj["content"].(string)
		if alg == "" || content == "" {
			return nil, fmt.Errorf("%s[%d] requires alg and content", key, i)
		}
		result = append(result, api.ChecksumInput{Alg: alg, Content: content})
	}
	return &result, nil
}

func getExternalURLInputsParam(args map[string]interface{}, key string) (*[]api.ExternalURLInput, error) {
	val, ok := args[key]
	if !ok {
		return nil, nil
	}
	items, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be an array", key)
	}
	result := make([]api.ExternalURLInput, 0, len(items))
	for i, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be an object", key, i)
		}
		name, _ := obj["name"].(string)
		url, _ := obj["url"].(string)
		if name == "" && url == "" {
			return nil, fmt.Errorf("%s[%d] requires name or url", key, i)
		}
		result = append(result, api.ExternalURLInput{Name: name, URL: url})
	}
	return &result, nil
}

func getComponentVulnCustomFieldInputsParam(args map[string]interface{}, key string) (*[]api.ComponentVulnCustomFieldAttributeInput, error) {
	val, ok := args[key]
	if !ok {
		return nil, nil
	}
	items, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be an array", key)
	}
	result := make([]api.ComponentVulnCustomFieldAttributeInput, 0, len(items))
	for i, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be an object", key, i)
		}
		input := api.ComponentVulnCustomFieldAttributeInput{}
		input.ID, _ = obj["id"].(string)
		input.ComponentVulnCustomFieldDefinitionID, _ = obj["componentVulnCustomFieldDefinitionId"].(string)
		if input.ComponentVulnCustomFieldDefinitionID == "" {
			input.ComponentVulnCustomFieldDefinitionID, _ = obj["component_vuln_custom_field_definition_id"].(string)
		}
		input.Value, _ = obj["value"].(string)
		if destroy, ok := obj["_destroy"].(bool); ok {
			input.Destroy = &destroy
		} else if destroy, ok := obj["destroy"].(bool); ok {
			input.Destroy = &destroy
		}
		result = append(result, input)
	}
	return &result, nil
}

func getSecurityIncidentMarkerInputsParam(args map[string]interface{}, key string) ([]api.SecurityIncidentMarkerInput, error) {
	val, ok := args[key]
	if !ok {
		return nil, nil
	}
	items, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be an array", key)
	}
	result := make([]api.SecurityIncidentMarkerInput, 0, len(items))
	for i, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be an object", key, i)
		}

		markerType, _ := obj["marker_type"].(string)
		if markerType == "" {
			markerType, _ = obj["markerType"].(string)
		}
		if markerType == "" {
			return nil, fmt.Errorf("%s[%d] requires marker_type", key, i)
		}

		input := api.SecurityIncidentMarkerInput{MarkerType: markerType}
		input.Purl, _ = obj["purl"].(string)
		input.ComponentName, _ = obj["component_name"].(string)
		if input.ComponentName == "" {
			input.ComponentName, _ = obj["componentName"].(string)
		}
		input.ComponentVersion, _ = obj["component_version"].(string)
		if input.ComponentVersion == "" {
			input.ComponentVersion, _ = obj["componentVersion"].(string)
		}
		input.GithubURL, _ = obj["github_url"].(string)
		if input.GithubURL == "" {
			input.GithubURL, _ = obj["githubUrl"].(string)
		}
		result = append(result, input)
	}
	return result, nil
}

func formatSecurityIncident(incident *api.SecurityIncident) map[string]interface{} {
	markers := make([]map[string]interface{}, len(incident.Markers))
	for i := range incident.Markers {
		markers[i] = formatSecurityIncidentMarker(&incident.Markers[i])
	}

	result := map[string]interface{}{
		"id":                 incident.ID,
		"title":              incident.Title,
		"slug":               incident.Slug,
		"summary":            incident.Summary,
		"severity":           incident.Severity,
		"status":             incident.Status,
		"confidence":         incident.Confidence,
		"recommendedActions": incident.RecommendedActions,
		"sourceUrls":         incident.SourceURLs,
		"firstSeenAt":        incident.FirstSeenAt,
		"publishedAt":        incident.PublishedAt,
		"lastUpdatedAt":      incident.LastUpdatedAt,
		"createdAt":          incident.CreatedAt,
		"updatedAt":          incident.UpdatedAt,
		"markers":            markers,
		"activeMarkerCount":  activeSecurityIncidentMarkerCount(incident.Markers),
	}

	if incident.OrgImpactState != nil {
		result["orgImpactState"] = map[string]interface{}{
			"status":                  incident.OrgImpactState.Status,
			"severity":                incident.OrgImpactState.Severity,
			"impactedProjectsCount":   incident.OrgImpactState.ImpactedProjectsCount,
			"impactedVersionsCount":   incident.OrgImpactState.ImpactedVersionsCount,
			"impactedComponentsCount": incident.OrgImpactState.ImpactedComponentsCount,
			"lastEvaluatedAt":         incident.OrgImpactState.LastEvaluatedAt,
			"lastNotifiedAt":          incident.OrgImpactState.LastNotifiedAt,
		}
	}

	return result
}

func formatSecurityIncidentUpdate(update *api.SecurityIncidentUpdate) map[string]interface{} {
	return map[string]interface{}{
		"title":      update.Title,
		"updateType": update.UpdateType,
		"body":       update.Body,
		"occurredAt": update.OccurredAt,
	}
}

func formatSecurityIncidentMarker(marker *api.SecurityIncidentMarker) map[string]interface{} {
	return map[string]interface{}{
		"id":               marker.ID,
		"markerType":       marker.MarkerType,
		"purl":             marker.Purl,
		"componentName":    marker.ComponentName,
		"componentVersion": marker.ComponentVersion,
		"githubUrl":        marker.GithubURL,
		"active":           marker.Active,
		"addedAt":          marker.AddedAt,
		"withdrawnAt":      marker.WithdrawnAt,
	}
}

func activeSecurityIncidentMarkerCount(markers []api.SecurityIncidentMarker) int {
	count := 0
	for _, marker := range markers {
		if marker.Active {
			count++
		}
	}
	return count
}

func formatSecurityIncidentFindingsResult(result *api.SecurityIncidentFindingsResult) map[string]interface{} {
	findings := make([]map[string]interface{}, len(result.Findings))
	for i := range result.Findings {
		findings[i] = formatSecurityIncidentFinding(&result.Findings[i])
	}

	return map[string]interface{}{
		"incidentId": result.IncidentID,
		"title":      result.Title,
		"slug":       result.Slug,
		"status":     result.Status,
		"severity":   result.Severity,
		"findings":   findings,
		"totalCount": len(findings),
	}
}

func formatSecurityIncidentFinding(finding *api.SecurityIncidentFinding) map[string]interface{} {
	result := map[string]interface{}{
		"id":              finding.ID,
		"status":          finding.Status,
		"matchMethod":     finding.MatchMethod,
		"matchedFields":   finding.MatchedFields,
		"firstDetectedAt": finding.FirstDetectedAt,
		"lastConfirmedAt": finding.LastConfirmedAt,
		"isPartSbom":      finding.IsPartSbom,
	}

	if finding.Component != nil {
		result["component"] = formatSecurityIncidentFindingComponent(finding.Component)
		result["componentId"] = finding.Component.ID
		result["componentName"] = finding.Component.Name
		result["componentVersion"] = finding.Component.Version
		result["componentPurl"] = finding.Component.Purl
		result["componentSbomId"] = finding.Component.SbomID
		if finding.Component.Sbom != nil {
			result["componentSbomProjectVersion"] = finding.Component.Sbom.ProjectVersion
		}
	}

	if finding.RootSbom != nil {
		result["rootSbom"] = formatSecurityIncidentFindingSbom(finding.RootSbom)
		result["rootSbomId"] = finding.RootSbom.ID
		result["rootProjectVersion"] = finding.RootSbom.ProjectVersion
		if finding.RootSbom.Project != nil {
			result["rootProjectId"] = finding.RootSbom.Project.ID
			result["rootProjectName"] = finding.RootSbom.Project.Name
		}
	}

	return result
}

func formatSecurityIncidentFindingComponent(component *api.SecurityIncidentFindingComponent) map[string]interface{} {
	result := map[string]interface{}{
		"id":          component.ID,
		"name":        component.Name,
		"version":     component.Version,
		"kind":        component.Kind,
		"purl":        component.Purl,
		"cpes":        component.Cpes,
		"licensesExp": component.LicensesExp,
		"group":       component.Group,
		"primary":     component.Primary,
		"internal":    component.Internal,
		"sbomId":      component.SbomID,
		"updatedAt":   component.UpdatedAt,
	}
	if component.Sbom != nil {
		result["sbom"] = formatSecurityIncidentFindingSbom(component.Sbom)
	}
	return result
}

func formatSecurityIncidentFindingSbom(sbom *api.SecurityIncidentFindingSbom) map[string]interface{} {
	result := map[string]interface{}{
		"id":             sbom.ID,
		"projectVersion": sbom.ProjectVersion,
	}
	if sbom.Project != nil {
		result["project"] = map[string]interface{}{
			"id":             sbom.Project.ID,
			"name":           sbom.Project.Name,
			"projectGroupId": sbom.Project.ProjectGroupID,
		}
	}
	return result
}

func formatSecurityIncidentDryRunResult(result *api.SecurityIncidentDryRunResult) map[string]interface{} {
	orgResults := make([]map[string]interface{}, len(result.OrgResults))
	for i := range result.OrgResults {
		orgResults[i] = formatSecurityIncidentDryRunOrgResult(&result.OrgResults[i])
	}

	formatted := map[string]interface{}{
		"status":            result.Status,
		"error":             result.Error,
		"completedAt":       result.CompletedAt,
		"totalOrgsImpacted": result.TotalOrgsImpacted,
		"orgResults":        orgResults,
	}

	if result.Org != nil {
		findings := make([]map[string]interface{}, len(result.Org.Findings))
		for i := range result.Org.Findings {
			findings[i] = formatSecurityIncidentDryRunFinding(&result.Org.Findings[i])
		}
		org := formatSecurityIncidentDryRunOrgResult(&result.Org.SecurityIncidentDryRunOrgResult)
		org["findings"] = findings
		org["pageInfo"] = map[string]interface{}{
			"hasNextPage": result.Org.HasNextPage,
			"endCursor":   result.Org.EndCursor,
		}
		formatted["org"] = org
	}

	return formatted
}

func formatSecurityIncidentDryRunOrgResult(org *api.SecurityIncidentDryRunOrgResult) map[string]interface{} {
	return map[string]interface{}{
		"organizationId":          org.OrganizationID,
		"organizationName":        org.OrganizationName,
		"impactedComponentsCount": org.ImpactedComponentsCount,
		"impactedProjectsCount":   org.ImpactedProjectsCount,
		"impactedVersionsCount":   org.ImpactedVersionsCount,
	}
}

func formatSecurityIncidentDryRunFinding(finding *api.SecurityIncidentDryRunFinding) map[string]interface{} {
	return map[string]interface{}{
		"id":                     finding.ID,
		"organizationId":         finding.OrganizationID,
		"projectId":              finding.ProjectID,
		"rootSbomId":             finding.RootSbomID,
		"componentSbomId":        finding.ComponentSbomID,
		"componentId":            finding.ComponentID,
		"componentName":          finding.ComponentName,
		"componentVersion":       finding.ComponentVersion,
		"componentPurl":          finding.ComponentPurl,
		"markerId":               finding.MarkerID,
		"markerType":             finding.MarkerType,
		"markerPurl":             finding.MarkerPurl,
		"markerComponentName":    finding.MarkerComponentName,
		"markerComponentVersion": finding.MarkerComponentVersion,
		"matchMethod":            finding.MatchMethod,
		"matchedFields":          finding.MatchedFields,
		"rootProjectName":        finding.RootProjectName,
		"rootProjectVersion":     finding.RootProjectVersion,
		"isPartSbom":             finding.IsPartSbom,
	}
}

func formatComponent(component *api.VersionComponent) map[string]interface{} {
	result := map[string]interface{}{
		"id":           component.ID,
		"name":         component.Name,
		"version":      component.Version,
		"kind":         component.Kind,
		"purl":         component.Purl,
		"cpes":         component.Cpes,
		"licensesExp":  component.LicensesExp,
		"group":        component.Group,
		"description":  component.Description,
		"scope":        component.Scope,
		"copyright":    component.Copyright,
		"primary":      component.Primary,
		"internal":     component.Internal,
		"uniqueId":     component.UniqueID,
		"versionId":    component.VersionID,
		"notice":       component.Notice,
		"supportLevel": component.SupportLevel,
		"endOfSupport": component.EndOfSupport,
	}
	if !component.UpdatedAt.IsZero() {
		result["updatedAt"] = component.UpdatedAt
	}
	if component.Checksums != nil {
		checksums := make([]map[string]interface{}, len(component.Checksums))
		for i, checksum := range component.Checksums {
			checksums[i] = map[string]interface{}{
				"alg":     checksum.Alg,
				"content": checksum.Content,
			}
		}
		result["checksums"] = checksums
	}
	if component.ExternalURLs != nil {
		externalURLs := make([]map[string]interface{}, len(component.ExternalURLs))
		for i, externalURL := range component.ExternalURLs {
			externalURLs[i] = map[string]interface{}{
				"name": externalURL.Name,
				"url":  externalURL.URL,
			}
		}
		result["externalUrls"] = externalURLs
	}
	return result
}

func formatVersionSummary(version *api.Version) map[string]interface{} {
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
		environment := map[string]interface{}{
			"id":        version.Environment.ID,
			"name":      version.Environment.Name,
			"productId": version.Environment.ProductID,
		}
		if version.Environment.Product != nil {
			environment["product"] = map[string]interface{}{
				"id":   version.Environment.Product.ID,
				"name": version.Environment.Product.Name,
			}
		}
		result["environment"] = environment
	}
	return result
}

func formatComponentVulnerabilitySummaries(result *api.ComponentsResult) map[string]interface{} {
	summaries := make([]map[string]interface{}, 0, len(result.Components))
	for _, component := range result.Components {
		if component.VulnSummary == nil {
			continue
		}
		summaries = append(summaries, map[string]interface{}{
			"componentId":      component.ID,
			"componentName":    component.Name,
			"componentVersion": component.Version,
			"purl":             component.Purl,
			"totalCount":       component.VulnSummary.TotalCount,
			"severityCounts":   component.VulnSummary.Stats,
		})
	}
	return map[string]interface{}{
		"components":   summaries,
		"totalCount":   result.TotalCount,
		"hasMore":      result.HasNextPage,
		"endCursor":    result.EndCursor,
		"scannedCount": len(result.Components),
	}
}

func formatComponentVuln(componentVuln *api.ComponentVuln) map[string]interface{} {
	result := map[string]interface{}{
		"id":          componentVuln.ID,
		"componentId": componentVuln.ComponentID,
		"vulnId":      componentVuln.VulnID,
		"versionId":   componentVuln.VersionID,
		"fixedIn":     componentVuln.FixedIn,
		"detail":      componentVuln.Detail,
		"impact":      componentVuln.Impact,
		"actionStmt":  componentVuln.ActionStmt,
	}
	if componentVuln.VexStatus != nil {
		result["vexStatus"] = map[string]interface{}{
			"id":   componentVuln.VexStatus.ID,
			"name": componentVuln.VexStatus.Name,
		}
	}
	if componentVuln.VexJustification != nil {
		result["vexJustification"] = map[string]interface{}{
			"id":   componentVuln.VexJustification.ID,
			"name": componentVuln.VexJustification.Name,
		}
	}
	return result
}

func formatBulkComponentVexResult(requestedIDs []string, result *api.BulkUpdateComponentVexResult) map[string]interface{} {
	updatedByID := make(map[string]api.ComponentVuln, len(result.ComponentVulns))
	updated := make([]map[string]interface{}, len(result.ComponentVulns))
	for i, componentVuln := range result.ComponentVulns {
		updatedByID[componentVuln.ID] = componentVuln
		updated[i] = formatComponentVuln(&componentVuln)
	}

	failed := make([]map[string]interface{}, 0)
	for _, id := range requestedIDs {
		if _, ok := updatedByID[id]; ok {
			continue
		}
		failure := map[string]interface{}{
			"componentVulnId": id,
		}
		if len(result.Errors) > 0 {
			failure["errors"] = result.Errors
		} else {
			failure["errors"] = []string{"API returned no updated component vulnerability for requested ID"}
		}
		failed = append(failed, failure)
	}

	return map[string]interface{}{
		"requestedCount": len(requestedIDs),
		"updatedCount":   len(updated),
		"failedCount":    len(failed),
		"updated":        updated,
		"failed":         failed,
		"errors":         result.Errors,
	}
}

func formatDoctorResults(versionID string, result *api.DoctorResultsResult) map[string]interface{} {
	findings := make([]map[string]interface{}, len(result.Findings))
	for i, f := range result.Findings {
		findings[i] = map[string]interface{}{
			"checkCode":        f.CheckCode,
			"checkName":        f.CheckName,
			"severity":         f.Severity,
			"domain":           f.Domain,
			"componentId":      f.ComponentID,
			"componentName":    f.ComponentName,
			"componentVersion": f.ComponentVersion,
			"autoFixable":      f.AutoFixable,
			"findings":         f.Findings,
		}
	}

	return map[string]interface{}{
		"versionId":  versionID,
		"findings":   findings,
		"totalCount": result.TotalCount,
		"hasMore":    result.PageInfo.HasNextPage,
		"pageInfo": map[string]interface{}{
			"endCursor":       result.PageInfo.EndCursor,
			"hasNextPage":     result.PageInfo.HasNextPage,
			"hasPreviousPage": result.PageInfo.HasPreviousPage,
			"startCursor":     result.PageInfo.StartCursor,
		},
	}
}
