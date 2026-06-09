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

func TestGetTicketingStatus_MapsJiraConfigProductsAndPolicies(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"organization": {
					"connections": {
						"nodes": [
							{
								"id": "connection-1",
								"enabled": true,
								"createdAt": "2026-04-16T23:00:09Z",
								"updatedAt": "2026-04-16T23:00:10Z",
								"connection": {
									"__typename": "JiraConnection",
									"id": "jira-1",
									"url": "https://example.atlassian.net",
									"userName": "bot@example.com",
									"healthCheckStatus": "Success",
									"lastHealthCheckAt": "2026-04-17T23:00:10Z",
									"updatedAt": "2026-04-16T23:00:10Z"
								}
							},
							{
								"id": "connection-2",
								"enabled": true,
								"createdAt": "2026-04-16T23:00:09Z",
								"updatedAt": "2026-04-16T23:00:10Z",
								"connection": {
									"__typename": "LinearConnection",
									"id": "linear-1",
									"url": "https://api.linear.app/graphql",
									"updatedAt": "2026-04-16T23:00:10Z"
								}
							}
						]
					},
					"projectGroups": {
						"nodes": [
							{
								"id": "product-1",
								"name": "Product 1",
								"enabled": true,
								"importedRepository": {
									"__typename": "GithubRepository",
									"id": "repo-1",
									"name": "repo",
									"fullName": "org/repo",
									"owner": "org",
									"defaultBranch": "main",
									"importStatus": "completed",
									"webhookEnabled": true
								},
								"projects": {
									"nodes": [
										{
											"id": "env-1",
											"name": "default",
											"enabled": true,
											"externalIssueTrackerSettings": [
												{
													"id": "setting-1",
													"provider": "jira",
													"projectKey": "SEC",
													"issueType": "Task",
													"assignee": "alice",
													"reporter": "bot",
													"epic": "SEC-1",
													"components": ["backend"],
													"teamId": "",
													"stateId": "",
													"enableJiraSync": true,
													"lastSyncedAt": "2026-04-18T23:00:10Z",
													"lastSyncStatus": "success",
													"updatedAt": "2026-04-18T23:00:10Z"
												},
												{
													"id": "setting-2",
													"provider": "linear",
													"projectKey": "",
													"issueType": "Issue",
													"assignee": "",
													"reporter": "",
													"epic": "",
													"components": null,
													"teamId": "team-1",
													"stateId": "state-1",
													"enableJiraSync": false,
													"lastSyncedAt": null,
													"lastSyncStatus": "",
													"updatedAt": "2026-04-18T23:00:10Z"
												}
											]
										}
									]
								}
							}
						],
						"totalCount": 1,
						"pageInfo": {
							"hasNextPage": true,
							"endCursor": "products-cursor-2"
						}
					}
				},
				"jiraVulnManagementConfig": {
					"id": "config-1",
					"enabled": true,
					"provisioningStatus": "completed",
					"provisioningStep": "",
					"provisioningErrors": "",
					"issueTypeId": "10001",
					"workflowId": "wf-1",
					"screenId": "screen-1",
					"updatedAt": "2026-04-18T23:00:10Z"
				},
				"policies": {
					"nodes": [
						{
							"id": "policy-1",
							"name": "High vulns",
							"isEnabled": true,
							"resultType": "fail",
							"createTicket": true,
							"policyInclusions": [
								{
									"projectId": "env-1",
									"project": {
										"id": "env-1",
										"name": "default",
										"projectGroup": {
											"id": "product-1",
											"name": "Product 1"
										}
									}
								}
							]
						}
					],
					"totalCount": 1,
					"pageInfo": {
						"hasNextPage": true,
						"endCursor": "policies-cursor-2"
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	includeCreatedTickets := false
	status, err := client.GetTicketingStatus(context.Background(), TicketingStatusInput{
		ProductsFirst:         5,
		ProductsAfter:         "products-cursor-1",
		PoliciesFirst:         10,
		PoliciesAfter:         "policies-cursor-1",
		IncludeCreatedTickets: &includeCreatedTickets,
	})
	if err != nil {
		t.Fatalf("GetTicketingStatus returned error: %v", err)
	}

	request := gql.requests[0]
	if request["productsFirst"] != 5 || request["policiesFirst"] != 10 {
		t.Fatalf("unexpected variables: %#v", request)
	}
	if request["productsAfter"] != "products-cursor-1" || request["policiesAfter"] != "policies-cursor-1" {
		t.Fatalf("unexpected cursor variables: %#v", request)
	}
	if request["includeTickets"] != false {
		t.Fatalf("includeTickets = %#v, want false", request["includeTickets"])
	}
	if status.ProductsEndCursor != "products-cursor-2" || !status.ProductsHasNextPage {
		t.Fatalf("unexpected products pagination: cursor=%q hasMore=%t", status.ProductsEndCursor, status.ProductsHasNextPage)
	}
	if status.PoliciesEndCursor != "policies-cursor-2" || !status.PoliciesHasNextPage {
		t.Fatalf("unexpected policies pagination: cursor=%q hasMore=%t", status.PoliciesEndCursor, status.PoliciesHasNextPage)
	}
	if len(status.Connections) != 2 {
		t.Fatalf("len(Connections) = %d, want 2", len(status.Connections))
	}
	if status.Connections[0].Provider != "jira" || status.Connections[0].HealthCheckStatus != "Success" {
		t.Fatalf("unexpected Jira connection: %#v", status.Connections[0])
	}
	if status.Connections[1].Provider != "linear" || status.Connections[1].ProviderID != "linear-1" {
		t.Fatalf("unexpected Linear connection: %#v", status.Connections[1])
	}
	if status.JiraVulnManagementConfig == nil || status.JiraVulnManagementConfig.ProvisioningStatus != "completed" {
		t.Fatalf("unexpected Jira config: %#v", status.JiraVulnManagementConfig)
	}
	if len(status.Products) != 1 || status.Products[0].Repository.FullName != "org/repo" {
		t.Fatalf("unexpected products: %#v", status.Products)
	}
	env := status.Products[0].Environments[0]
	if len(env.IssueTrackerSettings) != 2 {
		t.Fatalf("len(IssueTrackerSettings) = %d, want 2", len(env.IssueTrackerSettings))
	}
	if env.IssueTrackerSettings[0].Provider != "jira" || !env.IssueTrackerSettings[0].EnableSync {
		t.Fatalf("unexpected Jira settings: %#v", env.IssueTrackerSettings[0])
	}
	if env.IssueTrackerSettings[1].Provider != "linear" || env.IssueTrackerSettings[1].TeamID != "team-1" {
		t.Fatalf("unexpected Linear settings: %#v", env.IssueTrackerSettings[1])
	}
	if len(env.AppliedTicketPolicies) != 1 || env.AppliedTicketPolicies[0].ID != "policy-1" {
		t.Fatalf("unexpected applied policies: %#v", env.AppliedTicketPolicies)
	}
}

func TestGetTicketingStatus_ProductFilterKeepsOnlyMatchingPolicyInclusions(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"organization": {
					"connections": {
						"nodes": []
					}
				},
				"jiraVulnManagementConfig": null,
				"projectGroup": {
					"id": "product-1",
					"name": "Product 1",
					"enabled": true,
					"importedRepository": null,
					"projects": {
						"nodes": [
							{
								"id": "env-1",
								"name": "default",
								"enabled": true,
								"externalIssueTrackerSettings": []
							}
						]
					},
					"componentVulns": {
						"nodes": [],
						"totalCount": 42,
						"pageInfo": {
							"hasNextPage": true,
							"endCursor": "tickets-cursor-2"
						}
					}
				},
				"policies": {
					"nodes": [
						{
							"id": "policy-1",
							"name": "Ticket policy",
							"isEnabled": true,
							"resultType": "fail",
							"createTicket": true,
							"policyInclusions": [
								{
									"projectId": "env-1",
									"project": {
										"id": "env-1",
										"name": "default",
										"projectGroup": {
											"id": "product-1",
											"name": "Product 1"
										}
									}
								},
								{
									"projectId": "env-2",
									"project": {
										"id": "env-2",
										"name": "staging",
										"projectGroup": {
											"id": "product-2",
											"name": "Product 2"
										}
									}
								}
							]
						}
					],
					"totalCount": 1,
					"pageInfo": {
						"hasNextPage": true,
						"endCursor": "policies-cursor-2"
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	status, err := client.GetTicketingStatus(context.Background(), TicketingStatusInput{
		ProductID:     "product-1",
		PoliciesAfter: "policies-cursor-1",
		TicketsAfter:  "tickets-cursor-1",
	})
	if err != nil {
		t.Fatalf("GetTicketingStatus returned error: %v", err)
	}

	request := gql.requests[0]
	if request["productId"] != "product-1" {
		t.Fatalf("unexpected variables: %#v", request)
	}
	if _, ok := request["productsFirst"]; ok {
		t.Fatalf("productsFirst set for product-specific query: %#v", request)
	}
	if request["policiesAfter"] != "policies-cursor-1" || request["ticketsAfter"] != "tickets-cursor-1" {
		t.Fatalf("unexpected cursor variables: %#v", request)
	}
	if request["includeTickets"] != true {
		t.Fatalf("includeTickets = %#v, want true", request["includeTickets"])
	}
	if status.PoliciesEndCursor != "policies-cursor-2" || !status.PoliciesHasNextPage {
		t.Fatalf("unexpected policies pagination: cursor=%q hasMore=%t", status.PoliciesEndCursor, status.PoliciesHasNextPage)
	}
	if status.TicketsEndCursor != "tickets-cursor-2" || !status.TicketsHasNextPage || status.TicketsScannedCount != 42 {
		t.Fatalf("unexpected ticket pagination: cursor=%q hasMore=%t count=%d", status.TicketsEndCursor, status.TicketsHasNextPage, status.TicketsScannedCount)
	}
	if len(status.Policies) != 1 {
		t.Fatalf("len(Policies) = %d, want 1", len(status.Policies))
	}
	if len(status.Policies[0].Inclusions) != 1 || status.Policies[0].Inclusions[0].ProductID != "product-1" {
		t.Fatalf("unexpected inclusions: %#v", status.Policies[0].Inclusions)
	}
}
