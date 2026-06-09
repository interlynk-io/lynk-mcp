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
	listProductsInput api.ListProductsInput
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
