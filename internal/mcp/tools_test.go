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
	listProductsInput     api.ListProductsInput
	listVersionVulnsInput api.ListVersionVulnsInput
	ticketingStatusInput  api.TicketingStatusInput
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
			},
		},
		TotalCount:  3,
		HasNextPage: true,
		EndCursor:   "cursor-2",
	}, nil
}

func (f *fakeLynkClient) ListVersionVulns(ctx context.Context, input api.ListVersionVulnsInput) (*api.ComponentVulnsResult, error) {
	f.listVersionVulnsInput = input
	return &api.ComponentVulnsResult{
		ComponentVulns: []api.ComponentVuln{
			{
				ID:        "component-vuln-1",
				VersionID: input.VersionID,
				Component: &api.VersionComponent{
					ID:      "component-1",
					Name:    "openssl",
					Version: "3.0.0",
				},
				Vuln: &api.Vuln{
					ID:       "vuln-1",
					VulnID:   "CVE-2026-0001",
					Severity: "high",
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
