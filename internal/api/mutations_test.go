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
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestUpdateComponent_MapsInputsAndResults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"componentUpdate": {
					"component": {
						"id": "component-1",
						"name": "jest-diff",
						"version": "27.3.2",
						"kind": "library",
						"purl": "",
						"cpes": [],
						"licensesExp": "",
						"group": "npm",
						"description": "",
						"scope": "",
						"copyright": "",
						"primary": false,
						"internal": false,
						"uniqueId": "bom-ref-1",
						"sbomId": "version-1",
						"notice": "",
						"supportLevel": "",
						"endOfSupport": "",
						"checksums": [
							{"alg": "SHA-512", "content": "abc"}
						],
						"externalUrls": [
							{"name": "repo", "url": "https://example.com/repo"}
						]
					},
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	kind := "library"
	name := "jest-diff"
	empty := ""
	version := "27.3.2"
	primary := false
	internal := false
	generateUniqueID := true
	cpes := []string{}
	checksums := []ChecksumInput{{Alg: "SHA_512", Content: "abc"}}
	externalURLs := []ExternalURLInput{{Name: "repo", URL: "https://example.com/repo"}}

	result, err := client.UpdateComponent(context.Background(), UpdateComponentInput{
		ID:               "component-1",
		VersionID:        "version-1",
		Kind:             &kind,
		Name:             &name,
		Description:      &empty,
		Version:          &version,
		Licenses:         &LicenseInput{LicensesExp: ""},
		Cpes:             &cpes,
		Purl:             &empty,
		Primary:          &primary,
		Internal:         &internal,
		GenerateUniqueID: &generateUniqueID,
		Scope:            &empty,
		Checksums:        &checksums,
		ExternalURLs:     &externalURLs,
	})
	if err != nil {
		t.Fatalf("UpdateComponent returned error: %v", err)
	}

	request := gql.requests[0]
	if request["id"] != "component-1" || request["sbomId"] != "version-1" {
		t.Fatalf("unexpected required variables: %#v", request)
	}
	if request["primary"] != false || request["internal"] != false {
		t.Fatalf("false bools were not preserved: %#v", request)
	}
	if request["purl"] != "" || request["description"] != "" {
		t.Fatalf("empty strings were not preserved: %#v", request)
	}
	if got := request["cpes"].([]string); len(got) != 0 {
		t.Fatalf("cpes = %#v, want empty slice", got)
	}
	if request["generateUniqueId"] != true {
		t.Fatalf("generateUniqueId = %#v, want true", request["generateUniqueId"])
	}
	licenses := request["licenses"].(*LicenseInput)
	if licenses.LicensesExp != "" {
		t.Fatalf("licensesExp = %q, want empty string", licenses.LicensesExp)
	}

	if result.Component == nil {
		t.Fatal("expected updated component")
	}
	if result.Component.ID != "component-1" || result.Component.VersionID != "version-1" {
		t.Fatalf("unexpected component: %#v", result.Component)
	}
	if len(result.Component.Checksums) != 1 || result.Component.Checksums[0].Content != "abc" {
		t.Fatalf("unexpected checksums: %#v", result.Component.Checksums)
	}
}

func TestUpdateComponentSupplier_MapsInputsAndErrors(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"compSupplierUpdate": {
					"compSupplier": null,
					"errors": ["Component supplier not found"]
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	name := "Supplier"
	empty := ""
	result, err := client.UpdateComponentSupplier(context.Background(), UpdateComponentSupplierInput{
		ID:           "supplier-1",
		Name:         &name,
		URL:          &empty,
		ContactName:  &empty,
		ContactEmail: &empty,
	})
	if err != nil {
		t.Fatalf("UpdateComponentSupplier returned error: %v", err)
	}

	request := gql.requests[0]
	if request["id"] != "supplier-1" || request["name"] != "Supplier" {
		t.Fatalf("unexpected supplier variables: %#v", request)
	}
	if request["url"] != "" || request["contactName"] != "" || request["contactEmail"] != "" {
		t.Fatalf("empty supplier fields were not preserved: %#v", request)
	}
	if result.Supplier != nil {
		t.Fatalf("Supplier = %#v, want nil", result.Supplier)
	}
	if len(result.Errors) != 1 || result.Errors[0] != "Component supplier not found" {
		t.Fatalf("Errors = %#v, want API error", result.Errors)
	}
}

func TestUpdateComponentVex_MapsInputsAndResults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{"sbomsCompVuln": [{"id": "component-vuln-1", "isPart": false, "sbomId": "version-1"}]}`,
			`{
				"componentVexUpdate": {
					"componentVuln": {
						"id": "component-vuln-1",
						"componentId": "component-1",
						"vulnId": "vuln-1",
						"sbomId": "version-1",
						"fixedIn": "1.2.3",
						"detail": "detail",
						"impact": "impact",
						"actionStmt": "action",
						"vexStatus": {"id": "status-1", "name": "not_affected"},
						"vexJustification": {"id": "justification-1", "name": "vulnerable_code_not_present"}
					},
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	statusID := "status-1"
	justificationID := "justification-1"
	note := ""
	propagateVex := false
	resolutionDate := "2026-05-07"
	destroy := false
	customFields := []ComponentVulnCustomFieldAttributeInput{{
		ID:                                   "field-value-1",
		ComponentVulnCustomFieldDefinitionID: "field-def-1",
		Value:                                "field-value",
		Destroy:                              &destroy,
	}}

	result, err := client.UpdateComponentVex(context.Background(), UpdateComponentVexInput{
		ComponentVulnID:                    "component-vuln-1",
		CurrentVersionID:                   "version-1",
		VexStatusID:                        &statusID,
		VexJustificationID:                 &justificationID,
		Note:                               &note,
		PropagateVex:                       &propagateVex,
		ResolutionDate:                     &resolutionDate,
		ComponentVulnCustomFieldAttributes: &customFields,
	})
	if err != nil {
		t.Fatalf("UpdateComponentVex returned error: %v", err)
	}

	request := gql.requests[1]
	if request["componentVulnId"] != "component-vuln-1" || request["currentSbomId"] != "version-1" {
		t.Fatalf("unexpected VEX required variables: %#v", request)
	}
	if request["note"] != "" || request["propagateVex"] != false {
		t.Fatalf("destructive VEX values were not preserved: %#v", request)
	}
	if request["resolutionDate"] != "2026-05-07" {
		t.Fatalf("resolutionDate = %#v, want 2026-05-07", request["resolutionDate"])
	}
	if got := request["componentVulnCustomFieldAttributes"].([]ComponentVulnCustomFieldAttributeInput); len(got) != 1 || got[0].Destroy == nil || *got[0].Destroy {
		t.Fatalf("unexpected custom fields: %#v", got)
	}

	if result.ComponentVuln == nil {
		t.Fatal("expected updated component vuln")
	}
	if result.Scope != VexScopeDirect || result.ParentVersionID != "" || result.OwnerVersionID != "version-1" {
		t.Fatalf("unexpected scope fields: %#v", result)
	}
	if result.ComponentVuln.VexStatus == nil || result.ComponentVuln.VexStatus.ID != "status-1" {
		t.Fatalf("unexpected VEX result: %#v", result.ComponentVuln)
	}
}

func TestUpdateComponentVex_ParentVersionWritesParentDisposition(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{"sbomsCompVuln": [{"id": "component-vuln-1", "isPart": true, "sbomId": "owner-version"}]}`,
			`{
				"componentVexUpdate": {
					"componentVuln": {
						"id": "disposition-1",
						"componentId": "component-1",
						"vulnId": "vuln-1",
						"sbomId": "owner-version",
						"vexStatus": {"id": "status-1", "name": "In Triage"}
					},
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	statusID := "status-1"
	result, err := client.UpdateComponentVex(context.Background(), UpdateComponentVexInput{
		ComponentVulnID:  "component-vuln-1",
		CurrentVersionID: "parent-version",
		VexStatusID:      &statusID,
	})
	if err != nil {
		t.Fatalf("UpdateComponentVex returned error: %v", err)
	}

	if gql.requests[0]["sbomIds"].([]string)[0] != "parent-version" {
		t.Fatalf("scope check did not use the parent version: %#v", gql.requests[0])
	}
	if gql.requests[1]["currentSbomId"] != "parent-version" {
		t.Fatalf("mutation did not use the parent version: %#v", gql.requests[1])
	}
	if result.Scope != VexScopeParent || result.ParentVersionID != "parent-version" || result.OwnerVersionID != "owner-version" {
		t.Fatalf("unexpected scope fields: %#v", result)
	}
}

func TestUpdateComponentVex_NotVisibleVersionIsRefusedWithoutMutation(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{`{"sbomsCompVuln": []}`},
	}
	client := &Client{gql: gql}

	statusID := "status-1"
	result, err := client.UpdateComponentVex(context.Background(), UpdateComponentVexInput{
		ComponentVulnID:  "component-vuln-1",
		CurrentVersionID: "unrelated-version",
		VexStatusID:      &statusID,
	})
	if err == nil {
		t.Fatalf("expected a refusal, got result %#v", result)
	}
	var notVisible *VexNotVisibleError
	if !errors.As(err, &notVisible) {
		t.Fatalf("err = %T %v, want *VexNotVisibleError", err, err)
	}
	for _, want := range []string{"component-vuln-1", "unrelated-version", "retracted", "get_version_relations"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q does not mention %q", err.Error(), want)
		}
	}
	if len(gql.requests) != 1 {
		t.Fatalf("requests = %d, want only the scope check", len(gql.requests))
	}
}

func TestUpdateComponentVex_ScopeCheckFailureSendsNoMutation(t *testing.T) {
	gql := &failFirstGraphQLExecutor{err: errors.New("scope boom")}
	client := &Client{gql: gql}

	statusID := "status-1"
	_, err := client.UpdateComponentVex(context.Background(), UpdateComponentVexInput{
		ComponentVulnID:  "component-vuln-1",
		CurrentVersionID: "version-1",
		VexStatusID:      &statusID,
	})
	if err == nil || !strings.Contains(err.Error(), "scope boom") {
		t.Fatalf("err = %v, want the scope check error", err)
	}
	if gql.calls != 1 {
		t.Fatalf("calls = %d, want 1 (no mutation after a failed check)", gql.calls)
	}
}

// failFirstGraphQLExecutor fails the first call and counts every call.
type failFirstGraphQLExecutor struct {
	err   error
	calls int
}

func (f *failFirstGraphQLExecutor) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	f.calls++
	if f.calls == 1 {
		return f.err
	}
	return nil
}

func TestBulkUpdateComponentVex_MapsInputsAndResults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{"sbomsCompVuln": [
				{"id": "component-vuln-1", "isPart": false, "sbomId": "version-1"},
				{"id": "component-vuln-2", "isPart": false, "sbomId": "version-1"}
			]}`,
			`{
				"componentVexBulkUpdate": {
					"componentVulns": [
						{
							"id": "component-vuln-1",
							"componentId": "component-1",
							"vulnId": "vuln-1",
							"sbomId": "version-1",
							"fixedIn": "1.2.3",
							"detail": "detail",
							"impact": "impact",
							"actionStmt": "action",
							"vexStatus": {"id": "status-1", "name": "not_affected"},
							"vexJustification": {"id": "justification-1", "name": "vulnerable_code_not_present"}
						}
					],
					"errors": ["component-vuln-2 failed"]
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	statusID := "status-1"
	justificationID := "justification-1"
	currentVersionID := "version-1"
	note := ""
	propagateVex := false
	resolutionDate := "2026-05-07"

	result, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs:   []string{"component-vuln-1", "component-vuln-2"},
		CurrentVersionID:   &currentVersionID,
		VexStatusID:        &statusID,
		VexJustificationID: &justificationID,
		Note:               &note,
		PropagateVex:       &propagateVex,
		ResolutionDate:     &resolutionDate,
	})
	if err != nil {
		t.Fatalf("BulkUpdateComponentVex returned error: %v", err)
	}

	request := gql.requests[1]
	if got := request["componentVulnIds"].([]string); len(got) != 2 || got[0] != "component-vuln-1" || got[1] != "component-vuln-2" {
		t.Fatalf("componentVulnIds = %#v, want two requested IDs", request["componentVulnIds"])
	}
	if request["currentSbomId"] != "version-1" || request["vexStatusId"] != "status-1" {
		t.Fatalf("unexpected bulk VEX variables: %#v", request)
	}
	if request["note"] != "" || request["propagateVex"] != false {
		t.Fatalf("destructive bulk VEX values were not preserved: %#v", request)
	}

	if len(result.ComponentVulns) != 1 || result.ComponentVulns[0].ID != "component-vuln-1" {
		t.Fatalf("unexpected component vulns: %#v", result.ComponentVulns)
	}
	if len(result.Errors) != 1 || result.Errors[0] != "component-vuln-2 failed" {
		t.Fatalf("Errors = %#v, want partial failure", result.Errors)
	}
}

func bulkVexMutationPage(ids ...string) string {
	nodes := make([]string, len(ids))
	for i, id := range ids {
		nodes[i] = `{"id": "` + id + `", "componentId": "component-1", "vulnId": "vuln-1", "sbomId": "owner-version", "vexStatus": {"id": "status-1", "name": "In Triage"}}`
	}
	return `{"componentVexBulkUpdate": {"componentVulns": [` + strings.Join(nodes, ",") + `], "errors": []}}`
}

func TestBulkUpdateComponentVex_MixedDirectAndParentCarryPerItemScope(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{"sbomsCompVuln": [
				{"id": "direct-vuln", "isPart": false, "sbomId": "parent-version"},
				{"id": "part-vuln", "isPart": true, "sbomId": "child-version"}
			]}`,
			bulkVexMutationPage("direct-vuln", "part-vuln"),
		},
	}
	client := &Client{gql: gql}

	currentVersionID := "parent-version"
	statusID := "status-1"
	result, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs: []string{"direct-vuln", "part-vuln"},
		CurrentVersionID: &currentVersionID,
		VexStatusID:      &statusID,
	})
	if err != nil {
		t.Fatalf("BulkUpdateComponentVex returned error: %v", err)
	}

	if ids := gql.requests[1]["componentVulnIds"].([]string); len(ids) != 2 {
		t.Fatalf("mutation IDs = %#v, want both", gql.requests[1]["componentVulnIds"])
	}
	if got := result.Scopes["direct-vuln"]; got.Scope != VexScopeDirect {
		t.Fatalf("direct-vuln scope = %#v, want direct", got)
	}
	if got := result.Scopes["part-vuln"]; got.Scope != VexScopeParent || got.OwnerVersionID != "child-version" {
		t.Fatalf("part-vuln scope = %#v, want parent owned by child-version", got)
	}
	if result.ParentVersionID != "parent-version" {
		t.Fatalf("ParentVersionID = %q, want parent-version", result.ParentVersionID)
	}
	if len(result.Refused) != 0 {
		t.Fatalf("Refused = %#v, want none", result.Refused)
	}
}

func TestBulkUpdateComponentVex_NotVisibleItemIsRefusedAndTheRestAreSent(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{"sbomsCompVuln": [{"id": "good-vuln", "isPart": true, "sbomId": "child-version"}]}`,
			bulkVexMutationPage("good-vuln"),
		},
	}
	client := &Client{gql: gql}

	currentVersionID := "parent-version"
	statusID := "status-1"
	result, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs: []string{"good-vuln", "bad-vuln"},
		CurrentVersionID: &currentVersionID,
		VexStatusID:      &statusID,
	})
	if err != nil {
		t.Fatalf("BulkUpdateComponentVex returned error: %v", err)
	}

	ids := gql.requests[1]["componentVulnIds"].([]string)
	if len(ids) != 1 || ids[0] != "good-vuln" {
		t.Fatalf("mutation IDs = %#v, want only good-vuln", ids)
	}
	if len(result.Refused) != 1 || result.Refused[0].ComponentVulnID != "bad-vuln" || !strings.Contains(result.Refused[0].Reason, "get_version_relations") {
		t.Fatalf("Refused = %#v, want bad-vuln with the not-visible reason", result.Refused)
	}
	if len(result.ComponentVulns) != 1 {
		t.Fatalf("ComponentVulns = %#v, want one", result.ComponentVulns)
	}
}

func TestBulkUpdateComponentVex_AllNotVisibleSendsNoMutation(t *testing.T) {
	gql := &fakeGraphQLExecutor{pages: []string{`{"sbomsCompVuln": []}`}}
	client := &Client{gql: gql}

	currentVersionID := "unrelated-version"
	statusID := "status-1"
	result, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs: []string{"a", "b"},
		CurrentVersionID: &currentVersionID,
		VexStatusID:      &statusID,
	})
	if err != nil {
		t.Fatalf("BulkUpdateComponentVex returned error: %v", err)
	}
	if len(gql.requests) != 1 {
		t.Fatalf("requests = %d, want only the scope check", len(gql.requests))
	}
	if len(result.Refused) != 2 || len(result.ComponentVulns) != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestBulkUpdateComponentVex_WithoutCurrentVersionSkipsTheScopeCheck(t *testing.T) {
	gql := &fakeGraphQLExecutor{pages: []string{bulkVexMutationPage("vuln-1")}}
	client := &Client{gql: gql}

	statusID := "status-1"
	result, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs: []string{"vuln-1"},
		VexStatusID:      &statusID,
	})
	if err != nil {
		t.Fatalf("BulkUpdateComponentVex returned error: %v", err)
	}
	if len(gql.requests) != 1 {
		t.Fatalf("requests = %d, want only the mutation", len(gql.requests))
	}
	if result.ParentVersionID != "" || len(result.Scopes) != 0 {
		t.Fatalf("expected no scope data, got %#v", result)
	}
}

func TestBulkUpdateComponentVex_ScopeCheckFailureSendsNoMutation(t *testing.T) {
	gql := &failFirstGraphQLExecutor{err: errors.New("scope boom")}
	client := &Client{gql: gql}

	currentVersionID := "version-1"
	statusID := "status-1"
	_, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs: []string{"vuln-1"},
		CurrentVersionID: &currentVersionID,
		VexStatusID:      &statusID,
	})
	if err == nil || !strings.Contains(err.Error(), "scope boom") {
		t.Fatalf("err = %v, want the scope check error", err)
	}
	if gql.calls != 1 {
		t.Fatalf("calls = %d, want 1", gql.calls)
	}
}

func TestBulkUpdateComponentVex_MutationErrorNamesTheSentIDs(t *testing.T) {
	gql := &secondCallFailsGraphQLExecutor{
		first: `{"sbomsCompVuln": [{"id": "vuln-1", "isPart": false, "sbomId": "version-1"}, {"id": "vuln-2", "isPart": false, "sbomId": "version-1"}]}`,
		err:   errors.New("api boom"),
	}
	client := &Client{gql: gql}

	currentVersionID := "version-1"
	statusID := "status-1"
	_, err := client.BulkUpdateComponentVex(context.Background(), BulkUpdateComponentVexInput{
		ComponentVulnIDs: []string{"vuln-1", "vuln-2"},
		CurrentVersionID: &currentVersionID,
		VexStatusID:      &statusID,
	})
	if err == nil || !strings.Contains(err.Error(), "api boom") || !strings.Contains(err.Error(), "vuln-1") || !strings.Contains(err.Error(), "vuln-2") {
		t.Fatalf("err = %v, want the API error naming both sent IDs", err)
	}
}

// secondCallFailsGraphQLExecutor answers the first call with a fixed page and fails the next.
type secondCallFailsGraphQLExecutor struct {
	first string
	err   error
	calls int
}

func (f *secondCallFailsGraphQLExecutor) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	f.calls++
	if f.calls == 1 {
		return json.Unmarshal([]byte(f.first), result)
	}
	return f.err
}
