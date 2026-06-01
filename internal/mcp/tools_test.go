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
