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

package config

import (
	"strings"
	"testing"
)

func TestValidateTokenFormat_AcceptsKnownPrefixes(t *testing.T) {
	tokens := []string{
		"lynk_live_abc",
		"lynk_staging_abc",
		"lynk_test_abc",
		"lynk_service_test_abc",
	}

	for _, token := range tokens {
		if !ValidateTokenFormat(token) {
			t.Fatalf("expected token %q to be valid", token)
		}
	}
}

func TestValidateTokenFormat_RejectsUnknownPrefixes(t *testing.T) {
	tokens := []string{
		"",
		"lynk_service_live_abc",
		"lynk_service_staging_abc",
		"other_test_abc",
	}

	for _, token := range tokens {
		if ValidateTokenFormat(token) {
			t.Fatalf("expected token %q to be invalid", token)
		}
	}
}

func TestValidTokenPrefixesDescription_IncludesServiceTestPrefix(t *testing.T) {
	if !strings.Contains(ValidTokenPrefixesDescription(), "lynk_service_test_") {
		t.Fatalf("expected service test prefix in description")
	}
}
