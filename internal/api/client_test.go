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
	"strings"
	"testing"
)

func TestGetBasicOrganization_UsesMinimalQuery(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"organization": {
					"name": "Interlynk"
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	org, err := client.GetBasicOrganization(context.Background())
	if err != nil {
		t.Fatalf("GetBasicOrganization returned error: %v", err)
	}
	if org.Name != "Interlynk" {
		t.Fatalf("Name = %q, want Interlynk", org.Name)
	}

	query := gql.queries[0]
	if strings.Contains(query, "organizationMetric") {
		t.Fatalf("basic organization query should not request metrics: %s", query)
	}
	for _, field := range []string{"projectCount", "versionCount", "componentCount"} {
		if strings.Contains(query, field) {
			t.Fatalf("basic organization query should not request %s: %s", field, query)
		}
	}
}
