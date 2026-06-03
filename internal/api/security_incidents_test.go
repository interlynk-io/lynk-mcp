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

func TestUpdateSecurityIncident_MapsInputsAndResult(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"updateSecurityIncident": {
					"securityIncident": {
						"id": "incident-1",
						"title": "Updated incident",
						"slug": "updated-incident",
						"summary": "",
						"severity": "critical",
						"status": "active",
						"confidence": "confirmed",
						"recommendedActions": "Rotate tokens",
						"sourceUrls": "",
						"createdAt": "2026-05-01T00:00:00Z",
						"updatedAt": "2026-05-02T00:00:00Z",
						"markers": []
					},
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	summary := ""
	recommendedActions := "Rotate tokens"
	result, err := client.UpdateSecurityIncident(context.Background(), UpdateSecurityIncidentInput{
		ID:                 "incident-1",
		Summary:            &summary,
		RecommendedActions: &recommendedActions,
	})
	if err != nil {
		t.Fatalf("UpdateSecurityIncident returned error: %v", err)
	}

	request := gql.requests[0]
	if request["id"] != "incident-1" {
		t.Fatalf("id = %#v, want incident-1", request["id"])
	}
	if request["summary"] != "" {
		t.Fatalf("summary = %#v, want empty string", request["summary"])
	}
	if request["recommendedActions"] != "Rotate tokens" {
		t.Fatalf("recommendedActions = %#v, want Rotate tokens", request["recommendedActions"])
	}

	if result.SecurityIncident == nil {
		t.Fatal("expected updated security incident")
	}
	if result.SecurityIncident.ID != "incident-1" || result.SecurityIncident.RecommendedActions != "Rotate tokens" {
		t.Fatalf("unexpected incident: %#v", result.SecurityIncident)
	}
}

func TestWithdrawSecurityIncidentMarkers_MapsInputsAndResults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"withdrawSecurityIncidentMarkers": {
					"markers": [
						{
							"id": "marker-1",
							"markerType": "purl",
							"purl": "pkg:npm/example@1.0.0",
							"componentName": "",
							"componentVersion": "",
							"githubUrl": "",
							"active": false,
							"addedAt": "2026-05-01T00:00:00Z",
							"withdrawnAt": "2026-05-03T00:00:00Z"
						}
					],
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.WithdrawSecurityIncidentMarkers(
		context.Background(),
		"incident-1",
		[]string{"marker-1"},
	)
	if err != nil {
		t.Fatalf("WithdrawSecurityIncidentMarkers returned error: %v", err)
	}

	request := gql.requests[0]
	if request["securityIncidentId"] != "incident-1" {
		t.Fatalf("securityIncidentId = %#v, want incident-1", request["securityIncidentId"])
	}
	markerIDs := request["markerIds"].([]string)
	if len(markerIDs) != 1 || markerIDs[0] != "marker-1" {
		t.Fatalf("markerIds = %#v, want marker-1", markerIDs)
	}

	if len(result.Markers) != 1 || result.Markers[0].Active {
		t.Fatalf("unexpected markers: %#v", result.Markers)
	}
	if result.Markers[0].WithdrawnAt == nil {
		t.Fatal("expected withdrawnAt to be mapped")
	}
}

func TestCreateSecurityIncidentUpdate_MapsInputsAndResult(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"createSecurityIncidentUpdate": {
					"securityIncidentUpdate": {
						"title": "Containment complete",
						"updateType": "status_changed",
						"body": "",
						"occurredAt": "2026-05-03T10:00:00Z"
					},
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	body := ""
	customerVisible := false
	result, err := client.CreateSecurityIncidentUpdate(context.Background(), CreateSecurityIncidentUpdateInput{
		SecurityIncidentID: "incident-1",
		Title:              "Containment complete",
		UpdateType:         "status_changed",
		OccurredAt:         "2026-05-03T10:00:00Z",
		Body:               &body,
		CustomerVisible:    &customerVisible,
	})
	if err != nil {
		t.Fatalf("CreateSecurityIncidentUpdate returned error: %v", err)
	}

	request := gql.requests[0]
	if request["securityIncidentId"] != "incident-1" || request["updateType"] != "status_changed" {
		t.Fatalf("unexpected update variables: %#v", request)
	}
	if request["body"] != "" || request["customerVisible"] != false {
		t.Fatalf("empty body or false customerVisible not preserved: %#v", request)
	}

	if result.Update == nil {
		t.Fatal("expected timeline update")
	}
	if result.Update.UpdateType != "status_changed" || result.Update.OccurredAt == nil {
		t.Fatalf("unexpected update: %#v", result.Update)
	}
}

func TestGetSecurityIncidentFindings_MapsNestedComponentContext(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"securityIncident": {
					"id": "incident-1",
					"title": "Supply-chain attack",
					"slug": "supply-chain-attack",
					"status": "active",
					"severity": "high",
					"findings": [
						{
							"id": "finding-1",
							"status": "active",
							"matchMethod": "purl",
							"matchedFields": {"purl": true},
							"firstDetectedAt": "2026-05-01T00:00:00Z",
							"lastConfirmedAt": "2026-05-02T00:00:00Z",
							"isPartSbom": true,
							"component": {
								"id": "component-1",
								"name": "example",
								"version": "1.0.0",
								"kind": "library",
								"purl": "pkg:npm/example@1.0.0",
								"cpes": [],
								"licensesExp": "MIT",
								"group": "npm",
								"primary": false,
								"internal": false,
								"sbomId": "sbom-part-1",
								"updatedAt": "2026-05-02T00:00:00Z",
								"sbom": {
									"id": "sbom-part-1",
									"projectVersion": "1.0.0",
									"project": {
										"id": "project-part-1",
										"name": "library",
										"projectGroupId": "product-1"
									}
								}
							},
							"rootSbom": {
								"id": "sbom-root-1",
								"projectVersion": "2.0.0",
								"project": {
									"id": "project-root-1",
									"name": "service",
									"projectGroupId": "product-1"
								}
							}
						}
					]
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.GetSecurityIncidentFindings(context.Background(), SecurityIncidentFindingsInput{
		IncidentID: "incident-1",
		Statuses:   []string{"active"},
	})
	if err != nil {
		t.Fatalf("GetSecurityIncidentFindings returned error: %v", err)
	}

	request := gql.requests[0]
	if request["id"] != "incident-1" {
		t.Fatalf("id = %#v, want incident-1", request["id"])
	}
	statuses := request["statuses"].([]string)
	if len(statuses) != 1 || statuses[0] != "active" {
		t.Fatalf("statuses = %#v, want active", statuses)
	}

	if result.IncidentID != "incident-1" || len(result.Findings) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	finding := result.Findings[0]
	if finding.Component == nil || finding.Component.Name != "example" {
		t.Fatalf("unexpected component: %#v", finding.Component)
	}
	if finding.RootSbom == nil || finding.RootSbom.Project == nil || finding.RootSbom.Project.Name != "service" {
		t.Fatalf("unexpected root SBOM: %#v", finding.RootSbom)
	}
}

func TestSuppressSecurityIncidentFinding_MapsInputsAndFinding(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"suppressSecurityIncidentFinding": {
					"finding": {
						"id": "finding-1",
						"status": "suppressed",
						"matchMethod": "purl",
						"matchedFields": {},
						"firstDetectedAt": "2026-05-01T00:00:00Z",
						"isPartSbom": false,
						"component": null,
						"rootSbom": null
					},
					"errors": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.SuppressSecurityIncidentFinding(context.Background(), SuppressSecurityIncidentFindingInput{
		FindingID: "finding-1",
		Reason:    "False positive",
	})
	if err != nil {
		t.Fatalf("SuppressSecurityIncidentFinding returned error: %v", err)
	}

	request := gql.requests[0]
	if request["findingId"] != "finding-1" || request["reason"] != "False positive" {
		t.Fatalf("unexpected suppression variables: %#v", request)
	}
	if result.Finding == nil || result.Finding.Status != "suppressed" {
		t.Fatalf("unexpected finding: %#v", result.Finding)
	}
}

func TestSecurityIncidentEnabledOrganizations_MapsResult(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"securityIncidentEnabledOrganizations": [
					{
						"id": "org-1",
						"name": "Interlynk Inc",
						"securityIncidentsEnabled": true
					}
				]
			}`,
		},
	}
	client := &Client{gql: gql}

	orgs, err := client.SecurityIncidentEnabledOrganizations(context.Background())
	if err != nil {
		t.Fatalf("SecurityIncidentEnabledOrganizations returned error: %v", err)
	}

	if len(gql.requests) != 1 || gql.requests[0] != nil {
		t.Fatalf("unexpected GraphQL variables: %#v", gql.requests)
	}
	if len(orgs) != 1 || orgs[0].ID != "org-1" || !orgs[0].SecurityIncidentsEnabled {
		t.Fatalf("unexpected orgs: %#v", orgs)
	}
}
