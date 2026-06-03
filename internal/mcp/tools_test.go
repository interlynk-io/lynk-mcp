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

package mcp

import "testing"

func TestSameVexName_NormalizesCommonInputForms(t *testing.T) {
	tests := []struct {
		left  string
		right string
	}{
		{left: "not_affected", right: "Not Affected"},
		{left: "vulnerable_code_not_present", right: "vulnerable-code-not-present"},
		{left: " fixed ", right: "fixed"},
	}

	for _, tt := range tests {
		if !sameVexName(tt.left, tt.right) {
			t.Fatalf("sameVexName(%q, %q) = false, want true", tt.left, tt.right)
		}
	}
}

func TestGetSecurityIncidentMarkerInputsCSV_ParsesRows(t *testing.T) {
	markers, err := getSecurityIncidentMarkerInputsCSV(`marker_type,purl,component_name,component_version,github_url
purl,pkg:npm/example@1.0.0,,,
name_version,,example,1.0.0,https://github.com/example/repo
,,,,`)
	if err != nil {
		t.Fatalf("getSecurityIncidentMarkerInputsCSV returned error: %v", err)
	}

	if len(markers) != 2 {
		t.Fatalf("len(markers) = %d, want 2", len(markers))
	}
	if markers[0].MarkerType != "purl" || markers[0].Purl != "pkg:npm/example@1.0.0" {
		t.Fatalf("unexpected first marker: %#v", markers[0])
	}
	if markers[1].MarkerType != "name_version" || markers[1].ComponentName != "example" || markers[1].GithubURL == "" {
		t.Fatalf("unexpected second marker: %#v", markers[1])
	}
}

func TestGetSecurityIncidentMarkerInputsCSV_RequiresMarkerTypeHeader(t *testing.T) {
	_, err := getSecurityIncidentMarkerInputsCSV("purl\npkg:npm/example@1.0.0\n")
	if err == nil {
		t.Fatal("expected missing marker_type header error")
	}
}

func TestGetSecurityIncidentMarkerInputsCSV_RequiresMarkerTypeValue(t *testing.T) {
	_, err := getSecurityIncidentMarkerInputsCSV("marker_type,purl\n,pkg:npm/example@1.0.0\n")
	if err == nil {
		t.Fatal("expected missing marker_type value error")
	}
}
