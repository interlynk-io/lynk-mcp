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

func TestSearchVersions_MapsProductEnvironmentAndPagination(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"projectVersions": {
					"nodes": [
						{
							"id": "version-1",
							"projectVersion": "1.0.0",
							"spec": "CycloneDX",
							"specVersion": "1.6",
							"format": "json",
							"lifecycle": "released",
							"createdAt": "2026-06-01T00:00:00Z",
							"updatedAt": "2026-06-02T00:00:00Z",
							"projectId": "env-1",
							"stats": {
								"compCount": 2,
								"compPurlCount": 2,
								"compCpeCount": 1,
								"compLicenseCount": 2,
								"compSupplierCount": 1,
								"vulnStats": {"critical": 1}
							},
							"project": {
								"id": "env-1",
								"name": "production",
								"projectGroupId": "product-1",
								"projectGroup": {
									"id": "product-1",
									"name": "Product 1"
								}
							}
						}
					],
					"totalCount": 1,
					"pageInfo": {
						"hasNextPage": false,
						"endCursor": "cursor-1"
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.SearchVersions(context.Background(), VersionSearchInput{
		Search: "1.0.0",
		First:  25,
		After:  "cursor-0",
	})
	if err != nil {
		t.Fatalf("SearchVersions returned error: %v", err)
	}

	request := gql.requests[0]
	if request["search"] != "1.0.0" || request["first"] != 25 || request["after"] != "cursor-0" {
		t.Fatalf("unexpected request variables: %#v", request)
	}
	if len(result.Versions) != 1 {
		t.Fatalf("len(Versions) = %d, want 1", len(result.Versions))
	}
	version := result.Versions[0]
	if version.Environment == nil || version.Environment.Name != "production" {
		t.Fatalf("unexpected environment: %#v", version.Environment)
	}
	if version.Environment.Product == nil || version.Environment.Product.Name != "Product 1" {
		t.Fatalf("unexpected product: %#v", version.Environment.Product)
	}
	if version.Stats == nil || version.Stats.VulnStats["critical"] != float64(1) {
		t.Fatalf("unexpected stats: %#v", version.Stats)
	}
}

func TestListComponents_MapsVulnerabilitySummary(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbom": {
					"components": {
						"nodes": [
							{
								"id": "component-1",
								"name": "openssl",
								"version": "3.0.0",
								"kind": "library",
								"purl": "pkg:generic/openssl@3.0.0",
								"cpes": [],
								"licensesExp": "Apache-2.0",
								"group": "",
								"description": "",
								"primary": false,
								"internal": false,
								"sbomId": "version-1",
								"stats": {
									"vulnStats": {"high": 2},
									"vulnTotalCount": 2
								},
								"updatedAt": "2026-06-02T00:00:00Z"
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

	result, err := client.ListComponents(context.Background(), ListComponentsInput{VersionID: "version-1", First: 10})
	if err != nil {
		t.Fatalf("ListComponents returned error: %v", err)
	}
	if len(result.Components) != 1 || result.Components[0].VulnSummary == nil {
		t.Fatalf("expected component vuln summary: %#v", result.Components)
	}
	summary := result.Components[0].VulnSummary
	if summary.TotalCount != 2 || summary.Stats["high"] != float64(2) {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestDownloadSBOM_MapsCycloneDXDownload(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"sbom": {
					"download": {
						"ready": true,
						"filename": "bom.cdx.json",
						"contentType": "application/json",
						"content": "{\"bomFormat\":\"CycloneDX\"}",
						"processingStatus": {
							"automation": "FINISHED",
							"policyScan": "FINISHED",
							"vulnScan": "FINISHED"
						}
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	includeVulns := true
	result, err := client.DownloadSBOM(context.Background(), DownloadSBOMInput{
		VersionID:        "version-1",
		Spec:             "CycloneDX",
		SpecVersion:      "1.6",
		IncludeVulns:     &includeVulns,
		RequireCompleted: []string{"VULN_SCAN"},
	})
	if err != nil {
		t.Fatalf("DownloadSBOM returned error: %v", err)
	}
	request := gql.requests[0]
	if request["spec"] != "CycloneDX" || request["includeVulns"] != true {
		t.Fatalf("unexpected download variables: %#v", request)
	}
	if result.Filename != "bom.cdx.json" || !result.Ready || result.ProcessingStatus["vulnScan"] != "FINISHED" {
		t.Fatalf("unexpected download result: %#v", result)
	}
}
