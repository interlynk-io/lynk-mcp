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

	"github.com/interlynk-io/lynk-mcp/internal/api"
	mcpg "github.com/mark3labs/mcp-go/mcp"
)

type fakeLynkClient struct {
	lynkClient
	listProductsInput       api.ListProductsInput
	listVersionVulnsInput   api.ListVersionVulnsInput
	listComponentVulnsInput api.ListComponentVulnsInput
	ticketingStatusInput    api.TicketingStatusInput
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
				Name:          "default",
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

func (f *fakeLynkClient) ListVersionVulns(ctx context.Context, input api.ListVersionVulnsInput) (*api.ComponentVulnsResult, error) {
	f.listVersionVulnsInput = input
	return &api.ComponentVulnsResult{
		ComponentVulns: []api.ComponentVuln{
			{
				ID:          "component-vuln-1",
				ComponentID: "component-1",
				VersionID:   input.VersionID,
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

func TestHandleListProducts_PassesAfterAndReturnsEndCursor(t *testing.T) {
	client := &fakeLynkClient{}
	server := &Server{client: client}
	result, err := server.handleListProducts(context.Background(), mcpg.CallToolRequest{
		Params: mcpg.CallToolParams{
			Arguments: map[string]interface{}{
				"limit":  2,
				"after":  "cursor-1",
				"search": "Product",
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

	output := toolResultMap(t, result)
	if output["endCursor"] != "cursor-2" {
		t.Fatalf("endCursor = %#v, want cursor-2", output["endCursor"])
	}
	if output["hasMore"] != true {
		t.Fatalf("hasMore = %#v, want true", output["hasMore"])
	}
	products := output["products"].([]interface{})
	product := products[0].(map[string]interface{})
	repository := product["repository"].(map[string]interface{})
	if repository["importStatus"] != "completed" || repository["fullName"] != "org/repo" {
		t.Fatalf("unexpected repository output: %#v", repository)
	}
	summary := product["ticketingDefaultsSummary"].(map[string]interface{})
	if summary["jiraConfigured"] != true {
		t.Fatalf("unexpected ticketing summary output: %#v", summary)
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
	if !reflect.DeepEqual(client.listVersionVulnsInput.ComponentIDs, []string{"component-1"}) {
		t.Fatalf("ComponentIDs = %#v, want component-1", client.listVersionVulnsInput.ComponentIDs)
	}
	if client.listVersionVulnsInput.Purl != "pkg:generic/openssl@3.0.0" {
		t.Fatalf("Purl = %#v, want openssl purl", client.listVersionVulnsInput.Purl)
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
	if component["sbomId"] != "version-1" || component["versionId"] != "version-1" {
		t.Fatalf("missing stable version identifiers: %#v", component)
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
