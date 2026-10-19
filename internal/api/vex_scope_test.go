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

package api

import (
	"context"
	"errors"
	"testing"
)

func TestResolveVexScope_ClassifiesRows(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbomsCompVuln": [
					{"id": "direct-vuln", "isPart": false, "sbomId": "parent-version"},
					{"id": "part-vuln", "isPart": true, "sbomId": "child-version"}
				]
			}`,
		},
	}
	client := &Client{gql: gql}

	scopes, err := client.ResolveVexScope(context.Background(), "parent-version", []string{"direct-vuln", "part-vuln", "missing-vuln"})
	if err != nil {
		t.Fatalf("ResolveVexScope returned error: %v", err)
	}

	request := gql.requests[0]
	if ids := request["sbomIds"].([]string); len(ids) != 1 || ids[0] != "parent-version" {
		t.Fatalf("sbomIds = %#v, want [parent-version]", request["sbomIds"])
	}
	if ids := request["componentVulnIds"].([]string); len(ids) != 3 {
		t.Fatalf("componentVulnIds = %#v, want three IDs", request["componentVulnIds"])
	}

	if got := scopes["direct-vuln"]; got.Scope != VexScopeDirect || got.OwnerVersionID != "parent-version" {
		t.Fatalf("direct-vuln = %#v, want direct owned by parent-version", got)
	}
	if got := scopes["part-vuln"]; got.Scope != VexScopeParent || got.OwnerVersionID != "child-version" {
		t.Fatalf("part-vuln = %#v, want parent owned by child-version", got)
	}
	if got := scopes["missing-vuln"]; got.Scope != VexScopeNotVisible || got.OwnerVersionID != "" {
		t.Fatalf("missing-vuln = %#v, want not visible", got)
	}
}

func TestResolveVexScope_DuplicateRowsCollapseToOneScope(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbomsCompVuln": [
					{"id": "vuln-1", "isPart": true, "sbomId": "child-version"},
					{"id": "vuln-1", "isPart": true, "sbomId": "child-version"},
					{"id": "vuln-2", "isPart": true, "sbomId": "child-version"},
					{"id": "vuln-2", "isPart": false, "sbomId": "parent-version"}
				]
			}`,
		},
	}
	client := &Client{gql: gql}

	scopes, err := client.ResolveVexScope(context.Background(), "parent-version", []string{"vuln-1", "vuln-2"})
	if err != nil {
		t.Fatalf("ResolveVexScope returned error: %v", err)
	}
	if len(scopes) != 2 {
		t.Fatalf("len(scopes) = %d, want 2", len(scopes))
	}
	if scopes["vuln-1"].Scope != VexScopeParent {
		t.Fatalf("vuln-1 = %#v, want parent", scopes["vuln-1"])
	}
	if got := scopes["vuln-2"]; got.Scope != VexScopeDirect || got.OwnerVersionID != "parent-version" {
		t.Fatalf("vuln-2 = %#v, want direct (a direct row wins over a part row)", got)
	}
}

func TestResolveVexScope_ReturnsErrorAndNoScopes(t *testing.T) {
	client := &Client{gql: errGraphQLExecutor{err: errors.New("boom")}}

	scopes, err := client.ResolveVexScope(context.Background(), "version-1", []string{"vuln-1"})
	if err == nil || err.Error() != "boom" {
		t.Fatalf("err = %v, want boom", err)
	}
	if scopes != nil {
		t.Fatalf("scopes = %#v, want nil on error", scopes)
	}
}

func TestResolveVexScope_NoIDsSkipsTheCall(t *testing.T) {
	gql := &fakeGraphQLExecutor{}
	client := &Client{gql: gql}

	scopes, err := client.ResolveVexScope(context.Background(), "version-1", nil)
	if err != nil {
		t.Fatalf("ResolveVexScope returned error: %v", err)
	}
	if len(scopes) != 0 || len(gql.requests) != 0 {
		t.Fatalf("expected no scopes and no request, got %#v / %d requests", scopes, len(gql.requests))
	}
}
