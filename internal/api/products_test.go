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
	"reflect"
	"testing"
)

type fakeGraphQLExecutor struct {
	queries  []string
	requests []map[string]interface{}
	pages    []string
}

func (f *fakeGraphQLExecutor) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	f.queries = append(f.queries, query)
	f.requests = append(f.requests, variables)
	return json.Unmarshal([]byte(f.pages[len(f.requests)-1]), result)
}

func TestListProducts_MapsImportedRepositoryTypesAndTicketingSummary(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"organization": {
					"projectGroups": {
						"nodes": [
							{
								"id": "github-product",
								"name": "GitHub Product",
								"description": "",
								"enabled": true,
								"organizationId": "org-1",
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 1,
								"importedRepository": {
									"__typename": "GithubRepository",
									"id": "repo-github",
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
											"externalIssueTrackerSettings": [
												{"provider": "jira", "projectKey": "SEC"}
											]
										},
										{
											"externalIssueTrackerSettings": [
												{"provider": "jira", "projectKey": "OPS"}
											]
										}
									]
								}
							},
							{
								"id": "bitbucket-product",
								"name": "Bitbucket Product",
								"description": "",
								"enabled": true,
								"organizationId": "org-1",
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 2,
								"importedRepository": {
									"__typename": "BitbucketRepository",
									"id": "repo-bitbucket",
									"name": "repo",
									"fullName": "workspace/repo",
									"slug": "repo",
									"workspace": "workspace",
									"importStatus": "in_progress",
									"webhookEnabled": false
								},
								"projects": {"nodes": []}
							},
							{
								"id": "gitlab-product",
								"name": "GitLab Product",
								"description": "",
								"enabled": true,
								"organizationId": "org-1",
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 3,
								"importedRepository": {
									"__typename": "GitlabRepository",
									"id": "repo-gitlab",
									"name": "repo",
									"fullPath": "group/repo",
									"gitlabId": "123",
									"importStatus": "failed",
									"webhookEnabled": true
								},
								"projects": {"nodes": []}
							},
							{
								"id": "nil-product",
								"name": "Nil Product",
								"description": "",
								"enabled": true,
								"organizationId": "org-1",
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 4,
								"importedRepository": null,
								"projects": {"nodes": []}
							}
						],
						"totalCount": 4,
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

	result, err := client.ListProducts(context.Background(), ListProductsInput{First: 4})
	if err != nil {
		t.Fatalf("ListProducts returned error: %v", err)
	}
	products := result.Products
	if products[0].Repository.Type != "GithubRepository" || products[0].Repository.ImportStatus != "completed" || products[0].Repository.FullName != "org/repo" {
		t.Fatalf("unexpected GitHub repository: %#v", products[0].Repository)
	}
	if products[1].Repository.Type != "BitbucketRepository" || products[1].Repository.ImportStatus != "in_progress" || products[1].Repository.Workspace != "workspace" {
		t.Fatalf("unexpected Bitbucket repository: %#v", products[1].Repository)
	}
	if products[2].Repository.Type != "GitlabRepository" || products[2].Repository.ImportStatus != "failed" || products[2].Repository.FullPath != "group/repo" {
		t.Fatalf("unexpected GitLab repository: %#v", products[2].Repository)
	}
	if products[3].Repository != nil {
		t.Fatalf("nil repository product Repository = %#v, want nil", products[3].Repository)
	}
	summary := products[0].TicketingDefaultsSummary
	if summary == nil || !summary.JiraConfigured || summary.EnvironmentsWithJira != 2 {
		t.Fatalf("unexpected ticketing summary: %#v", summary)
	}
	wantKeys := []string{"OPS", "SEC"}
	if len(summary.JiraProjectKeys) != len(wantKeys) {
		t.Fatalf("JiraProjectKeys = %#v, want %#v", summary.JiraProjectKeys, wantKeys)
	}
	for i := range wantKeys {
		if summary.JiraProjectKeys[i] != wantKeys[i] {
			t.Fatalf("JiraProjectKeys = %#v, want %#v", summary.JiraProjectKeys, wantKeys)
		}
	}
}

func TestListLabels_MapsInputsAndResults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"labels": {
					"nodes": [
						{
							"id": "label-1",
							"name": "Aidash",
							"color": "#0052cc",
							"organizationId": "org-1",
							"createdAt": "2026-01-01T00:00:00Z",
							"updatedAt": "2026-01-02T00:00:00Z"
						}
					],
					"totalCount": 2,
					"pageInfo": {
						"hasNextPage": true,
						"endCursor": "label-cursor-2"
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	result, err := client.ListLabels(context.Background(), ListLabelsInput{
		First:  10,
		After:  "label-cursor-1",
		Search: "Aidash",
	})
	if err != nil {
		t.Fatalf("ListLabels returned error: %v", err)
	}
	if gql.requests[0]["first"] != 10 || gql.requests[0]["after"] != "label-cursor-1" || gql.requests[0]["search"] != "Aidash" {
		t.Fatalf("unexpected request variables: %#v", gql.requests[0])
	}
	if result.TotalCount != 2 || !result.HasNextPage || result.EndCursor != "label-cursor-2" {
		t.Fatalf("unexpected pagination result: %#v", result)
	}
	if len(result.Labels) != 1 || result.Labels[0].ID != "label-1" || result.Labels[0].Name != "Aidash" || result.Labels[0].Color != "#0052cc" {
		t.Fatalf("unexpected labels: %#v", result.Labels)
	}
}

func TestListProducts_SendsLabelIDsAndMapsLabels(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"organization": {
					"projectGroups": {
						"nodes": [
							{
								"id": "product-1",
								"name": "Product 1",
								"description": "",
								"enabled": true,
								"organizationId": "org-1",
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 1,
								"labels": [
									{
										"id": "label-1",
										"name": "Aidash",
										"color": "#0052cc",
										"organizationId": "org-1",
										"createdAt": "2026-01-01T00:00:00Z",
										"updatedAt": "2026-01-02T00:00:00Z"
									}
								],
								"importedRepository": null,
								"projects": {"nodes": []}
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

	result, err := client.ListProducts(context.Background(), ListProductsInput{
		First:    20,
		LabelIDs: []string{"label-1", "label-2"},
	})
	if err != nil {
		t.Fatalf("ListProducts returned error: %v", err)
	}
	if !reflect.DeepEqual(gql.requests[0]["labelIds"], []string{"label-1", "label-2"}) {
		t.Fatalf("labelIds = %#v, want label-1,label-2", gql.requests[0]["labelIds"])
	}
	if len(result.Products) != 1 || len(result.Products[0].Labels) != 1 {
		t.Fatalf("unexpected products: %#v", result.Products)
	}
	label := result.Products[0].Labels[0]
	if label.ID != "label-1" || label.Name != "Aidash" || label.Color != "#0052cc" {
		t.Fatalf("unexpected label: %#v", label)
	}
}

func TestGetProduct_PaginatesEnvironments(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"projectGroup": {
					"id": "product-1",
					"name": "Product 1",
					"description": "",
					"enabled": true,
					"organizationId": "org-1",
					"updatedAt": "2026-04-16T23:00:09Z",
					"sbomsCount": 2,
					"projects": {
						"nodes": [
							{
								"id": "env-1",
								"name": "default",
								"description": "default environment",
								"enabled": true,
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 1
							}
						],
						"pageInfo": {
							"hasNextPage": true,
							"endCursor": "cursor-1"
						}
					}
				}
			}`,
			`{
				"projectGroup": {
					"id": "product-1",
					"name": "Product 1",
					"description": "",
					"enabled": true,
					"organizationId": "org-1",
					"updatedAt": "2026-04-16T23:00:09Z",
					"sbomsCount": 2,
					"projects": {
						"nodes": [
							{
								"id": "env-2",
								"name": "production",
								"description": "production environment",
								"enabled": true,
								"updatedAt": "2026-04-16T23:00:10Z",
								"sbomsCount": 1
							}
						],
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
	product, err := client.GetProduct(context.Background(), "product-1")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}

	if len(product.Environments) != 2 {
		t.Fatalf("expected 2 environments, got %d", len(product.Environments))
	}
	if product.Environments[0].ID != "env-1" || product.Environments[1].ID != "env-2" {
		t.Fatalf("unexpected environments: %#v", product.Environments)
	}
	if len(gql.requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(gql.requests))
	}
	if gql.requests[0]["projectsAfter"] != nil {
		t.Fatalf("first request should not include projectsAfter, got %#v", gql.requests[0]["projectsAfter"])
	}
	if gql.requests[1]["projectsAfter"] != "cursor-1" {
		t.Fatalf("second request projectsAfter = %#v, want cursor-1", gql.requests[1]["projectsAfter"])
	}
}

func TestGetProduct_OrdersEnvironmentsWithVersionsFirst(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"projectGroup": {
					"id": "product-1",
					"name": "Product 1",
					"description": "",
					"enabled": true,
					"organizationId": "org-1",
					"updatedAt": "2026-04-16T23:00:09Z",
					"sbomsCount": 8,
					"projects": {
						"nodes": [
							{
								"id": "env-development",
								"name": "development",
								"description": "development environment",
								"enabled": true,
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 0
							},
							{
								"id": "env-production",
								"name": "production",
								"description": "production environment",
								"enabled": true,
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 0
							},
							{
								"id": "env-default",
								"name": "default",
								"description": "default environment",
								"enabled": true,
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 8
							}
						],
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
	product, err := client.GetProduct(context.Background(), "product-1")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}

	got := []string{
		product.Environments[0].ID,
		product.Environments[1].ID,
		product.Environments[2].ID,
	}
	want := []string{"env-default", "env-development", "env-production"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("environment order = %#v, want %#v", got, want)
		}
	}
}

func TestGetProduct_MapsRepositoryAndJiraDefaults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"projectGroup": {
					"id": "product-1",
					"name": "Product 1",
					"description": "",
					"enabled": true,
					"organizationId": "org-1",
					"updatedAt": "2026-04-16T23:00:09Z",
					"sbomsCount": 1,
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
								"description": "",
								"enabled": true,
								"updatedAt": "2026-04-16T23:00:09Z",
								"sbomsCount": 1,
								"externalIssueTrackerSettings": [
									{
										"id": "linear-setting",
										"provider": "linear",
										"projectKey": "",
										"issueType": "",
										"assignee": "",
										"reporter": "",
										"epic": "",
										"components": null,
										"enableJiraSync": false,
										"updatedAt": "2026-04-16T23:00:09Z"
									},
									{
										"id": "jira-setting",
										"provider": "jira",
										"projectKey": "SEC",
										"issueType": "10001",
										"assignee": "alice",
										"reporter": "bot",
										"epic": "SEC-1",
										"components": ["backend"],
										"enableJiraSync": true,
										"updatedAt": "2026-04-16T23:00:09Z"
									}
								]
							}
						],
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

	product, err := client.GetProduct(context.Background(), "product-1")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}
	if product.Repository == nil || product.Repository.ImportStatus != "completed" {
		t.Fatalf("unexpected repository: %#v", product.Repository)
	}
	defaults := product.Environments[0].JiraDefaults
	if defaults == nil {
		t.Fatal("JiraDefaults = nil, want defaults")
	}
	if defaults.ProjectKey != "SEC" || defaults.IssueType != "10001" || !defaults.EnableSync {
		t.Fatalf("unexpected Jira defaults: %#v", defaults)
	}
}

func TestGetEnvironment_MapsJiraDefaultsAndProductRepository(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"project": {
					"id": "env-1",
					"name": "default",
					"description": "",
					"enabled": true,
					"projectGroupId": "product-1",
					"updatedAt": "2026-04-16T23:00:09Z",
					"sbomsCount": 1,
					"externalIssueTrackerSettings": [
						{
							"id": "jira-setting",
							"provider": "jira",
							"projectKey": "SEC",
							"issueType": "10001",
							"assignee": "alice",
							"reporter": "bot",
							"epic": "SEC-1",
							"components": ["backend"],
							"enableJiraSync": true,
							"updatedAt": "2026-04-16T23:00:09Z"
						}
					],
					"projectGroup": {
						"id": "product-1",
						"name": "Product 1",
						"importedRepository": {
							"__typename": "BitbucketRepository",
							"id": "repo-1",
							"name": "repo",
							"fullName": "workspace/repo",
							"slug": "repo",
							"workspace": "workspace",
							"importStatus": "completed",
							"webhookEnabled": true
						}
					}
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	env, err := client.GetEnvironment(context.Background(), "env-1")
	if err != nil {
		t.Fatalf("GetEnvironment returned error: %v", err)
	}
	if env.JiraDefaults == nil || env.JiraDefaults.ProjectKey != "SEC" {
		t.Fatalf("unexpected Jira defaults: %#v", env.JiraDefaults)
	}
	if env.Product == nil || env.Product.Repository == nil || env.Product.Repository.Type != "BitbucketRepository" {
		t.Fatalf("unexpected product repository: %#v", env.Product)
	}
}
