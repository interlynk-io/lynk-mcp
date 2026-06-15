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
	"reflect"
	"testing"
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/api"
	mcpg "github.com/mark3labs/mcp-go/mcp"
)

type fakeLynkClient struct {
	lynkClient
	listLabelsInput             api.ListLabelsInput
	listProductsInput           api.ListProductsInput
	listVersionVulnsInput       api.ListVersionVulnsInput
	listVersionVulnsInputs      []api.ListVersionVulnsInput
	listVersionVulnsResults     map[string]*api.ComponentVulnsResult
	listComponentVulnsInput     api.ListComponentVulnsInput
	listVersionsInput           api.ListVersionsInput
	searchVersionsInput         api.VersionSearchInput
	listComponentsInput         api.ListComponentsInput
	downloadSBOMInput           api.DownloadSBOMInput
	bulkUpdateComponentVexInput api.BulkUpdateComponentVexInput
	ticketingStatusInput        api.TicketingStatusInput
}

func (f *fakeLynkClient) ListLabels(ctx context.Context, input api.ListLabelsInput) (*api.LabelsResult, error) {
	f.listLabelsInput = input
	return &api.LabelsResult{
		Labels: []api.Label{
			{
				ID:             "label-1",
				Name:           "Aidash",
				Color:          "#0052cc",
				OrganizationID: "org-1",
				CreatedAt:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		TotalCount:  2,
		HasNextPage: true,
		EndCursor:   "label-cursor-2",
	}, nil
}

func (f *fakeLynkClient) ListProducts(ctx context.Context, input api.ListProductsInput) (*api.ProductsResult, error) {
	f.listProductsInput = input
	return &api.ProductsResult{
		Products: []api.Product{
			{
				ID:            "product-1",
				Name:          "Product 1",
				Description:   "First product",
				Enabled:       true,
				VersionsCount: 7,
				Labels: []api.Label{
					{
						ID:             "label-1",
						Name:           "Aidash",
						Color:          "#0052cc",
						OrganizationID: "org-1",
					},
				},
				Repository: &api.ImportedRepository{
					Type:           "GithubRepository",
					ID:             "repo-1",
					Name:           "repo",
					FullName:       "org/repo",
					ImportStatus:   "completed",
					WebhookEnabled: true,
				},
				TicketingDefaultsSummary: &api.TicketingDefaultsSummary{
					JiraConfigured:       true,
					JiraProjectKeys:      []string{"SEC"},
					EnvironmentsWithJira: 1,
				},
			},
		},
		TotalCount:  3,
		HasNextPage: true,
		EndCursor:   "cursor-2",
	}, nil
}

func (f *fakeLynkClient) GetProduct(ctx context.Context, id string) (*api.Product, error) {
	return &api.Product{
		ID:            id,
		Name:          "Product 1",
		Description:   "First product",
		Enabled:       true,
		VersionsCount: 7,
		Labels: []api.Label{
			{
				ID:             "label-1",
				Name:           "Aidash",
				Color:          "#0052cc",
				OrganizationID: "org-1",
			},
		},
		Repository: &api.ImportedRepository{
			Type:           "BitbucketRepository",
			ID:             "repo-1",
			Name:           "repo",
			FullName:       "workspace/repo",
			ImportStatus:   "completed",
			WebhookEnabled: true,
		},
		Environments: []api.Environment{
			{
				ID:            "env-1",
				Name:          "production",
				Enabled:       true,
				VersionsCount: 1,
				JiraDefaults: &api.JiraDefaults{
					ID:         "jira-setting",
					ProjectKey: "SEC",
					IssueType:  "10001",
					EnableSync: true,
				},
			},
		},
	}, nil
}

func (f *fakeLynkClient) ListVersions(ctx context.Context, input api.ListVersionsInput) (*api.VersionsResult, error) {
	f.listVersionsInput = input
	return &api.VersionsResult{
		Versions: []api.Version{
			{
				ID:            "version-old",
				Version:       "1.0.0",
				EnvironmentID: input.EnvironmentID,
				CreatedAt:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:            "version-new",
				Version:       "1.0.1",
				EnvironmentID: input.EnvironmentID,
				CreatedAt:     time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		TotalCount:  2,
		HasNextPage: false,
	}, nil
}

func (f *fakeLynkClient) GetEnvironment(ctx context.Context, id string) (*api.Environment, error) {
	return &api.Environment{
		ID:            id,
		Name:          "default",
		Enabled:       true,
		ProductID:     "product-1",
		VersionsCount: 1,
		JiraDefaults: &api.JiraDefaults{
			ID:         "jira-setting",
			ProjectKey: "SEC",
			IssueType:  "10001",
			EnableSync: true,
		},
		Product: &api.Product{
			ID:   "product-1",
			Name: "Product 1",
			Repository: &api.ImportedRepository{
				Type:         "BitbucketRepository",
				ID:           "repo-1",
				Name:         "repo",
				FullName:     "workspace/repo",
				ImportStatus: "completed",
			},
		},
	}, nil
}

func (f *fakeLynkClient) GetVersion(ctx context.Context, id string) (*api.Version, error) {
	return &api.Version{
		ID:            id,
		Version:       "1.0.0",
		Spec:          "CycloneDX",
		SpecVersion:   "1.6",
		Format:        "json",
		Lifecycle:     "released",
		EnvironmentID: "env-1",
		Stats: &api.VersionStats{
			CompCount: 2,
			VulnStats: map[string]interface{}{
				"critical": 1,
			},
		},
		Environment: &api.Environment{
			ID:        "env-1",
			Name:      "production",
			ProductID: "product-1",
		},
	}, nil
}

func (f *fakeLynkClient) SearchVersions(ctx context.Context, input api.VersionSearchInput) (*api.VersionsResult, error) {
	f.searchVersionsInput = input
	return &api.VersionsResult{
		Versions: []api.Version{
			{
				ID:            "version-1",
				Version:       "1.0.0",
				Spec:          "CycloneDX",
				SpecVersion:   "1.6",
				Format:        "json",
				Lifecycle:     "released",
				EnvironmentID: "env-1",
				Environment: &api.Environment{
					ID:        "env-1",
					Name:      "production",
					ProductID: "product-1",
					Product: &api.Product{
						ID:   "product-1",
						Name: "Product 1",
					},
				},
			},
			{
				ID:            "version-2",
				Version:       "1.0.1",
				EnvironmentID: "env-2",
				Environment: &api.Environment{
					ID:   "env-2",
					Name: "development",
				},
			},
		},
		TotalCount:  2,
		HasNextPage: true,
		EndCursor:   "version-cursor-2",
	}, nil
}

func (f *fakeLynkClient) ListComponents(ctx context.Context, input api.ListComponentsInput) (*api.ComponentsResult, error) {
	f.listComponentsInput = input
	return &api.ComponentsResult{
		Components: []api.VersionComponent{
			{
				ID:      "component-1",
				Name:    "openssl",
				Version: "3.0.0",
				Purl:    "pkg:generic/openssl@3.0.0",
				VulnSummary: &api.ComponentVulnerabilitySummary{
					TotalCount: 2,
					Stats: map[string]interface{}{
						"high": 2,
					},
				},
			},
		},
		TotalCount:  1,
		HasNextPage: false,
	}, nil
}

func (f *fakeLynkClient) DownloadSBOM(ctx context.Context, input api.DownloadSBOMInput) (*api.DownloadSBOMResult, error) {
	f.downloadSBOMInput = input
	return &api.DownloadSBOMResult{
		Ready:       true,
		Filename:    "bom.cdx.json",
		ContentType: "application/json",
		Content:     `{"bomFormat":"CycloneDX"}`,
		ProcessingStatus: map[string]string{
			"vulnScan": "FINISHED",
		},
	}, nil
}

func (f *fakeLynkClient) ListVersionVulns(ctx context.Context, input api.ListVersionVulnsInput) (*api.ComponentVulnsResult, error) {
	f.listVersionVulnsInput = input
	f.listVersionVulnsInputs = append(f.listVersionVulnsInputs, input)
	if f.listVersionVulnsResults != nil {
		if result, ok := f.listVersionVulnsResults[input.After]; ok {
			return result, nil
		}
	}
	return &api.ComponentVulnsResult{
		ComponentVulns: []api.ComponentVuln{
			{
				ID:          "component-vuln-1",
				ComponentID: "component-1",
				VersionID:   input.VersionID,
				FixedVersions: []string{
					"3.0.1",
				},
				Component: &api.VersionComponent{
					ID:        "component-1",
					Name:      "openssl",
					Version:   "3.0.0",
					Purl:      "pkg:generic/openssl@3.0.0",
					VersionID: input.VersionID,
				},
				Vuln: &api.Vuln{
					ID:       "vuln-1",
					VulnID:   "CVE-2026-0001",
					Severity: "high",
				},
			},
			{
				ID:          "component-vuln-2",
				ComponentID: "component-2",
				VersionID:   input.VersionID,
				Component: &api.VersionComponent{
					ID:        "component-2",
					Name:      "zlib",
					Version:   "1.2.13",
					Purl:      "pkg:generic/zlib@1.2.13",
					VersionID: input.VersionID,
				},
				Vuln: &api.Vuln{
					ID:       "vuln-2",
					VulnID:   "CVE-2026-0002",
					Severity: "medium",
				},
			},
		},
		TotalCount:  4,
		HasNextPage: true,
		EndCursor:   "vuln-cursor-2",
	}, nil
}

func (f *fakeLynkClient) GetTicketingStatus(ctx context.Context, input api.TicketingStatusInput) (*api.TicketingStatus, error) {
	f.ticketingStatusInput = input
	return &api.TicketingStatus{
		ProductsTotalCount:  178,
		ProductsHasNextPage: true,
		ProductsEndCursor:   "products-cursor-2",
		PoliciesTotalCount:  6,
		PoliciesHasNextPage: true,
		PoliciesEndCursor:   "policies-cursor-2",
		TicketsScannedCount: 500,
		TicketsHasNextPage:  true,
		TicketsEndCursor:    "tickets-cursor-2",
	}, nil
}

func (f *fakeLynkClient) ListComponentVulns(ctx context.Context, input api.ListComponentVulnsInput) (*api.ComponentVulnsResult, error) {
	f.listComponentVulnsInput = input
	return &api.ComponentVulnsResult{
		ComponentVulns: []api.ComponentVuln{
			{
				ID:          "component-vuln-1",
				ComponentID: "component-1",
				VersionID:   "version-1",
				Component: &api.VersionComponent{
					ID:        "component-1",
					Name:      "openssl",
					Version:   "3.0.0",
					Purl:      "pkg:generic/openssl@3.0.0",
					VersionID: "version-1",
				},
				Vuln: &api.Vuln{
					ID:       "vuln-1",
					VulnID:   "CVE-2026-0001",
					Severity: "high",
				},
			},
		},
		TotalCount:  1,
		HasNextPage: true,
		EndCursor:   "search-cursor-2",
	}, nil
}

func (f *fakeLynkClient) GetVexStatuses(ctx context.Context) ([]api.VexStatus, error) {
	return []api.VexStatus{
		{ID: "status-1", Name: "not_affected"},
	}, nil
}

func (f *fakeLynkClient) GetVexJustifications(ctx context.Context) ([]api.VexJustification, error) {
	return []api.VexJustification{
		{ID: "justification-1", Name: "vulnerable_code_not_present"},
	}, nil
}

func (f *fakeLynkClient) BulkUpdateComponentVex(ctx context.Context, input api.BulkUpdateComponentVexInput) (*api.BulkUpdateComponentVexResult, error) {
	f.bulkUpdateComponentVexInput = input
	return &api.BulkUpdateComponentVexResult{
		ComponentVulns: []api.ComponentVuln{
			{
				ID:          "component-vuln-1",
				ComponentID: "component-1",
				VulnID:      "vuln-1",
				VersionID:   "version-1",
				FixedIn:     "1.2.3",
				VexStatus: &api.VexStatus{
					ID:   "status-1",
					Name: "not_affected",
				},
				VexJustification: &api.VexJustification{
					ID:   "justification-1",
					Name: "vulnerable_code_not_present",
				},
			},
		},
		Errors: []string{"component-vuln-2 failed"},
	}, nil
}

func TestHandleListProducts_PassesAfterAndReturnsEndCursor(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleListProducts(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"limit":  2,
				"after":  "cursor-1",
				"search": "Product",
				"label_ids": []interface{}{
					"label-1",
					"label-2",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListProducts returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListProducts returned tool error: %#v", result.Content)
	}

	if client.listProductsInput.After != "cursor-1" {
		t.Fatalf("After = %#v, want cursor-1", client.listProductsInput.After)
	}
	if client.listProductsInput.Search != "Product" {
		t.Fatalf("Search = %#v, want Product", client.listProductsInput.Search)
	}
	if client.listProductsInput.First != 2 {
		t.Fatalf("First = %#v, want 2", client.listProductsInput.First)
	}
	if !reflect.DeepEqual(client.listProductsInput.LabelIDs, []string{"label-1", "label-2"}) {
		t.Fatalf("LabelIDs = %#v, want label-1,label-2", client.listProductsInput.LabelIDs)
	}

	output := toolResultMap(t, result)
	if output["endCursor"] != "cursor-2" {
		t.Fatalf("endCursor = %#v, want cursor-2", output["endCursor"])
	}
	if output["hasMore"] != true {
		t.Fatalf("hasMore = %#v, want true", output["hasMore"])
	}
	products := output["products"].([]interface{})
	product := products[0].(map[string]interface{})
	labels := product["labels"].([]interface{})
	label := labels[0].(map[string]interface{})
	if label["id"] != "label-1" || label["name"] != "Aidash" {
		t.Fatalf("unexpected labels output: %#v", labels)
	}
	repository := product["repository"].(map[string]interface{})
	if repository["importStatus"] != "completed" || repository["fullName"] != "org/repo" {
		t.Fatalf("unexpected repository output: %#v", repository)
	}
	summary := product["ticketingDefaultsSummary"].(map[string]interface{})
	if summary["jiraConfigured"] != true {
		t.Fatalf("unexpected ticketing summary output: %#v", summary)
	}
}

func TestHandleListLabels_PassesFiltersAndReturnsLabels(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleListLabels(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"limit":  10,
				"after":  "label-cursor-1",
				"search": "Aidash",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListLabels returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListLabels returned tool error: %#v", result.Content)
	}
	if client.listLabelsInput.First != 10 || client.listLabelsInput.After != "label-cursor-1" || client.listLabelsInput.Search != "Aidash" {
		t.Fatalf("unexpected ListLabels input: %#v", client.listLabelsInput)
	}

	output := toolResultMap(t, result)
	if output["endCursor"] != "label-cursor-2" || output["hasMore"] != true {
		t.Fatalf("unexpected pagination output: %#v", output)
	}
	labels := output["labels"].([]interface{})
	label := labels[0].(map[string]interface{})
	if label["id"] != "label-1" || label["name"] != "Aidash" || label["color"] != "#0052cc" {
		t.Fatalf("unexpected labels output: %#v", labels)
	}
}

func TestHandleGetProduct_ReturnsRepositoryAndJiraDefaults(t *testing.T) {
	server := &Server{client: &fakeLynkClient{}}
	result, err := server.handleGetProduct(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{"id": "product-1"},
		},
	})
	if err != nil {
		t.Fatalf("handleGetProduct returned error: %v", err)
	}
	output := toolResultMap(t, result)
	repository := output["repository"].(map[string]interface{})
	if repository["type"] != "BitbucketRepository" || repository["importStatus"] != "completed" {
		t.Fatalf("unexpected repository output: %#v", repository)
	}
	labels := output["labels"].([]interface{})
	label := labels[0].(map[string]interface{})
	if label["id"] != "label-1" || label["name"] != "Aidash" {
		t.Fatalf("unexpected labels output: %#v", labels)
	}
	environments := output["environments"].([]interface{})
	env := environments[0].(map[string]interface{})
	defaults := env["jiraDefaults"].(map[string]interface{})
	if defaults["projectKey"] != "SEC" || defaults["enableSync"] != true {
		t.Fatalf("unexpected Jira defaults output: %#v", defaults)
	}
}

func TestHandleGetEnvironment_ReturnsJiraDefaultsAndProductRepository(t *testing.T) {
	server := &Server{client: &fakeLynkClient{}}
	result, err := server.handleGetEnvironment(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{"id": "env-1"},
		},
	})
	if err != nil {
		t.Fatalf("handleGetEnvironment returned error: %v", err)
	}
	output := toolResultMap(t, result)
	defaults := output["jiraDefaults"].(map[string]interface{})
	if defaults["projectKey"] != "SEC" || defaults["enableSync"] != true {
		t.Fatalf("unexpected Jira defaults output: %#v", defaults)
	}
	product := output["product"].(map[string]interface{})
	repository := product["repository"].(map[string]interface{})
	if repository["type"] != "BitbucketRepository" || repository["importStatus"] != "completed" {
		t.Fatalf("unexpected product repository output: %#v", repository)
	}
}

func TestHandleListVulnerabilities_PassesAfterAndReturnsEndCursor(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleListVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version_id": "version-1",
				"limit":      2,
				"after":      "vuln-cursor-1",
				"severity":   "high",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListVulnerabilities returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListVulnerabilities returned tool error: %#v", result.Content)
	}

	if client.listVersionVulnsInput.VersionID != "version-1" {
		t.Fatalf("VersionID = %#v, want version-1", client.listVersionVulnsInput.VersionID)
	}
	if client.listVersionVulnsInput.After != "vuln-cursor-1" {
		t.Fatalf("After = %#v, want vuln-cursor-1", client.listVersionVulnsInput.After)
	}
	if client.listVersionVulnsInput.First != 2 {
		t.Fatalf("First = %#v, want 2", client.listVersionVulnsInput.First)
	}
	if !reflect.DeepEqual(client.listVersionVulnsInput.Severity, []string{"high"}) {
		t.Fatalf("Severity = %#v, want high", client.listVersionVulnsInput.Severity)
	}

	output := toolResultMap(t, result)
	if output["endCursor"] != "vuln-cursor-2" {
		t.Fatalf("endCursor = %#v, want vuln-cursor-2", output["endCursor"])
	}
	if output["hasMore"] != true {
		t.Fatalf("hasMore = %#v, want true", output["hasMore"])
	}
}

func TestHandleGetVersion_IncludesComponentVulnerabilitySummary(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleGetVersion(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"id":                             "version-1",
				"include_component_vuln_summary": true,
				"component_summary_limit":        25,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleGetVersion returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleGetVersion returned tool error: %#v", result.Content)
	}
	if client.listComponentsInput.VersionID != "version-1" || client.listComponentsInput.First != 25 {
		t.Fatalf("unexpected ListComponents input: %#v", client.listComponentsInput)
	}

	output := toolResultMap(t, result)
	summary := output["componentVulnerabilitySummary"].(map[string]interface{})
	components := summary["components"].([]interface{})
	if len(components) != 1 {
		t.Fatalf("len(summary components) = %d, want 1", len(components))
	}
	component := components[0].(map[string]interface{})
	if component["componentId"] != "component-1" || component["totalCount"] != float64(2) {
		t.Fatalf("unexpected component summary: %#v", component)
	}
}

func TestHandleFindVersion_FiltersExactProductAndEnvironment(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleFindVersion(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version":          "1.0.0",
				"product_name":     "product 1",
				"environment_name": "Production",
				"limit":            10,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleFindVersion returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleFindVersion returned tool error: %#v", result.Content)
	}
	if client.searchVersionsInput.Search != "1.0.0" || client.searchVersionsInput.First != 10 {
		t.Fatalf("unexpected SearchVersions input: %#v", client.searchVersionsInput)
	}

	output := toolResultMap(t, result)
	if output["matchCount"] != float64(1) || output["searchedCount"] != float64(2) {
		t.Fatalf("unexpected find counts: %#v", output)
	}
	version := output["version"].(map[string]interface{})
	if version["id"] != "version-1" {
		t.Fatalf("unexpected matched version: %#v", version)
	}
}

func TestHandleDownloadSBOM_PassesGenericDownloadOptions(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleDownloadSBOM(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version_id":                  "version-1",
				"spec":                        "SPDX",
				"spec_version":                "2.3",
				"include_vulns":               true,
				"include_files":               true,
				"lite":                        true,
				"dont_package_sbom":           true,
				"original":                    false,
				"exclude_parts":               true,
				"include_support_status":      true,
				"support_level_only":          false,
				"redact_internal_components":  true,
				"tlp_classification_override": "AMBER",
				"require_completed":           []interface{}{"AUTOMATION", "VULN_SCAN"},
				"include_content":             false,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleDownloadSBOM returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleDownloadSBOM returned tool error: %#v", result.Content)
	}
	input := client.downloadSBOMInput
	if input.VersionID != "version-1" || input.Spec != "SPDX" || input.SpecVersion != "2.3" {
		t.Fatalf("unexpected DownloadSBOM input: %#v", input)
	}
	if input.IncludeVulns == nil || !*input.IncludeVulns || input.IncludeFiles == nil || !*input.IncludeFiles {
		t.Fatalf("expected vuln/files options to be true: %#v", input)
	}
	if input.Original == nil || *input.Original || input.SupportLevelOnly == nil || *input.SupportLevelOnly {
		t.Fatalf("expected explicit false options to be preserved: %#v", input)
	}
	if input.TLPClassificationOverride != "AMBER" || !reflect.DeepEqual(input.RequireCompleted, []string{"AUTOMATION", "VULN_SCAN"}) {
		t.Fatalf("unexpected advanced options: %#v", input)
	}

	output := toolResultMap(t, result)
	if output["versionId"] != "version-1" || output["contentLength"] != float64(len(`{"bomFormat":"CycloneDX"}`)) {
		t.Fatalf("unexpected output: %#v", output)
	}
	if _, ok := output["content"]; ok {
		t.Fatalf("content should be omitted when include_content=false: %#v", output)
	}
}

func TestHandleDownloadSBOM_ResolvesLatestVersionByProductName(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleDownloadSBOM(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"product_name":     "Product 1",
				"environment_name": "production",
				"include_content":  false,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleDownloadSBOM returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleDownloadSBOM returned tool error: %#v", result.Content)
	}
	if client.listProductsInput.Search != "Product 1" {
		t.Fatalf("unexpected ListProducts input: %#v", client.listProductsInput)
	}
	if client.downloadSBOMInput.VersionID != "version-new" {
		t.Fatalf("download version = %q, want latest version-new", client.downloadSBOMInput.VersionID)
	}
	if client.listVersionsInput.OrderByField != "SBOMS_UPDATED_AT" || client.listVersionsInput.OrderByDir != "DESC" {
		t.Fatalf("expected latest lookup to request updated-at descending order: %#v", client.listVersionsInput)
	}
	output := toolResultMap(t, result)
	if output["versionId"] != "version-new" {
		t.Fatalf("unexpected output version: %#v", output)
	}
	if _, ok := output["resolvedVersion"].(map[string]interface{}); !ok {
		t.Fatalf("expected resolvedVersion in output: %#v", output)
	}
}

func TestHandleListVulnerabilities_FiltersByComponentIDAndPurl(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleListVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version_id":   "version-1",
				"component_id": "component-1",
				"purl":         "pkg:generic/openssl@3.0.0",
				"limit":        10,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListVulnerabilities returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListVulnerabilities returned tool error: %#v", result.Content)
	}
	if len(client.listVersionVulnsInput.ComponentIDs) != 0 || client.listVersionVulnsInput.Purl != "" {
		t.Fatalf("component identity filters should be applied before MCP pagination, not forwarded to ListVersionVulns: %#v", client.listVersionVulnsInput)
	}

	output := toolResultMap(t, result)
	vulns := output["vulnerabilities"].([]interface{})
	if len(vulns) != 1 {
		t.Fatalf("len(vulnerabilities) = %d, want 1", len(vulns))
	}
	component := vulns[0].(map[string]interface{})["component"].(map[string]interface{})
	if component["id"] != "component-1" || component["purl"] != "pkg:generic/openssl@3.0.0" {
		t.Fatalf("unexpected component output: %#v", component)
	}
	fixedVersions := vulns[0].(map[string]interface{})["fixedVersions"].([]interface{})
	if len(fixedVersions) != 1 || fixedVersions[0] != "3.0.1" {
		t.Fatalf("unexpected fixedVersions output: %#v", fixedVersions)
	}
	if component["sbomId"] != "version-1" || component["versionId"] != "version-1" {
		t.Fatalf("missing stable version identifiers: %#v", component)
	}
}

func TestHandleListVulnerabilities_ComponentFilterScansBeforePaging(t *testing.T) {
	mkVuln := func(id, componentID, purl string) api.ComponentVuln {
		return api.ComponentVuln{
			ID:          id,
			ComponentID: componentID,
			VersionID:   "version-1",
			Component: &api.VersionComponent{
				ID:        componentID,
				Name:      "axios",
				Version:   "0.30.0",
				Purl:      purl,
				VersionID: "version-1",
			},
			Vuln: &api.Vuln{
				ID:       "vuln-" + id,
				VulnID:   "CVE-2026-" + id,
				Severity: "high",
			},
		}
	}
	client := &fakeLynkClient{
		listVersionVulnsResults: map[string]*api.ComponentVulnsResult{
			"": {
				ComponentVulns: []api.ComponentVuln{
					mkVuln("1", "component-1", "pkg:generic/openssl@3.0.0"),
				},
				TotalCount:  4,
				HasNextPage: true,
				EndCursor:   "api-cursor-1",
			},
			"api-cursor-1": {
				ComponentVulns: []api.ComponentVuln{
					mkVuln("2", "axios-component", "pkg:npm/axios@0.30.0"),
					mkVuln("3", "axios-component", "pkg:npm/axios@0.30.0"),
					mkVuln("4", "axios-component", "pkg:npm/axios@0.30.0"),
				},
				TotalCount:  4,
				HasNextPage: false,
				EndCursor:   "api-cursor-2",
			},
		},
	}
	server := &Server{client: client}

	result, err := server.handleListVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version_id": "version-1",
				"purl":       "pkg:npm/axios@0.30.0",
				"limit":      2,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListVulnerabilities returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListVulnerabilities returned tool error: %#v", result.Content)
	}

	output := toolResultMap(t, result)
	vulns := output["vulnerabilities"].([]interface{})
	if len(vulns) != 2 {
		t.Fatalf("len(vulnerabilities) = %d, want 2", len(vulns))
	}
	if output["totalCount"] != float64(3) || output["hasMore"] != true || output["endCursor"] != "Mg==" {
		t.Fatalf("unexpected filtered pagination output: %#v", output)
	}
	if len(client.listVersionVulnsInputs) != 2 {
		t.Fatalf("ListVersionVulns call count = %d, want 2", len(client.listVersionVulnsInputs))
	}
	firstInput := client.listVersionVulnsInputs[0]
	if firstInput.First != 100 || firstInput.After != "" || firstInput.Purl != "" || len(firstInput.ComponentIDs) != 0 {
		t.Fatalf("unexpected first ListVersionVulns input: %#v", firstInput)
	}
	secondInput := client.listVersionVulnsInputs[1]
	if secondInput.After != "api-cursor-1" {
		t.Fatalf("second page After = %#v, want api-cursor-1", secondInput.After)
	}

	client.listVersionVulnsInputs = nil
	result, err = server.handleListVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version_id": "version-1",
				"purl":       "pkg:npm/axios@0.30.0",
				"limit":      2,
				"after":      "Mg==",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListVulnerabilities returned error on second filtered page: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleListVulnerabilities returned tool error on second filtered page: %#v", result.Content)
	}

	output = toolResultMap(t, result)
	vulns = output["vulnerabilities"].([]interface{})
	if len(vulns) != 1 {
		t.Fatalf("len(second page vulnerabilities) = %d, want 1", len(vulns))
	}
	if output["totalCount"] != float64(3) || output["hasMore"] != false || output["endCursor"] != "Mw==" {
		t.Fatalf("unexpected second filtered pagination output: %#v", output)
	}
}

func TestHandleListVulnerabilities_RejectsEmptyPurl(t *testing.T) {
	server := &Server{client: &fakeLynkClient{}}
	result, err := server.handleListVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"version_id": "version-1",
				"purl":       "",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleListVulnerabilities returned error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected empty purl to return a tool error")
	}
}

func TestHandleSearchVulnerabilities_FiltersComponentIDsAndPaginates(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleSearchVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"component_ids": []interface{}{"component-1", "component-2"},
				"component_id":  "component-3",
				"after":         "search-cursor-1",
				"limit":         25,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleSearchVulnerabilities returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleSearchVulnerabilities returned tool error: %#v", result.Content)
	}
	wantComponentIDs := []string{"component-1", "component-2", "component-3"}
	if !reflect.DeepEqual(client.listComponentVulnsInput.ComponentIDs, wantComponentIDs) {
		t.Fatalf("ComponentIDs = %#v, want %#v", client.listComponentVulnsInput.ComponentIDs, wantComponentIDs)
	}
	if client.listComponentVulnsInput.After != "search-cursor-1" || client.listComponentVulnsInput.First != 25 {
		t.Fatalf("unexpected pagination input: %#v", client.listComponentVulnsInput)
	}

	output := toolResultMap(t, result)
	if output["endCursor"] != "search-cursor-2" || output["hasMore"] != true {
		t.Fatalf("unexpected pagination output: %#v", output)
	}
}

func TestHandleSearchVulnerabilities_RejectsEmptyPurl(t *testing.T) {
	server := &Server{client: &fakeLynkClient{}}
	result, err := server.handleSearchVulnerabilities(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"purl": "",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleSearchVulnerabilities returned error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected empty purl to return a tool error")
	}
}

func TestHandleBulkUpdateComponentVex_RequiresConfirm(t *testing.T) {
	server := &Server{client: &fakeLynkClient{}}
	result, err := server.handleBulkUpdateComponentVex(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"component_vuln_ids": []interface{}{"component-vuln-1"},
				"vex_status_id":      "status-1",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleBulkUpdateComponentVex returned error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected missing confirm to return a tool error")
	}
}

func TestHandleBulkUpdateComponentVex_RejectsEmptyIDs(t *testing.T) {
	server := &Server{client: &fakeLynkClient{}}
	result, err := server.handleBulkUpdateComponentVex(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"confirm":            true,
				"component_vuln_ids": []interface{}{" "},
				"vex_status_id":      "status-1",
			},
		},
	})
	if err != nil {
		t.Fatalf("handleBulkUpdateComponentVex returned error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected empty component_vuln_ids to return a tool error")
	}
}

func TestHandleBulkUpdateComponentVex_ResolvesNamesAndReturnsPartialFailures(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleBulkUpdateComponentVex(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"confirm":            true,
				"component_vuln_ids": []interface{}{"component-vuln-1", "component-vuln-2"},
				"current_version_id": "version-1",
				"vex_status":         "not affected",
				"vex_justification":  "vulnerable code not present",
				"fixed_in":           "1.2.3",
				"propagate_vex":      false,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleBulkUpdateComponentVex returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleBulkUpdateComponentVex returned tool error: %#v", result.Content)
	}

	input := client.bulkUpdateComponentVexInput
	if !reflect.DeepEqual(input.ComponentVulnIDs, []string{"component-vuln-1", "component-vuln-2"}) {
		t.Fatalf("ComponentVulnIDs = %#v, want requested IDs", input.ComponentVulnIDs)
	}
	if input.CurrentVersionID == nil || *input.CurrentVersionID != "version-1" {
		t.Fatalf("CurrentVersionID = %#v, want version-1", input.CurrentVersionID)
	}
	if input.VexStatusID == nil || *input.VexStatusID != "status-1" {
		t.Fatalf("VexStatusID = %#v, want status-1", input.VexStatusID)
	}
	if input.VexJustificationID == nil || *input.VexJustificationID != "justification-1" {
		t.Fatalf("VexJustificationID = %#v, want justification-1", input.VexJustificationID)
	}
	if input.PropagateVex == nil || *input.PropagateVex {
		t.Fatalf("PropagateVex = %#v, want false", input.PropagateVex)
	}

	output := toolResultMap(t, result)
	if output["requestedCount"] != float64(2) || output["updatedCount"] != float64(1) || output["failedCount"] != float64(1) {
		t.Fatalf("unexpected counts: %#v", output)
	}
	updated := output["updated"].([]interface{})
	if len(updated) != 1 || updated[0].(map[string]interface{})["id"] != "component-vuln-1" {
		t.Fatalf("unexpected updated entries: %#v", updated)
	}
	failed := output["failed"].([]interface{})
	if len(failed) != 1 || failed[0].(map[string]interface{})["componentVulnId"] != "component-vuln-2" {
		t.Fatalf("unexpected failed entries: %#v", failed)
	}
}

func TestHandleGetTicketingStatus_PassesCursorsAndIncludeCreatedTickets(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleGetTicketingStatus(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"product_id":              "product-1",
				"products_limit":          25,
				"products_after":          "products-cursor-1",
				"policies_limit":          10,
				"policies_after":          "policies-cursor-1",
				"ticket_links_limit":      5,
				"ticket_links_after":      "tickets-cursor-1",
				"include_created_tickets": false,
			},
		},
	})
	if err != nil {
		t.Fatalf("handleGetTicketingStatus returned error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handleGetTicketingStatus returned tool error: %#v", result.Content)
	}

	input := client.ticketingStatusInput
	if input.ProductID != "product-1" {
		t.Fatalf("ProductID = %#v, want product-1", input.ProductID)
	}
	if input.ProductsFirst != 25 || input.ProductsAfter != "products-cursor-1" {
		t.Fatalf("unexpected products input: %#v", input)
	}
	if input.PoliciesFirst != 10 || input.PoliciesAfter != "policies-cursor-1" {
		t.Fatalf("unexpected policies input: %#v", input)
	}
	if input.TicketsFirst != 5 || input.TicketsAfter != "tickets-cursor-1" {
		t.Fatalf("unexpected tickets input: %#v", input)
	}
	if input.IncludeCreatedTickets == nil || *input.IncludeCreatedTickets {
		t.Fatalf("IncludeCreatedTickets = %#v, want false", input.IncludeCreatedTickets)
	}

	output := toolResultMap(t, result)
	for key, want := range map[string]interface{}{
		"productsEndCursor": "products-cursor-2",
		"policiesEndCursor": "policies-cursor-2",
		"ticketsEndCursor":  "tickets-cursor-2",
	} {
		if output[key] != want {
			t.Fatalf("%s = %#v, want %#v", key, output[key], want)
		}
	}
}

func TestSameVexName_NormalizesCommonInputForms(t *testing.T) {
	tests := []struct {
		left  string
		right string
	}{
		{left: "not_affected", right: "Not Affected"},
		{left: "vulnerable_code_not_present", right: "vulnerable-code-not-present"},
		{left: " fixed ", right: "fixed"},
	}

	for _, tt := range tests {
		if !sameVexName(tt.left, tt.right) {
			t.Fatalf("sameVexName(%q, %q) = false, want true", tt.left, tt.right)
		}
	}
}

func toolResultMap(t *testing.T, result *mcpg.CallToolResult) map[string]interface{} {
	t.Helper()
	if len(result.Content) != 1 {
		t.Fatalf("expected one content item, got %d", len(result.Content))
	}
	text, ok := result.Content[0].(mcpg.TextContent)
	if !ok {
		t.Fatalf("content type = %T, want TextContent", result.Content[0])
	}
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(text.Text), &output); err != nil {
		t.Fatalf("failed to decode tool result: %v", err)
	}
	return output
}

func TestGetVulnerabilityMetadataFilters_ExceptionalShortcut(t *testing.T) {
	filters := getVulnerabilityMetadataFilters(map[string]interface{}{
		"exceptional": true,
	})

	if !filters.MatchAny || !filters.Exceptional {
		t.Fatalf("shortcut flags = MatchAny:%t Exceptional:%t, want both true", filters.MatchAny, filters.Exceptional)
	}
	if filters.Kev == nil || !*filters.Kev {
		t.Fatalf("Kev = %#v, want true", filters.Kev)
	}
	if filters.EpssMin == nil || *filters.EpssMin != 0.05 {
		t.Fatalf("EpssMin = %#v, want 0.05", filters.EpssMin)
	}
	if filters.CvssMin == nil || *filters.CvssMin != 9.0 {
		t.Fatalf("CvssMin = %#v, want 9.0", filters.CvssMin)
	}
}

func TestComponentVulnMatchReasons_ReportsMatchingExceptionalMetadata(t *testing.T) {
	epssMin := 0.05
	cvssMin := 9.0
	kev := true
	filters := vulnerabilityMetadataFilters{
		EpssMin: &epssMin,
		CvssMin: &cvssMin,
		Kev:     &kev,
	}
	vuln := api.ComponentVuln{
		ID: "component-vuln-1",
		Vuln: &api.Vuln{
			CvssScore: 9.8,
			VulnInfo: &api.VulnInfo{
				EpssScore: 0.07,
				Kev:       true,
			},
		},
	}

	reasons := componentVulnMatchReasons(vuln, filters)

	if !reflect.DeepEqual(reasons, []string{"kev", "epss", "cvss"}) {
		t.Fatalf("reasons = %#v, want kev, epss, cvss", reasons)
	}
}

func TestFormatComponentVulns_IncludesCustomFieldAttributes(t *testing.T) {
	maxValue := 14
	vulns := []api.ComponentVuln{
		{
			ID:        "component-vuln-1",
			VersionID: "version-1",
			CustomFields: []api.ComponentVulnCustomField{
				{
					ID:                                   "custom-field-1",
					ComponentVulnCustomFieldDefinitionID: "field-def-1",
					Value:                                "12",
					VexableID:                            "component-vuln-1",
					VexableType:                          "ComponentVuln",
					ComponentVulnCustomFieldDefinition: &api.ComponentVulnCustomFieldDefinition{
						ID:           "field-def-1",
						DisplayName:  "CRM age",
						FieldType:    "RANGE",
						InternalName: "crm_age",
						MaxValue:     &maxValue,
					},
				},
			},
		},
	}

	formatted := formatComponentVulns(vulns, nil, false)
	customFields, ok := formatted[0]["customFieldAttributes"].([]map[string]interface{})
	if !ok {
		t.Fatalf("customFieldAttributes = %T %#v, want []map[string]interface{}", formatted[0]["customFieldAttributes"], formatted[0]["customFieldAttributes"])
	}
	if len(customFields) != 1 {
		t.Fatalf("len(customFieldAttributes) = %d, want 1", len(customFields))
	}
	field := customFields[0]
	if field["componentVulnCustomFieldDefinitionId"] != "field-def-1" || field["value"] != "12" {
		t.Fatalf("unexpected custom field output: %#v", field)
	}
	definition := field["definition"].(map[string]interface{})
	if definition["displayName"] != "CRM age" || definition["maxValue"] != 14 {
		t.Fatalf("unexpected definition output: %#v", definition)
	}
}

func TestFilterComponentVulnsByCvss_AppliesThresholds(t *testing.T) {
	cvssMin := 9.0
	filters := vulnerabilityMetadataFilters{CvssMin: &cvssMin}
	vulns := []api.ComponentVuln{
		{ID: "low", Vuln: &api.Vuln{CvssScore: 8.9}},
		{ID: "high", Vuln: &api.Vuln{CvssScore: 9.0}},
	}

	filtered := filterComponentVulnsByCvss(vulns, filters)

	if len(filtered) != 1 || filtered[0].ID != "high" {
		t.Fatalf("filtered = %#v, want only high", filtered)
	}
}
