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

func TestListDoctorResults_MapsInputsAndResults(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"doctorResults": {
					"totalCount": 1,
					"pageInfo": {
						"endCursor": "cursor-2",
						"hasNextPage": true,
						"hasPreviousPage": false,
						"startCursor": "cursor-1"
					},
					"nodes": [
						{
							"checkCode": "IDT-PURL-001",
							"checkName": "PURL conforms to the package-url specification",
							"severity": "high",
							"domain": "identity",
							"componentId": "component-1",
							"componentName": "openssl",
							"componentVersion": "3.0.0",
							"autoFixable": false,
							"findings": [{"field": "purl", "message": "invalid"}]
						}
					]
				}
			}`,
		},
	}
	forceRefresh := false
	client := &Client{gql: gql}

	result, err := client.ListDoctorResults(context.Background(), ListDoctorResultsInput{
		VersionID:     "version-1",
		Search:        "open",
		ComponentID:   "component-1",
		Severity:      []string{"high"},
		Domain:        []string{"identity"},
		CheckCode:     []string{"IDT-PURL-001"},
		ComponentName: []string{"openssl"},
		ForceRefresh:  &forceRefresh,
		First:         25,
		After:         "cursor-1",
	})
	if err != nil {
		t.Fatalf("ListDoctorResults returned error: %v", err)
	}

	request := gql.requests[0]
	if request["sbomId"] != "version-1" {
		t.Fatalf("sbomId = %#v, want version-1", request["sbomId"])
	}
	if request["forceRefresh"] != false {
		t.Fatalf("forceRefresh = %#v, want false", request["forceRefresh"])
	}
	if request["first"] != 25 {
		t.Fatalf("first = %#v, want 25", request["first"])
	}
	if request["after"] != "cursor-1" {
		t.Fatalf("after = %#v, want cursor-1", request["after"])
	}

	if result.TotalCount != 1 {
		t.Fatalf("TotalCount = %d, want 1", result.TotalCount)
	}
	if !result.PageInfo.HasNextPage || result.PageInfo.StartCursor != "cursor-1" || result.PageInfo.EndCursor != "cursor-2" {
		t.Fatalf("unexpected PageInfo: %#v", result.PageInfo)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
	finding := result.Findings[0]
	if finding.CheckCode != "IDT-PURL-001" || finding.ComponentName != "openssl" || finding.Severity != "high" {
		t.Fatalf("unexpected finding: %#v", finding)
	}
}

func TestListDoctorResults_DefaultsFirst(t *testing.T) {
	gql := &fakeGraphQLExecutor{
		pages: []string{
			`{
				"doctorResults": {
					"totalCount": 0,
					"pageInfo": {
						"endCursor": "",
						"hasNextPage": false,
						"hasPreviousPage": false,
						"startCursor": ""
					},
					"nodes": []
				}
			}`,
		},
	}
	client := &Client{gql: gql}

	_, err := client.ListDoctorResults(context.Background(), ListDoctorResultsInput{VersionID: "version-1"})
	if err != nil {
		t.Fatalf("ListDoctorResults returned error: %v", err)
	}

	if gql.requests[0]["first"] != 25 {
		t.Fatalf("first = %#v, want default 25", gql.requests[0]["first"])
	}
}
