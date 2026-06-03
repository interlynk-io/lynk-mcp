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

package graphql

import (
	"regexp"
	"strings"
	"testing"
)

func TestProjectGroupQuery_ProjectsUsesConnectionNodes(t *testing.T) {
	if !regexp.MustCompile(`projects\s*\([^)]*\)\s*\{\s*nodes\s*\{`).MatchString(ProjectGroupQuery) {
		t.Fatal("ProjectGroupQuery must select projects through the ProjectConnection nodes field")
	}
}

func TestSecurityIncidentCreateUpdateMutation_UsesISO8601DateTime(t *testing.T) {
	if !strings.Contains(SecurityIncidentCreateUpdateMutation, "$occurredAt: ISO8601DateTime!") {
		t.Fatal("SecurityIncidentCreateUpdateMutation must use the ISO8601DateTime scalar for occurredAt")
	}
}

func TestSecurityIncidentFindingsQuery_SelectsNestedComponentContext(t *testing.T) {
	for _, required := range []string{"findings(statuses: $statuses)", "component {", "rootSbom {", "projectVersion"} {
		if !strings.Contains(SecurityIncidentFindingsQuery, required) {
			t.Fatalf("SecurityIncidentFindingsQuery missing %q", required)
		}
	}
}
