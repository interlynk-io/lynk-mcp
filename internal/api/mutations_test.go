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

	request := gql.requests[0]
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
	if result.ComponentVuln.VexStatus == nil || result.ComponentVuln.VexStatus.ID != "status-1" {
		t.Fatalf("unexpected VEX result: %#v", result.ComponentVuln)
	}
}
