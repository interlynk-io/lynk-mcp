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
