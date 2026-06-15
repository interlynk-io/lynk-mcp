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
	"reflect"
	"testing"
	"time"
)

func TestGetVuln_ByCveIDAcceptsTimestampWithoutTimezone(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"cveLookup": {
					"vulnId": "CVE-2025-32463",
					"description": "sudo vulnerability",
					"severity": "high",
					"published": "2025-06-30T21:15:30.257",
					"lastModified": "2025-07-01T10:11:12.123",
					"cvssScore": 9.3,
					"cvssVector": "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:C/C:H/I:H/A:H",
					"cwes": ["CWE-863"],
					"advisories": ["https://example.com/advisory"]
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	vuln, err := client.GetVuln(context.Background(), "", "CVE-2025-32463")
	if err != nil {
		t.Fatalf("GetVuln returned error: %v", err)
	}

	wantPublished := time.Date(2025, 6, 30, 21, 15, 30, 257000000, time.UTC)
	if !vuln.PublishedAt.Equal(wantPublished) {
		t.Fatalf("PublishedAt = %s, want %s", vuln.PublishedAt, wantPublished)
	}
	if vuln.VulnID != "CVE-2025-32463" || vuln.Severity != "high" {
		t.Fatalf("unexpected vulnerability: %#v", vuln)
	}
}

func TestListVersionVulns_AcceptsNestedVulnTimestampsWithoutTimezone(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbom": {
					"vulns": {
						"nodes": [
							{
								"id": "component-vuln-1",
								"componentId": "component-1",
								"vulnId": "vuln-1",
								"sbomId": "version-1",
								"fixedIn": "",
								"fixedVersions": [],
								"detail": "",
								"impact": "",
								"actionStmt": "",
								"createdAt": "2026-04-16T23:00:09Z",
								"updatedAt": "2026-04-16T23:00:09Z",
								"component": null,
								"vuln": {
									"id": "vuln-1",
									"vulnId": "CVE-2025-32463",
									"desc": "sudo vulnerability",
									"sev": "high",
									"cvssScore": 9.3,
									"cvssVector": "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:C/C:H/I:H/A:H",
									"source": "nvd",
									"publishedAt": "2025-06-30T21:15:30.257",
									"lastModifiedAt": "2025-07-01T10:11:12.123",
									"vulnInfo": {
										"cveId": "CVE-2025-32463",
										"epssScore": 0.1,
										"epssPercentile": 0.2,
										"kev": false,
										"cwes": ["CWE-863"]
									}
								},
								"vexStatus": null,
								"vexJustification": null
							}
						],
						"totalCount": 1,
						"pageInfo": {
							"hasNextPage": false,
							"endCursor": ""
						}
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.ListVersionVulns(context.Background(), ListVersionVulnsInput{VersionID: "version-1"})
	if err != nil {
		t.Fatalf("ListVersionVulns returned error: %v", err)
	}

	if len(result.ComponentVulns) != 1 {
		t.Fatalf("len(ComponentVulns) = %d, want 1", len(result.ComponentVulns))
	}
	wantPublished := time.Date(2025, 6, 30, 21, 15, 30, 257000000, time.UTC)
	if !result.ComponentVulns[0].Vuln.PublishedAt.Equal(wantPublished) {
		t.Fatalf("PublishedAt = %s, want %s", result.ComponentVulns[0].Vuln.PublishedAt, wantPublished)
	}
}

func TestListVersionVulns_MapsCustomFieldAttributes(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbom": {
					"vulns": {
						"nodes": [
							{
								"id": "component-vuln-1",
								"componentId": "component-1",
								"vulnId": "vuln-1",
								"sbomId": "version-1",
								"fixedIn": "",
								"fixedVersions": [],
								"detail": "",
								"impact": "",
								"actionStmt": "",
								"createdAt": "2026-04-16T23:00:09Z",
								"updatedAt": "2026-04-16T23:00:09Z",
								"component": null,
								"vuln": null,
								"vexStatus": null,
								"vexJustification": null,
								"componentVulnCustomFields": [
									{
										"id": "custom-field-1",
										"componentVulnCustomFieldDefinitionId": "field-def-1",
										"value": "12",
										"vexableId": "component-vuln-1",
										"vexableType": "ComponentVuln",
										"createdAt": "2026-04-17T00:00:00Z",
										"updatedAt": "2026-04-18T00:00:00Z",
										"componentVulnCustomFieldDefinition": {
											"id": "field-def-1",
											"displayName": "CRM age",
											"fieldType": "RANGE",
											"internalName": "crm_age",
											"minValue": 0,
											"maxValue": 14,
											"organizationId": "org-1"
										}
									}
								]
							}
						],
						"totalCount": 1,
						"pageInfo": {
							"hasNextPage": false,
							"endCursor": ""
						}
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.ListVersionVulns(context.Background(), ListVersionVulnsInput{VersionID: "version-1"})
	if err != nil {
		t.Fatalf("ListVersionVulns returned error: %v", err)
	}

	if len(result.ComponentVulns) != 1 {
		t.Fatalf("len(ComponentVulns) = %d, want 1", len(result.ComponentVulns))
	}
	customFields := result.ComponentVulns[0].CustomFields
	if len(customFields) != 1 {
		t.Fatalf("len(CustomFields) = %d, want 1", len(customFields))
	}
	field := customFields[0]
	if field.ID != "custom-field-1" || field.ComponentVulnCustomFieldDefinitionID != "field-def-1" || field.Value != "12" {
		t.Fatalf("unexpected custom field: %#v", field)
	}
	if field.ComponentVulnCustomFieldDefinition == nil {
		t.Fatal("expected custom field definition")
	}
	definition := field.ComponentVulnCustomFieldDefinition
	if definition.DisplayName != "CRM age" || definition.FieldType != "RANGE" || definition.InternalName != "crm_age" {
		t.Fatalf("unexpected custom field definition: %#v", definition)
	}
	if definition.MaxValue == nil || *definition.MaxValue != 14 {
		t.Fatalf("MaxValue = %#v, want 14", definition.MaxValue)
	}
}

func TestListComponentVulns_PassesComponentIDs(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"componentVulns": {
					"nodes": [],
					"totalCount": 0,
					"pageInfo": {
						"hasNextPage": false,
						"endCursor": ""
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	_, err := client.ListComponentVulns(context.Background(), ListComponentVulnsInput{
		ComponentIDs: []string{"component-1", "component-2"},
		After:        "cursor-1",
		First:        25,
	})
	if err != nil {
		t.Fatalf("ListComponentVulns returned error: %v", err)
	}

	request := gql.requests[0]
	if request["after"] != "cursor-1" || request["first"] != 25 {
		t.Fatalf("unexpected pagination variables: %#v", request)
	}
	got, ok := request["componentIds"].([]string)
	if !ok {
		t.Fatalf("componentIds = %T %#v, want []string", request["componentIds"], request["componentIds"])
	}
	want := []string{"component-1", "component-2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("componentIds = %#v, want %#v", got, want)
	}
}

func TestListVersionVulns_SendsEpssRangeFilter(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbom": {
					"vulns": {
						"nodes": [],
						"totalCount": 0,
						"pageInfo": {
							"hasNextPage": false,
							"endCursor": ""
						}
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}
	epssMin := 0.05

	_, err := client.ListVersionVulns(context.Background(), ListVersionVulnsInput{
		VersionID: "version-1",
		EpssMin:   &epssMin,
	})
	if err != nil {
		t.Fatalf("ListVersionVulns returned error: %v", err)
	}

	epss, ok := gql.requests[0]["epss"].(map[string]interface{})
	if !ok {
		t.Fatalf("epss variable = %#v, want range map", gql.requests[0]["epss"])
	}
	if epss["min"] != 0.05 || epss["max"] != 1.0 {
		t.Fatalf("epss variable = %#v, want min 0.05 max 1.0", epss)
	}
}

func TestListComponentVulns_SendsMetadataAndScopeFilters(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"componentVulns": {
					"nodes": [],
					"totalCount": 0,
					"pageInfo": {
						"hasNextPage": false,
						"endCursor": ""
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}
	epssMax := 0.2
	kev := true

	_, err := client.ListComponentVulns(context.Background(), ListComponentVulnsInput{
		EpssMax:        &epssMax,
		Kev:            &kev,
		ProductIDs:     []string{"product-1"},
		EnvironmentIDs: []string{"environment-1"},
	})
	if err != nil {
		t.Fatalf("ListComponentVulns returned error: %v", err)
	}

	epss, ok := gql.requests[0]["epss"].(map[string]interface{})
	if !ok {
		t.Fatalf("epss variable = %#v, want range map", gql.requests[0]["epss"])
	}
	if epss["min"] != 0.0 || epss["max"] != 0.2 {
		t.Fatalf("epss variable = %#v, want min 0.0 max 0.2", epss)
	}
	if gql.requests[0]["kev"] != true {
		t.Fatalf("kev variable = %#v, want true", gql.requests[0]["kev"])
	}
	if got := gql.requests[0]["projectGroupIds"]; !stringSlicesEqual(got, []string{"product-1"}) {
		t.Fatalf("projectGroupIds = %#v, want product-1", got)
	}
	if got := gql.requests[0]["projectIds"]; !stringSlicesEqual(got, []string{"environment-1"}) {
		t.Fatalf("projectIds = %#v, want environment-1", got)
	}
}

func TestListComponentVulns_MapsCustomFieldAttributes(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"componentVulns": {
					"nodes": [
						{
							"id": "component-vuln-1",
							"componentId": "component-1",
							"vulnId": "vuln-1",
							"sbomId": "version-1",
							"fixedIn": "",
							"fixedVersions": [],
							"createdAt": "2026-04-16T23:00:09Z",
							"updatedAt": "2026-04-16T23:00:09Z",
							"component": null,
							"vuln": null,
							"vexStatus": null,
							"componentVulnCustomFields": [
								{
									"id": "custom-field-1",
									"componentVulnCustomFieldDefinitionId": "field-def-1",
									"value": "7",
									"vexableId": "component-vuln-1",
									"vexableType": "ComponentVuln",
									"createdAt": "2026-04-17T00:00:00Z",
									"updatedAt": "2026-04-18T00:00:00Z",
									"componentVulnCustomFieldDefinition": {
										"id": "field-def-1",
										"displayName": "CRM age",
										"fieldType": "RANGE",
										"internalName": "crm_age",
										"minValue": 0,
										"maxValue": 14,
										"organizationId": "org-1"
									}
								}
							]
						}
					],
					"totalCount": 1,
					"pageInfo": {
						"hasNextPage": false,
						"endCursor": ""
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.ListComponentVulns(context.Background(), ListComponentVulnsInput{})
	if err != nil {
		t.Fatalf("ListComponentVulns returned error: %v", err)
	}

	if len(result.ComponentVulns) != 1 {
		t.Fatalf("len(ComponentVulns) = %d, want 1", len(result.ComponentVulns))
	}
	customFields := result.ComponentVulns[0].CustomFields
	if len(customFields) != 1 {
		t.Fatalf("len(CustomFields) = %d, want 1", len(customFields))
	}
	if customFields[0].Value != "7" || customFields[0].ComponentVulnCustomFieldDefinition == nil {
		t.Fatalf("unexpected custom field: %#v", customFields[0])
	}
}

func TestGetVexStatuses(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"vexStatuses": [
					{"id": "status-1", "name": "affected"},
					{"id": "status-2", "name": "not_affected"}
				]
			}`,
		},
	}
	client := &Client{gql: gql}

	statuses, err := client.GetVexStatuses(context.Background())
	if err != nil {
		t.Fatalf("GetVexStatuses returned error: %v", err)
	}

	if len(statuses) != 2 {
		t.Fatalf("len(statuses) = %d, want 2", len(statuses))
	}
	if statuses[1].ID != "status-2" || statuses[1].Name != "not_affected" {
		t.Fatalf("unexpected statuses: %#v", statuses)
	}
}

func stringSlicesEqual(got interface{}, want []string) bool {
	gotStrings, ok := got.([]string)
	if !ok || len(gotStrings) != len(want) {
		return false
	}
	for i := range want {
		if gotStrings[i] != want[i] {
			return false
		}
	}
	return true
}

func TestGetVexJustifications(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"vexJustifications": [
					{"id": "justification-1", "name": "vulnerable_code_not_present"},
					{"id": "justification-2", "name": "component_not_present"}
				]
			}`,
		},
	}
	client := &Client{gql: gql}

	justifications, err := client.GetVexJustifications(context.Background())
	if err != nil {
		t.Fatalf("GetVexJustifications returned error: %v", err)
	}

	if len(justifications) != 2 {
		t.Fatalf("len(justifications) = %d, want 2", len(justifications))
	}
	if justifications[0].ID != "justification-1" || justifications[0].Name != "vulnerable_code_not_present" {
		t.Fatalf("unexpected justifications: %#v", justifications)
	}
}
