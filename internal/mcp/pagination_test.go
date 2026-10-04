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
	"testing"

	"github.com/interlynk-io/lynk-mcp/internal/api"
	mcpg "github.com/mark3labs/mcp-go/mcp"
)

type paginationClient struct {
	lynkClient
	first int
	after string
}

func (f *paginationClient) ListComponents(ctx context.Context, input api.ListComponentsInput) (*api.ComponentsResult, error) {
	f.first, f.after = input.First, input.After
	return &api.ComponentsResult{TotalCount: 5, HasNextPage: true, EndCursor: "next"}, nil
}
func (f *paginationClient) ListPolicies(ctx context.Context, input api.ListPoliciesInput) (*api.PoliciesResult, error) {
	f.first, f.after = input.First, input.After
	return &api.PoliciesResult{TotalCount: 5, HasNextPage: true, EndCursor: "next"}, nil
}
func (f *paginationClient) ListPolicyResults(ctx context.Context, input api.ListPolicyResultsInput) (*api.PolicyResultsResult, error) {
	f.first, f.after = input.First, input.After
	return &api.PolicyResultsResult{TotalCount: 5, HasNextPage: true, EndCursor: "next"}, nil
}
func (f *paginationClient) ListLicenses(ctx context.Context, input api.ListLicensesInput) (*api.LicensesResult, error) {
	f.first, f.after = input.First, input.After
	return &api.LicensesResult{TotalCount: 5, HasNextPage: true, EndCursor: "next"}, nil
}
func TestListTools_CursorPagination(t *testing.T) {
	for _, name := range []string{"components", "policies", "policy_violations", "licenses"} {
		t.Run(name, func(t *testing.T) {
			for _, limit := range []int{2, 1000} {
				client := &paginationClient{}
				server := &Server{client: client}
				handler := server.handleListComponents
				switch name {
				case "policies":
					handler = server.handleListPolicies
				case "policy_violations":
					handler = server.handleListPolicyViolations
				case "licenses":
					handler = server.handleListLicenses
				}
				result, err := handler(context.Background(), mcpg.CallToolRequest{Params: mcpg.CallToolParams{Arguments: map[string]interface{}{"version_id": "version-1", "after": "previous", "limit": float64(limit)}}})
				if err != nil || result.IsError {
					t.Fatalf("call failed: %v %#v", err, result)
				}
				if client.first != min(limit, 100) || client.after != "previous" {
					t.Fatalf("pagination input: first=%d after=%q", client.first, client.after)
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(result.Content[0].(mcpg.TextContent).Text), &payload); err != nil {
					t.Fatal(err)
				}
				if payload["endCursor"] != "next" || payload["hasMore"] != true || payload["totalCount"] != float64(5) {
					t.Fatalf("pagination output: %#v", payload)
				}
			}
		})
	}
}
