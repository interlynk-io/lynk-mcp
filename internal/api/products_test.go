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
	"testing"
)

type fakeGraphQLExecutor struct {
	requests []map[string]interface{}
	pages    []string
}

func (f *fakeGraphQLExecutor) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	f.requests = append(f.requests, variables)
	return json.Unmarshal([]byte(f.pages[len(f.requests)-1]), result)
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
